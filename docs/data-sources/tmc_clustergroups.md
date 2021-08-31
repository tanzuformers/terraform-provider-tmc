---
page_title: "tmc_cluster_groups Data Source - terraform-provider-tmc"
subcategory: ""
description: |-
  Data Source to get the details of multiple Tanzu Cluster Groups.
---

# Data Source `tmc_cluster_groups`

Data Source to get the details of multiple Tanzu Cluster Groups.

## Example Usage

```terraform
data "tmc_clustergroups" "name" {
  
}
```

## Schema

### Optional

- **match_labels** (Map, Optional) Set of labels to filter the available Cluster Groups.