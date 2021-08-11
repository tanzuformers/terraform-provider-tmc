terraform {
  required_providers {
    tmc = {
      source  = "coda-global/tanzu/tmc"
    }
  }
}

resource "tmc_clustergroup" "cluster_group" {
  name = "tf-group"
  description = "test desc"
  labels = {
      "test" = "one"
  }
}
