package tmc

import (
	"context"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceClusterGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterGroupsRead,
		Schema: map[string]*schema.Schema{
			"match_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Set of labels to search for in the cluster groups",
			},
			"cluster_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Tanzu Cluster Groups",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of the Tanzu Cluster Group",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique Name of the Tanzu Cluster Group in your Org",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional Description for the Tanzu Cluster Group",
						},
						"labels": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Optional Labels for the Tanzu Cluster Group",
						},
					},
				},
			},
		},
	}
}

func dataSourceClusterGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	var diags diag.Diagnostics

	labels := d.Get("match_labels").(map[string]interface{})

	res, err := client.GetAllClusterGroups(labels)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterGroups := make([]interface{}, len(*res))

	for i, clusterGroup := range *res {
		clusterGroups[i] = flattenClusterGroup(&clusterGroup)
	}

	if err := d.Set("cluster_groups", clusterGroups); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().String())

	return diags
}

func flattenClusterGroup(cg *tanzuclient.ClusterGroup) *map[string]interface{} {
	return &map[string]interface{}{
		"id":          cg.Meta.UID,
		"name":        cg.FullName.Name,
		"description": cg.Meta.Description,
		"labels":      cg.Meta.Labels,
	}
}
