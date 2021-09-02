terraform {
  required_version = ">= 0.15"
}

provider "tmc" {
}

resource "tmc_workspace" "example" {
  name = "foo"

  description = "Terraform provider acceptance testing workspace"

  labels = {
    env = "test"
    createdby = "Terraform"
    purpose = "Acceptance testing for Terraform tmc provider"
    repo_url = "https://github.com/codaglobal/terraform-provider-tmc"
  }
}

data "tmc_workspace" "example" {
  name = tmc_workspace.example.name
}

data "tmc_workspaces" "all" {}