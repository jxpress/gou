terraform {
  backend "s3" {
    key    = "terraform/dynamodb.tfstate"
    region = "ap-northeast-1"
  }
}
