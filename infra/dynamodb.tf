resource "aws_dynamodb_table" "this" {
  name = join(local.delimiter, [local.name_tag_prefix, "dynamodb"])
  billing_mode   = "PAY_PER_REQUEST"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "CommentID"
  table_class    = "STANDARD"

  attribute {
    name = "CommentID"
    type = "N"
  }

  attribute {
    name = "MRID"
    type = "N"
  }

  attribute {
    name = "MRIID"
    type = "N"
  }

  attribute {
    name = "PipelineID"
    type = "N"
  }

  attribute {
    name = "ActionType"
    type = "S"
  }

  attribute {
    name = "ActionOptions"
    type = "S"
  }

  tags = {
    Name        = join(local.delimiter, [local.name_tag_prefix, "dynamodb"])
  }
}