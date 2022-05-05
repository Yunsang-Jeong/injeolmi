################################################################################
# API-Gateway

resource "aws_api_gateway_rest_api" "this" {
  name = join(local.delimiter, [local.name_tag_prefix, "apigateway"])
}

resource "aws_api_gateway_rest_api_policy" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id

  policy = data.aws_iam_policy_document.apigateway_resource_policy.json
}
################################################################################


################################################################################
# Resource

resource "aws_api_gateway_resource" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = "events"
}

resource "aws_api_gateway_request_validator" "this" {
  name                        = "this"
  rest_api_id                 = aws_api_gateway_rest_api.this.id
  validate_request_body       = false
  validate_request_parameters = true
}
################################################################################


################################################################################
# Method

resource "aws_api_gateway_method" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this.id

  authorization = "NONE"
  http_method   = "POST"

  request_validator_id = aws_api_gateway_request_validator.this.id
  # https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-method-settings-method-request.html#setup-method-request-parameters
  # A method parameter can be a path parameter, a header, or a query string parameter.
  # method.request.{location}.{name}, where location is querystring, path, or header and name is a valid and unique parameter name.
  request_parameters = {
    "method.request.header.X-Gitlab-Event" = true,
    "method.request.header.X-Gitlab-Token" = true
  }
}

resource "aws_api_gateway_method_response" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this.id
  http_method = aws_api_gateway_method.this.http_method
  status_code = "200"
}

################################################################################


################################################################################
# Integration

resource "aws_api_gateway_integration" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this.id
  http_method = aws_api_gateway_method.this.http_method
  type        = "AWS_PROXY"

  integration_http_method = "POST"
  uri                     = aws_lambda_function.this.invoke_arn
  cache_key_parameters    = []
  request_parameters      = {}
  request_templates       = {}
}
################################################################################



################################################################################
# Deployment

resource "aws_api_gateway_stage" "this" {
  deployment_id = aws_api_gateway_deployment.this.id
  rest_api_id   = aws_api_gateway_rest_api.this.id
  stage_name    = "prod"
}

resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.this.body))
  }

  depends_on = [
    aws_api_gateway_integration.this
  ]

  lifecycle {
    create_before_destroy = true
  }
}
################################################################################


################################################################################
# Enable logging to Cloudwatch

resource "aws_api_gateway_account" "this" {
  cloudwatch_role_arn = aws_iam_role.apigateway_service_role.arn
}

resource "aws_iam_role" "apigateway_service_role" {
  name               = "${var.service_name}-apigateway-service-role"
  assume_role_policy = data.aws_iam_policy_document.apigateway_service_role.json
}

resource "aws_iam_role_policy" "apigateway_cloudwatch" {
  name   = "cloudwatch"
  role   = aws_iam_role.apigateway_service_role.id
  policy = data.aws_iam_policy_document.cloudwatch_policy.json
}

resource "aws_api_gateway_method_settings" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  method_path = "*/*"

  settings {
    logging_level = "INFO"
  }
}
################################################################################