---
page_title: "TMC: tmc_workspace"
layout: "tmc"
subcategory: "Workspaces"
description: |-
  Creates and manages a workspace in the TMC platform
---

# Resource `tmc_workspace`

The TMC Workspace resource allows requesting the creation of a workspace in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the workspace.

## Example Usage

```terraform
resource "tmc_workspace" "example" {
  name = "example"

  description = "Example description"

  labels = {
    Environment = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu Workspace. Changing the name forces recreation of this resource.
* `description` - (Optional) The description of the Tanzu Workspace.
* `labels` - (Optional) A map of labels to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Workspace

