---
page_title: "TMC: tmc_workspace"
layout: "tmc"
subcategory: "ClusterGroups"
description: |-
  Creates and manages a clustergroup in the TMC platform
---

# Resource `tmc_cluster_group`

Create a Tanzu Cluster Group for your organisation. The name of the Cluster Group is unique and cannot be modified.

```terraform
resource "tmc_clustergroup" "cluster_group" {
  name        = "tf-group"
  description = "test desc"
  labels = {
    "test" = "one"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu ClusterGroup.
* `description` - (Optional) The description of the Tanzu ClusterGroup.
* `labels` - (Optional) A map of labels to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu ClusterGroup