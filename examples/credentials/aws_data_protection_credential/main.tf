terraform {
  required_version = ">= 0.15"
}

provider "tmc" {
}

resource "tmc_aws_data_protection_credential" "example" {
  name         = "foo"
  iam_role_arn = "arn:aws:iam::xxxxxxxxxx:role/mock_role"
}

data "tmc_aws_data_protection_credential" "example" {
  name = tmc_aws_data_protection_credential.example.name
}