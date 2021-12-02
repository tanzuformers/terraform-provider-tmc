package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcClusterGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmcClusterGroupCreate,
		ReadContext:   resourceTmcClusterGroupRead,
		UpdateContext: resourceTmcClusterGroupUpdate,
		DeleteContext: resourceTmcClusterGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Cluster Group",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique Name of the Tanzu Cluster Group in your Org",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !IsValidTanzuName(v) {
						errs = append(errs, fmt.Errorf("name should contain only lowercase letters, numbers or hyphens and should begin with either an alphabet or number"))
					}
					return
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional Description for the Tanzu Cluster Group",
			},
			"labels": labelsSchema(),
		},
	}
}

func resourceTmcClusterGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterGroupName := d.Get("name").(string)
	desc := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})

	if !IsValidTanzuName(clusterGroupName) {
		return InvalidTanzuNameError("cluster group")
	}

	clusterGroup, err := client.CreateClusterGroup(clusterGroupName, desc, labels)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create cluster group",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(clusterGroup.Meta.UID)

	resourceTmcClusterGroupRead(ctx, d, m)

	return nil
}

func resourceTmcClusterGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterGroupName := d.Get("name").(string)

	clusterGroup, err := client.GetClusterGroup(clusterGroupName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster group",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("description", clusterGroup.Meta.Description)
	if err := d.Set("labels", clusterGroup.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster group",
			Detail:   fmt.Sprintf("Error setting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.Set("id", clusterGroup.Meta.UID)

	return nil
}

func resourceTmcClusterGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterGroupName := d.Get("name").(string)

	if d.HasChange("description") || d.HasChange("labels") {
		desc := d.Get("description").(string)
		labels := d.Get("labels").(map[string]interface{})

		_, err := client.UpdateClusterGroup(clusterGroupName, desc, labels)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to update cluster group",
				Detail:   fmt.Sprintf("Error updating resource %s: %s", d.Get("name"), err),
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceTmcClusterGroupRead(ctx, d, m)
}

func resourceTmcClusterGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	cgName := d.Get("name").(string)

	err := client.DeleteClusterGroup(cgName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete cluster group",
			Detail:   fmt.Sprintf("Error deleting resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId("")

	return nil
}
