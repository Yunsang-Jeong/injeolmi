resource "aws_dynamodb_table" "this" {
  name = join(local.delimiter, [local.name_tag_prefix, "dynamodb"])
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "CommentID"
  table_class    = "STANDARD"

  attribute {
    name = "CommentID"
    type = "N"
  }

  tags = {
    Name        = join(local.delimiter, [local.name_tag_prefix, "dynamodb"])
  }
}