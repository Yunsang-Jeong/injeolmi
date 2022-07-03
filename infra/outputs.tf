output "apigateway_execution_arn" {
  value = aws_api_gateway_stage.this.execution_arn
}

output "apigateway_invoke_url" {
  value = aws_api_gateway_stage.this.invoke_url
}

output "dynamodb_id" {
  value = aws_dynamodb_table.this.id
}