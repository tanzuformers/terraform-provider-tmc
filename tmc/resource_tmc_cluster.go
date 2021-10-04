package tmc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceTmcClusterRead,
		CreateContext: resourceTmcClusterCreate,
		UpdateContext: resourceTmcClusterUpdate,
		DeleteContext: resourceTmcClusterDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Cluster",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "",
				Description: "Description of the Cluster",
			},
			"management_cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the management cluster used",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the provisioner",
			},
			"cluster_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Name of the cluster group",
			},
			"installer_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to install the agent",
			},
			"labels": labelsSchema(),
			"tkg_aws": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Details of Cluster hosted on AWS",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Region of the AWS Cluster",
							Optional:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Provisioner credential used to create the cluster",
							Optional:    true,
						},
						"credential_name": {
							Type:        schema.TypeString,
							Description: "Kubernetes version of the AWS Cluster",
							Optional:    true,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Description: "Availability zones of the control plane node",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type used to deploy the control plane node",
							Optional:    true,
						},
						"vpc_cidrblock": {
							Type:        schema.TypeString,
							Description: "CIDR block used by the Cluster's VPC",
							Optional:    true,
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Description: "Name of the SSH Keypair used in the AWS Cluster",
							Optional:    true,
						},
						// "pods_cidrblocks": {
						// 	Type:        schema.TypeList,
						// 	Description: "CIDR blocks allocated to the pods in the cluster",
						// 	Optional:    true,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"cidr_blocks": {
						// 				Type:     schema.TypeString,
						// 				Computed: true,
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},
		},
	}
}

func resourceTmcClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	var diags diag.Diagnostics

	cluster, err := client.GetCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", cluster.Meta.Description)
	d.Set("cluster_group_name", cluster.Spec.ClusterGroupName)
	d.Set("installer_link", cluster.Status.InstallerLink)

	if err := d.Set("labels", cluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read clustergroup",
			Detail:   fmt.Sprintf("Error getting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	tkgAws := make([]interface{}, 0)

	awsData := flattenAwsData(cluster)

	tkgAws = append(tkgAws, awsData)

	if err := d.Set("tkg_aws", tkgAws); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster",
			Detail:   fmt.Sprintf("Error setting spec for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.SetId(string(cluster.Meta.UID))

	return diags
}

func resourceTmcClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	description := d.Get("description").(string)
	managementCluster := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	clusterGroupName := d.Get("cluster_group_name").(string)
	labels := d.Get("labels").(map[string]interface{})

	cluster, err := client.CreateCluster(clusterName, description, managementCluster, provisionerName, clusterGroupName, labels)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(cluster.Meta.UID))

	return resourceTmcClusterRead(ctx, d, meta)
}

func resourceTmcClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	clusterName := d.Get("name").(string)
	managementCluster := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	if d.HasChange("description") || d.HasChange("labels") || d.HasChange("cluster_group_name") {
		description := d.Get("description").(string)
		labels := d.Get("labels").(map[string]interface{})
		clusterGroupName := d.Get("cluster_group_name").(string)

		_, err := client.UpdateCluster(clusterName, description, managementCluster, provisionerName, clusterGroupName, labels)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to update cluster",
				Detail:   fmt.Sprintf("Cannot update the cluster %s with the new values: %s", d.Get("name").(string), err),
			})
			return diags
		}
	}

	return resourceTmcClusterRead(ctx, d, meta)
}

func resourceTmcClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := client.DeleteCluster(d.Get("name").(string), d.Get("management_cluster").(string), d.Get("provisioner_name").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete cluster",
			Detail:   fmt.Sprintf("Cannot delete given cluster %s: %s", d.Get("name").(string), err),
		})
		return diags
	}

	d.SetId("")

	return diags
}
