package tmc

import (
	"context"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTmcWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcWorkspacesRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Names of the All Tanzu Workspaces",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Set of labels to filter the workspaces",
			},
			"ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "UID of the All Tanzu Workspaces",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTmcWorkspacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	labels := d.Get("labels").(map[string]interface{})

	res, err := client.GetAllWorkspaces(labels)
	if err != nil {
		return diag.FromErr(err)
	}

	workspaceNames := make([]interface{}, len(res))

	for i, workspace := range res {
		workspaceNames[i] = workspace.FullName.Name
	}

	if err := d.Set("names", workspaceNames); err != nil {
		return diag.FromErr(err)
	}

	// Check if a different suitable is available
	d.SetId(time.Now().UTC().String())
	return diags
}
