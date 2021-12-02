terraform {
  required_version = ">= 0.15"
}

provider "tmc" {
}

resource "tmc_observability_credential" "example" {
  name              = "foo"
  observability_url = "mock-url"
  api_token         = "mock-token"
}

data "tmc_observability_credential" "example" {
  name = tmc_observability_credential.example.name
}