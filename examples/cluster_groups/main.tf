terraform {
  required_version = ">= 0.15"
}

provider "tmc" {
}

resource "tmc_cluster_groups" "example" {
  name = "foo"

  description = "Terraform provider acceptance testing cluster group"

  labels = {
    env = "test"
    createdby = "Terraform"
    purpose = "Acceptance testing for Terraform tmc provider"
    repo_url = "https://github.com/codaglobal/terraform-provider-tmc"
  }
}

data "tmc_cluster_group" "example" {
  name = tmc_cluster_group.example.name
}

data "tmc_cluster_groups" "all" {}