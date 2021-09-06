package tmc

import (
	"context"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTmcProvisioners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcProvisionersRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Names of the All Tanzu Provisioners under a Management Cluster",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "UID of the All Tanzu Provisioners under a Management Cluster",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"management_cluster_name": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Management Cluster Name of the Tanzu Provisioners",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"labels": labelsSchema(),
		},
	}
}

func dataSourceTmcProvisionersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	labels := d.Get("labels").(map[string]interface{})
	mgmtClusterName := d.Get("management_cluster_name").(string)

	res, err := client.GetAllProvisioners(mgmtClusterName, labels)
	if err != nil {
		return diag.FromErr(err)
	}

	provisionerNames := make([]interface{}, len(res))
	provisionerIds := make([]interface{}, len(res))

	for i, provisioner := range res {
		provisionerNames[i] = provisioner.FullName.SimpleFullName.Name
		provisionerIds[i] = provisioner.Meta.UID
	}

	if err := d.Set("names", provisionerNames); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ids", provisionerIds); err != nil {
		return diag.FromErr(err)
	}

	// Check if a different suitable is available
	d.SetId(time.Now().UTC().String())
	return diags
}
