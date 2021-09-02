---
page_title: "Tanzu Mission Control Provider"
subcategory: ""
layout: "tmc"
description: |-
  Use the Tanzu Mission Control (TMC) provider to interact with the many resources supported by TMC. You must configure the provider with the proper credentials before you can use it.

---

# Tanzu Mission Control (TMC) provider

Use the Tanzu Mission Control (TMC) provider to interact with the many resources supported by TMC. You must configure the provider with the proper credentials before you can use it.

Use the navigation to the left to read about the available resources.

To learn more about the Tanzu Mission Control and the resources supported by it, refer the offical [VMware Tanzu Mission Control documentation](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/index.html).


## Example Usage

Terraform 0.15 and later:

```terraform
terraform {
  required_providers {
    tmc = {
      source  = "codaglobal/tmc"
      version = "~> 0.1"
    }
  }
}

# Configure the TMC Provider
provider "tmc" {
  org_url = "my-org-url"
}

# Create a Tanzu Workspace
resource "tmc_workspace" "example" {
  name = "example"
}
```

## Authentication

We generally require the VMware Cloud Console service url unique to every organization and an API token generated from the console to authenticate to Tanzu Mission Control.
The TMC provider offers a the following methods of providing credentials for
authentication, in the given order:

- Static credentials
- Environment variables

### Static Credentials

!> **Warning:** Hard-coded credentials are not recommended in any Terraform
configuration and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be added in-line as `org_url` and `api_token` to the provider configuration block:

Usage:

```terraform
provider "tmc" {
  org_url = "my-org-url"
  api_token = "my-api-token"
}
```

### Environment Variables

You can provide your credentials via the `TMC_ORG_URL` and
`TMC_API_TOKEN`, environment variables, representing your AWS
VMware Cloud Console url and API token, respectively:

```terraform
provider "tmc" {}
```

Usage:

```sh
$ export TMC_API_TOKEN="yourapitoken"
$ export TMC_ORG_URL="yourorgurl"
$ terraform plan
```

