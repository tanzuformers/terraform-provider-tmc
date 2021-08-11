data "tmc_clustergroup" "name" {
  name = "sandbox-orl"
}

output "clustergroup_name" {
  value = data.tmc_clustergroup.name
}
