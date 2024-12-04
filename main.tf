# Create the API Gateway HTTP API
resource "aws_apigatewayv2_api" "http_api" {
  name          = "product-api"
  protocol_type = "HTTP"
}


locals {
  lambda_functions = ["productHandler"]
}

data "archive_file" "zip_the_lambda_code" {
  for_each = toset(local.lambda_functions)

  type        = "zip"
  source_dir  = "${path.module}/functions/${each.key}"
  output_path = "${path.module}/functions/${each.key}.zip"
  excludes    = ["**/*.zip"]
}


# Define the Lambda function
resource "aws_lambda_function" "product-handler" {
  function_name = "productHandler"
  handler       = "index.handler"
  runtime       = "nodejs14.x"
  filename      = data.archive_file.zip_the_lambda_code["productHandler"].output_path
  role          = "arn:aws:iam::123456789012:role/service-role/lambda-execution-role"
}

# Create the API Gateway stage
resource "aws_apigatewayv2_stage" "api_gateway" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "v1"
  auto_deploy = true
}

# Create the API Gateway Lambda integration
resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.http_api.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.product-handler.invoke_arn
}

# Define the /product route and associate it with the Lambda integration
resource "aws_apigatewayv2_route" "product_route" {
  api_id        = aws_apigatewayv2_api.http_api.id
  route_key     = "POST /product"
  target        = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

# Allow API Gateway to invoke the Lambda function
resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowApiGatewayInvoke"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  function_name = aws_lambda_function.product-handler.function_name
}
