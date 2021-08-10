package tmc

import (
	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceClusterGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster Group",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Cluster Group",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Cluster Group",
			},
		},
	}
}

func dataSourceClusterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*tanzuclient.Client)

	clusterGroup, err := client.GetClusterGroup(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.Set("description", clusterGroup.Meta.Description)
	d.SetId(string(clusterGroup.Meta.UID))

	return nil
}
