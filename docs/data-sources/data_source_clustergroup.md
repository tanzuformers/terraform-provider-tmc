---
page_title: "TMC: tmc_clustergroup"
subcategory: "ClusterGroups"
description: |-
  Get information on a list of Tanzu Mission Control (TMC) ClusterGroup.
---

# Data Source `tmc_cluster_group`

Data Source to get the details of a Tanzu Cluster Group.

## Example Usage

```terraform
data "tmc_clustergroup" "name" {
  name = "cluster-name"
}
```

## Argument Reference

* `name` - (Required) The name of the clustergroup to lookup in the TMC platform. If no clustergroup is found with this name, an error will be returned.

## Attributes Reference

* `description` - Description of the found clustergroup.
* `labels` - A mapping of labels of the resource.