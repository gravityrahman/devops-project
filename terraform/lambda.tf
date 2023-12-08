# Variable to get a unique name across Gitlab.
variable "gl_project" {
  type = string
}

locals {
  # The name can have weird characters so we just strip everything.
  safe_gl_project = replace(var.gl_project, "/[^a-zA-Z]/", "_")
}


# Creates a policy the Lambda runtime can assume.
resource "aws_iam_role" "lambda" {
  name_prefix = "lr-"

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

# Adds the basic AWS execution role
# https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html#permissions-executionrole-features
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Read the lambda out of S3 so we can dynamically make changes.
data "aws_s3_bucket_object" "lambda_version" {
  bucket = "altr-interview-data" # This is hard coded for easiness.
  key    = "${var.gl_project}-app.zip"
}

# Define the lambda function.
resource "aws_lambda_function" "go_app" {
  function_name = "gofunc-${local.safe_gl_project}"
  role          = aws_iam_role.lambda.arn
  handler       = "app" # This is the built package name inside the zip.

  runtime = "go1.x"

  s3_bucket = data.aws_s3_bucket_object.lambda_version.bucket
  s3_key = data.aws_s3_bucket_object.lambda_version.key
  s3_object_version = data.aws_s3_bucket_object.lambda_version.version_id

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = aws_dynamodb_table.table.name
      DYNAMODB_ITEM_HASH = "this-is-a-key"
    }
  }
}

# Allow the API Gateway to invoke the lambda.
resource "aws_lambda_permission" "lambda" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.go_app.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*"
}
