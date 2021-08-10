package tmc

import (
	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTmcWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTmcWorkspaceRead,
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
		},
	}
}

func dataSourceTmcWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*tanzuclient.Client)

	workspace, err := client.GetWorkspace(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.Set("description", workspace.Meta.Description)
	d.SetId(string(workspace.Meta.UID))

	return nil
}
