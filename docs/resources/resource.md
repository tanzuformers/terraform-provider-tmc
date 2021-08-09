---
page_title: "tmc_resource Resource - terraform-provider-tmc"
subcategory: ""
description: |-
  Sample resource in the Terraform provider tmc.
---

# Resource `tmc_resource`

Sample resource in the Terraform provider tmc.

## Example Usage

```terraform
resource "tmc_resource" "example" {
  sample_attribute = "foo"
}
```

## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **sample_attribute** (String, Optional) Sample attribute.


