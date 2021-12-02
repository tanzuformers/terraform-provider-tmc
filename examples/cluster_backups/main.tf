terraform {
  required_version = ">= 0.15"
}

provider "tmc" {
}

resource "tmc_cluster_backup" "example" {
  name                    = "cluster-bkp-example"
  management_cluster_name = "aws-hosted"
  provisioner_name        = "aws-provisioner"
  cluster_name            = "tmc-test-cluster"
  retention_period        = "3600s"
  snapshot_volumes        = false
  storage_location        = "aws-cluster-backup-target"
}

data "tmc_cluster_backup" "example" {
  name                    = tmc_cluster_backup.example.name
  management_cluster_name = tmc_cluster_backup.example.management_cluster_name
  provisioner_name        = tmc_cluster_backup.example.provisioner_name
  cluster_name            = tmc_cluster_backup.example.cluster_name
}