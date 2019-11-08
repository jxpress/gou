#####################################
# DynamoDB
#####################################
resource "aws_dynamodb_table" "karma-dynamodb-table" {
  name           = "Karma"
  read_capacity  = "${var.karma_read_capacity}"
  write_capacity = "${var.karma_write_capacity}"
  hash_key       = "identifier"

  attribute {
    name = "identifier"
    type = "S"
  }

  tags {
    Name = "Karma"
    env  = "${var.env}"
  }
}
