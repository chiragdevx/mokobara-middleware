locals {
  lambda_functions = ["productHandler", "orderHandler"]
}

# ============================= S3 BUCKET ============================
resource "aws_s3_bucket" "mokobara_state" {
  bucket = "mokobara-state-bucket"
  tags = {
    Name        = "Terraform State Bucket"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket_versioning" "versioning" {
  bucket = aws_s3_bucket.mokobara_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "encryption" {
  bucket = aws_s3_bucket.mokobara_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "public_access" {
  bucket = aws_s3_bucket.mokobara_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# ============================= DYNAMODB TABLE ============================
resource "aws_dynamodb_table" "mokobara_lock" {
  name         = "mokobara-lock-table"
  hash_key     = "LockID"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name        = "Terraform Lock Table"
    Environment = "Dev"
  }
}

terraform {
  backend "s3" {
    bucket         = "mokobara-state-bucket"
    key            = "./terraform.tfstate"
    region         = "us-west-1"
    dynamodb_table = "mokobara-lock-table"
    encrypt        = true
  }
}

# ============================= API GATEWAY ============================
resource "aws_apigatewayv2_api" "http_api" {
  name          = "mokobara-api-gateway"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.http_api.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.lambda_function["productHandler"].invoke_arn
}

resource "aws_apigatewayv2_integration" "order_lambda_integration" {
  api_id             = aws_apigatewayv2_api.http_api.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.lambda_function["orderHandler"].invoke_arn
}

resource "aws_apigatewayv2_stage" "api_gateway_stage" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "v1"
  auto_deploy = true
}

resource "aws_apigatewayv2_route" "product_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "POST /product"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "order_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "POST /order"
  target    = "integrations/${aws_apigatewayv2_integration.order_lambda_integration.id}"
}

resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowApiGatewayInvoke"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  function_name = aws_lambda_function.lambda_function["productHandler"].function_name
}

resource "aws_lambda_permission" "allow_api_gateway_order" {
  statement_id  = "AllowApiGatewayInvokeOrder"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  function_name = aws_lambda_function.lambda_function["orderHandler"].function_name
}

# ======================= LAMBDA ROLE AND POLICY =======================
resource "aws_iam_role" "lambda_execution_role" {
  name               = "aws_lambda_role"
  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
 ]
}
EOF
}

resource "aws_iam_policy" "iam_policy_for_lambda" {
  name        = "aws_iam_policy_for_aws_lambda_role"
  path        = "/"
  description = "AWS IAM Policy for managing AWS Lambda role"
  policy      = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_default_policy_attachment" {
  role       = aws_iam_role.lambda_execution_role.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda.arn
}

# ======================= LAMBDA FUNCTION =======================
data "archive_file" "zip_the_lambda_code" {
  for_each = toset(local.lambda_functions)

  type        = "zip"
  source_dir  = "${path.module}/functions/${each.key}"
  output_path = "${path.module}/functions/${each.key}.zip"
  excludes    = ["**/*.zip"]
}

resource "aws_lambda_function" "lambda_function" {
  for_each = toset(local.lambda_functions)

  function_name = each.key
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = data.archive_file.zip_the_lambda_code[each.key].output_path
  role          = aws_iam_role.lambda_execution_role.arn
  timeout       = 180

  architectures    = ["x86_64"]
  source_code_hash = data.archive_file.zip_the_lambda_code[each.key].output_base64sha256

  environment {
    variables = {
      BASE_URL      = var.base_url
      URL_TOKEN     = var.url_token
      STORE_NAME    = var.store_name
      SHOPIFY_TOKEN = var.shopify_token
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}
