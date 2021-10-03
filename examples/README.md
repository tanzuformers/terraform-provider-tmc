# TMC Provider Examples

This directory contains examples of TMC Resources that can be run/tested manually via the Terraform CLI. The examples each have their own README containing more details on what the example does.

To run any example, clone the repository and run `terraform apply` within the example's own directory.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool.

* **provider/provider.tf** example file for the provider index page
* **data-sources/<full data source name>/data-source.tf** example file for the named data source page
* **resources/<full resource name>/resource.tf** example file for the named data source page

In this scenario, we embed example usage for resources directly in the documentation itself and use these examples for purposes of testing.