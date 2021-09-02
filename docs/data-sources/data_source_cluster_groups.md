---
page_title: "TMC: tmc_cluster_groups"
layout: "tmc"
subcategory: "Cluster Groups"
description: |-
  Get information on a list of Tanzu Mission Control (TMC) Cluster Groups.
---

# Data Source `tmc_cluster_groups`

Use this data source to get list of all cluster groups in Tanzu Mission Control (TMC).

## Example Usage
# List all the cluster groups in the Tanzu platform.
```terraform
data "tmc_cluster_groups" "example" {
  lables = {
    Environment = "test" // Fetches only the cluster groups which match the labels
  }
}
```

## Arugment Reference

* `labels` - (Optional) Map of labels to filter only the cluster groups that match them.


## Attributes Reference

* `names` - List of all the cluster group names in the TMC platform, suitable for referencing in other resources.

* `ids` - List of Unique Identifiers (UID) of all the cluster groups in the TMC platform, suitable for referencing in other resources.