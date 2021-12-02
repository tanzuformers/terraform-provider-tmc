package tmc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// labelsSchema returns the schema to use for labels.
//
func labelsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// Ignore changes to the creator label added automatically added by TMC and
			// also ignore changes when the labels field itself is deleted when updating
			return k == "labels.tmc.cloud.vmware.com/creator" || k == "labels.%"
		},
	}
}

func labelsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

func labelsSchemaImmutable() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
	}
}
