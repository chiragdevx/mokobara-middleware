terraform {
  backend "s3" {
    bucket         = "mokobara-state-bucket"
    key            = "terraform.tfstate"
    region         = "us-west-1"
    dynamodb_table = "mokobara-lock-table"
    encrypt        = true
  }
}

provider "aws" {
  region  = "us-west-1"
  profile = "devx"
}