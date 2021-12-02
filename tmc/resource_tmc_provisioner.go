package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcProvisioner() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmcProvisionerCreate,
		ReadContext:   resourceTmcProvisionerRead,
		UpdateContext: resourceTmcProvisionerUpdate,
		DeleteContext: resourceTmcProvisionerDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Provisioner",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique Name of the Tanzu Provisioner in your Org",
			},
			"management_cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Management Cluster which contains the Tanzu Provisioner",
			},
			"labels": labelsSchema(),
		},
	}
}

func resourceTmcProvisionerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	provisionerName := d.Get("name").(string)
	mgmtClusterName := d.Get("management_cluster_name").(string)
	labels := d.Get("labels").(map[string]interface{})

	if !IsValidTanzuName(provisionerName) {
		return InvalidTanzuNameError("provisioner")
	}

	provisioner, err := client.CreateProvisioner(mgmtClusterName, provisionerName, labels)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create Provisioner Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(provisioner.Meta.UID)

	resourceTmcProvisionerRead(ctx, d, m)

	return nil
}

func resourceTmcProvisionerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	provisionerName := d.Get("name").(string)
	mgmtClusterName := d.Get("management_cluster_name").(string)

	provisioner, err := client.GetProvisioner(mgmtClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Read Provisioner Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	if err := d.Set("labels", provisioner.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read Provisioner",
			Detail:   fmt.Sprintf("Error setting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.Set("id", provisioner.Meta.UID)

	return nil
}

func resourceTmcProvisionerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	provisionerName := d.Get("name").(string)

	if d.HasChange("labels") || d.HasChange("management_cluster_name") {
		mgmtClusterName := d.Get("management_cluster_name").(string)
		labels := d.Get("labels").(map[string]interface{})

		_, err := client.UpdateProvisioner(mgmtClusterName, provisionerName, labels)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Update Provisioner Failed",
				Detail:   "Cannot Update Provisioner with the given values",
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceTmcProvisionerRead(ctx, d, m)
}

func resourceTmcProvisionerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	provisionerName := d.Get("name").(string)
	mgmtClusterName := d.Get("management_cluster_name").(string)

	err := client.DeleteProvisioner(mgmtClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete Provisioner Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}
