---
page_title: "tmc_cluster_group Resource - terraform-provider-tmc"
subcategory: ""
description: |-
    Resource to create a Tanzu Cluster Group for your organisation
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

## Schema

### Required

- **name** (String, Required) Name of the cluster group. This name is unique and an organisation cannot have more than one cluster group with the same name.

### Optional

- **description** (String, Optional) Description for the cluster group
- **labels** (Map, Optional) Key-Value pairs representing the labels to be applied to the cluster group