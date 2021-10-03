package tmc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func dataSourceTmcProvisioner() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcProvisionerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Provisioner",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Provisioner",
			},
			"management_cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Management Cluster which contains the Tanzu Provisioner",
			},
			"labels": labelsSchemaComputed(),
		},
	}
}

func dataSourceTmcProvisionerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	provisioner, err := client.GetProvisioner(d.Get("management_cluster_name").(string), d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", provisioner.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read provisioner",
			Detail:   fmt.Sprintf("Error setting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.SetId(string(provisioner.Meta.UID))

	return diags
}
