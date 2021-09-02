---
page_title: "TMC: tmc_cluster_group"
layout: "tmc"
subcategory: "Cluster Groups"
description: |-
  Creates and manages a Cluster Group in the TMC platform
---

# Resource `tmc_cluster_group`

The TMC Cluster Group resource allows requesting the creation of a cluster group in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the cluster group.

```terraform
resource "tmc_cluster_group" "example" {
  name        = "example"
  description = "Example description"
  labels = {
    Environment = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu Cluster Group. Changing the name forces recreation of this resource.
* `description` - (Optional) The description of the Tanzu Cluster Group.
* `labels` - (Optional) A map of labels to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Cluster Group.