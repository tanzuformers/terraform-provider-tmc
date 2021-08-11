package tmc

import (
	"context"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceClusterGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterGroupCreate,
		ReadContext:   resourceClusterGroupRead,
		UpdateContext: resourceClusterGroupUpdate,
		DeleteContext: resourceClusterGroupDelete,
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
	}
}

func resourceClusterGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	name := d.Get("name").(string)
	desc := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})

	clusterGroup, err := client.CreateClusterGroup(name, desc, labels)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create ClusterGroup Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(clusterGroup.Meta.UID)

	resourceClusterGroupRead(ctx, d, m)

	return nil
}

func resourceClusterGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	cgName := d.Get("name").(string)

	clusterGroup, err := client.GetClusterGroup(cgName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Read ClusterGroup Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	err = d.Set("description", clusterGroup.Meta.Description)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Set Description ClusterGroup Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	err = d.Set("labels", clusterGroup.Meta.Labels)
	if err != nil {
		message := err.Error()
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  message,
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("id", clusterGroup.Meta.UID)

	return nil
}
func resourceClusterGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	cgName := d.Get("name").(string)

	if d.HasChange("description") || d.HasChange("labels") {
		desc := d.Get("description").(string)
		labels := d.Get("labels").(map[string]interface{})

		_, err := client.UpdateClusterGroup(cgName, desc, labels)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Update ClusterGroup Failed",
				Detail:   "Cannot Update Cluster Group with the given values",
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceClusterGroupRead(ctx, d, m)
}
func resourceClusterGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	cgName := d.Get("name").(string)

	err := client.DeleteClusterGroup(cgName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete ClusterGroup Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}
