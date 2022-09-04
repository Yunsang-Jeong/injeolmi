################################################################################
# lambda

resource "aws_lambda_function" "this" {
  # General
  function_name = join(local.delimiter, [local.name_tag_prefix, "lambda"])

  # Runtime
  role    = aws_iam_role.lambda_service_role.arn
  runtime = "go1.x"
  handler = "main"

  # Source code
  filename         = "../bin/main.zip"
  package_type     = "Zip"
  publish          = true
  source_code_hash = data.archive_file.zip.output_base64sha256
  # filebase64sha256("../bin/injeolmi.zip")

  environment {
    variables = {
      "GITLAB_TOKEN"          = var.gitlab_token
      "GITLAB_WEBHOOK_SECRET" = var.gitlab_webhook_secret
    }
  }
}
################################################################################


################################################################################
# Allow api-gateway to execute the lambda

resource "aws_lambda_permission" "apigateway" {
  # General
  function_name = aws_lambda_function.this.function_name

  # Permission
  principal    = "apigateway.amazonaws.com"
  action       = "lambda:InvokeFunction"
  statement_id = "AllowExecutionFromAPIGateway"
  source_arn   = "${aws_api_gateway_rest_api.this.execution_arn}/*/*"
}
################################################################################


################################################################################
# Service role for the lambda
resource "aws_iam_role" "lambda_service_role" {
  name               = join(local.delimiter, [local.name_tag_prefix, "lambda", "service", "role"])
  assume_role_policy = data.aws_iam_policy_document.lambda_service_role.json
}

resource "aws_iam_role_policy" "lambda_cloudwatch" {
  name   = "cloudwatch"
  role   = aws_iam_role.lambda_service_role.id
  policy = data.aws_iam_policy_document.cloudwatch_policy.json
}

resource "aws_iam_role_policy" "lambda_dynamodb" {
  name   = "dynamodb"
  role   = aws_iam_role.lambda_service_role.id
  policy = data.aws_iam_policy_document.dynamodb_policy.json
}

################################################################################

resource "null_resource" "makefile" {
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = "make -f ../Makefile build"
  }
}

data "archive_file" "zip" {
  type        = "zip"
  source_dir  = "../bin/"
  output_path = "../bin/main.zip"

  depends_on = [
    null_resource.makefile
  ]
}