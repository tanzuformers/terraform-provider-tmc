package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster",
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource version of the Cluster",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Cluster",
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
				Description: "Description of the Cluster",
			},
			"management_cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of an existing management cluster to be used",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of an existing provisioner to be used",
			},
			"cluster_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the cluster group",
			},
			"labels": labelsSchema(),
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})
	cluster_group := d.Get("cluster_group").(string)

	opts := &tanzuclient.ClusterOpts{}

	cluster, err := client.CreateCluster(clusterName, managementClusterName, provisionerName, cluster_group, description, labels, opts)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"PENDING",
			"CREATING",
			"PROCESSING",
		},
		Target: []string{
			"READY",
		},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.GetCluster(clusterName, managementClusterName, provisionerName)
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status.Phase, nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(cluster.Meta.UID)

	resourceClusterRead(ctx, d, m)

	return diags
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	cluster, err := client.GetCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("resource_version", cluster.Meta.ResourceVersion)
	d.Set("description", cluster.Meta.Description)
	d.Set("cluster_group", cluster.Spec.ClusterGroupName)

	if err := d.Set("labels", cluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster",
			Detail:   fmt.Sprintf("Error getting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	return diags
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})
	cluster_group := d.Get("cluster_group").(string)
	resourceVersion := d.Get("resource_version").(string)

	opts := &tanzuclient.ClusterOpts{}

	if d.HasChange("labels") || d.HasChange("cluster_group") {
		_, err := client.UpdateCluster(clusterName, managementClusterName, provisionerName, cluster_group, description, resourceVersion, labels, opts)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to update cluster",
				Detail:   fmt.Sprintf("Error updating resource %s: %s", d.Get("name"), err),
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceClusterRead(ctx, d, m)

}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	err := client.DeleteCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete cluster",
			Detail:   fmt.Sprintf("Error deleting resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"DELETING",
		},
		Target: []string{
			"DELETED",
		},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.DescribeCluster(clusterName, managementClusterName, provisionerName)
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Phase, nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete cluster",
			Detail:   fmt.Sprintf("Error waiting to delete resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId("")

	return diags
}
