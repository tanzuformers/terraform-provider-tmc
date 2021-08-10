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
		},
	}
}

func dataSourceTmcWorkspacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	res, err := client.GetAllWorkspaces()
	if err != nil {
		return diag.FromErr(err)
	}

	workspaces := make([]interface{}, len(res))

	for i, workspace := range res {
		workspaces[i] = workspace.FullName.Name
	}

	if err := d.Set("names", workspaces); err != nil {
		return diag.FromErr(err)
	}

	// Check if a different suitable is available
	d.SetId(time.Now().UTC().String())
	return diags
}
