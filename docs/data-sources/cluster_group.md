---
page_title: "TMC: tmc_cluster_group"
layout: "tmc"
subcategory: "Cluster Groups"
description: |-
  Get information on a list of Tanzu Mission Control (TMC) Cluster Group.
---

# Data Source `tmc_cluster_group`

Use this data source to get the details about a cluster group in TMC platform.

## Example Usage
# Get details of a cluster group in the Tanzu platform.
```terraform
data "tmc_cluster_group" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the clustergroup to lookup in the TMC platform. If no cluster group is found with this name, an error will be returned.

## Attributes Reference

* `id` - Unique Identifiers (UID) of the found cluster group in the TMC platform.
* `description` - Description of the found cluster group.
* `labels` - A mapping of labels of the resource.