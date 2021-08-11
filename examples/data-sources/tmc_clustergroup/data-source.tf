terraform {
  required_providers {
    tmc = {
      source = "coda-global/tanzu/tmc"
    }
  }
}

data "tmc_clustergroups" "name" {
  match_labels = {
    "cloud" = "aws"
    "test" = "one"
  }  
}

output "clustergroup_name" {
  value = data.tmc_clustergroups.name
}
