---
page_title: "TMC: tmc_clustergroups"
subcategory: "ClusterGroups"
description: |-
  Get information on a list of Tanzu Mission Control (TMC) ClusterGroups.
---

# Data Source `tmc_cluster_groups`

Use this data source to get list of all clustergroups in Tanzu Mission Control (TMC).

## Example Usage

```terraform
data "tmc_clustergroups" "name" {
  
}
```

## Attributes Reference

* `names` - List of all the clustergroup names in the TMC platform, suitable for referencing in other resources.

* `labels` - (Optional) Map of labels to filter clustergroups matching them
