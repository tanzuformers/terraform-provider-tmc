data "tmc_clustergroups" "test" {
  match_labels = {
    "cloud" = "aws"
  }  
}
data "tmc_clustergroups" "test2" {
  
}

output "clustergroup_name" {
  value = data.tmc_clustergroups.test
}

output "clustergroup_name" {
  value = data.tmc_clustergroups.test2
}
