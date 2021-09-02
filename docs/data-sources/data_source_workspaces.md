---
page_title: "TMC: tmc_workspaces"
layout: "tmc"
subcategory: "Workspaces"
description: |-
  Get information on a list of Tanzu Mission Control (TMC) Workspaces
---

# Data Source `tmc_workspaces`

Use this data source to get list of all workspaces in Tanzu Mission Control (TMC).

## Example Usage
# List all the workspaces in the Tanzu platform.
```terraform
data "tmc_workspaces" "example" {
}
```

## Attributes Reference

* `names` - List of all the workspace names in the TMC platform, suitable for referencing in other resources.

