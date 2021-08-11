resource "tmc_clustergroup" "cluster_group" {
  name        = "tf-group"
  description = "test desc"
  labels = {
    "test" = "one"
  }
}
