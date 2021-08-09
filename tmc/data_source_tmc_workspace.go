package spotify

import (
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
	client := meta.(*tmc.Client)

	workspace, err := client.GetWorkspace(d.Get("name"))
	if err != nil {
		return err
	}

	d.Set("name", workspace.Name)
	d.Set("description", workspace.Description)
	d.SetId(string(workspace.UID))

	return nil
}
