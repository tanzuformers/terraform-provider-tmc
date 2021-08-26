---
page_title: "tmc_cluster_group Data Source - terraform-provider-tmc"
subcategory: ""
description: |-
  Data Source to get the details of a Tanzu Cluster Group.
---

# Data Source `tmc_cluster_group`

Data Source to get the details of a Tanzu Cluster Group.

## Example Usage

```terraform
data "tmc_clustergroup" "name" {
  name = "cluster-name"
}
```

## Schema

### Required

- **name** (String, Required) Name of the Cluster Group.