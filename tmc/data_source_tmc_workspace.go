package tmc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func dataSourceTmcWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Workspace",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Workspace",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Tanzu Workspace",
			},
			"labels": labelsSchemaComputed(),
		},
	}
}

func dataSourceTmcWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workspace, err := client.GetWorkspace(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", workspace.Meta.Description)
	if err := d.Set("labels", workspace.Meta.SimpleMetaData.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read workspace",
			Detail:   fmt.Sprintf("Error setting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.SetId(string(workspace.Meta.SimpleMetaData.UID))

	return diags
}
