---
page_title: "TMC: tmc_workspace"
layout: "tmc"
subcategory: "Workspaces"
description: |-
  Get information on a specific Tanzu Mission Control (TMC) Workspace
---

# Data Source `tmc_workspace`

Use this data source to get the details about a workspace in TMC platform.

## Example Usage
# Get details of a workspace in the Tanzu platform.
```terraform
data "tmc_workspace" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the workspace to lookup in the TMC platform. If no workspace is found with this name, an error will be returned.

## Attributes Reference

* `id` - Unique Identifiers (UID) of the found workspace in the TMC platform.
* `description` - Description of the found workspace.
* `labels` - A mapping of labels of the resource.