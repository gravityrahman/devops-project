resource "aws_dynamodb_table" "table" {
  name = "my-table-${local.safe_gl_project}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "MyKey"

  attribute {
    name = "MyKey"
    type = "S"
  }
}

resource "aws_dynamodb_table_item" "item" {
  table_name = aws_dynamodb_table.table.name
  hash_key   = aws_dynamodb_table.table.hash_key

  item = <<ITEM
{
  "MyKey": {"S": "this-is-a-key"},
  "MyContent": {"S": "This was read from DynamoDB!"}
}
ITEM
}
