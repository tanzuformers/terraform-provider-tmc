package tmc

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func IsValidTanzuName(name string) bool {
	// Only allow lowercase alphanumric characters and hyphens to be present in the names of resources
	// Begin with only letters or numbers
	return regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`).MatchString(name)
}

func InvalidTanzuNameError(resourceName string) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("Failed to create %s", resourceName),
		Detail:   "Invalid Resource Name. Name must start and end with a letter or number, and can contain only lowercase letters, numbers, and hyphens.",
	})

	return diags
}
