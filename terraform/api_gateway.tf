# Stands up an API Gateway that sends everything to our Go Lambda.
# https://registry.terraform.io/modules/terraform-aws-modules/apigateway-v2/aws/latest
module "api_gateway" {
  source = "terraform-aws-modules/apigateway-v2/aws"
  version = "1.1.0"

  name          = "dev-http-${local.safe_gl_project}"
  description   = "My awesome HTTP API Gateway"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }

  create_api_domain_name = false

  # Routes and integrations
  integrations = {
    "POST /" = {
      lambda_arn             = aws_lambda_function.go_app.arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 12000
    }

    "GET /" = {
      lambda_arn             = aws_lambda_function.go_app.arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 12000
    }

    "$default" = {
      lambda_arn = aws_lambda_function.go_app.arn
    }
  }
}

output "api_gateway_uri" {
  value = module.api_gateway.apigatewayv2_api_api_endpoint
}
