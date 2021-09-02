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
  lables = {
    env = "test" // Fetches only the workspaces which match the labels
  }
}
```

## Argument Reference

* `labels` - (Optional) Map of labels to filter only the workspaces that match them.

## Attributes Reference

* `names` - List of all the workspace names in the TMC platform, suitable for referencing in other resources.

* `ids` - List of Unique Identifiers (UID) of all the workspaces in the TMC platform, suitable for referencing in other resources.
