---
page_title: "tmc_data_source Data Source - terraform-provider-tmc"
subcategory: ""
description: |-
  Sample data source in the Terraform provider tmc.
---

# Data Source `tmc_data_source`

Sample data source in the Terraform provider tmc.

## Example Usage

```terraform
data "tmc_data_source" "example" {
  sample_attribute = "foo"
}
```

## Schema

### Required

- **sample_attribute** (String, Required) Sample attribute.

### Optional

- **id** (String, Optional) The ID of this resource.


