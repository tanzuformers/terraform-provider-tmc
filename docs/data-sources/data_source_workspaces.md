---
page_title: "TMC: tmc_workspaces"
layout: "tmc"
subcategory: "Workspaces"
description: |-
  Get information on a list of Tanzu Mission Control (TMC) Workspaces
---

# Data Source `tmc_workspaces`

Use this data source to get the details about  of a certificate in AWS Certificate
Manager (ACM), you can reference
it by domain without having to hard code the ARNs as input.

## Example Usage
# List all the workspaces in the Tanzu platform.
```terraform
data "tmc_workspaces" "example" {
}
```

## Attributes Reference

* `names` - List of all the workspace names in the TMC platform, suitable for referencing in other resources.

