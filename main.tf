# ============================= API GATEWAY ============================
resource "aws_apigatewayv2_api" "http_api" {
  name          = "product-api"
  protocol_type = "HTTP"
}

locals {
  lambda_functions = ["productHandler"]
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.http_api.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.product-handler.invoke_arn
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

resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowApiGatewayInvoke"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  function_name = aws_lambda_function.product-handler.function_name
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

resource "aws_lambda_function" "product-handler" {
  function_name = "productHandler"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = data.archive_file.zip_the_lambda_code["productHandler"].output_path
  role          = aws_iam_role.lambda_execution_role.arn


  architectures    = ["x86_64"]
  source_code_hash = data.archive_file.zip_the_lambda_code["productHandler"].output_base64sha256

  lifecycle {
    create_before_destroy = true
  }
}
