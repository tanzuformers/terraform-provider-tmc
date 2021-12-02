terraform {
  required_version = ">= 0.15"
}

provider "tmc" {
}

resource "tmc_aws_storage_credential" "example" {
  name              = "foo"
  access_key_id     = "mock"
  secret_access_key = "mock-key"
}

data "tmc_aws_storage_credential" "example" {
  name = tmc_aws_storage_credential.example.name
}