terraform {
  required_version = ">= 0.15"
}

provider "tmc" {
}

resource "tmc_clustergroups" "example" {
  name = "foo"

  description = "Terraform provider acceptance testing clustergroup"

  labels = {
    env = "test"
    createdby = "Terraform"
    purpose = "Acceptance testing for Terraform tmc provider"
    repo_url = "https://github.com/codaglobal/terraform-provider-tmc"
  }
}

data "tmc_clustergroup" "example" {
  name = tmc_clustergroup.example.name
}

data "tmc_clustergroups" "all" {}