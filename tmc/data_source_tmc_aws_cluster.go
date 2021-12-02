package tmc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func dataSourceAwsCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAwsClusterRead,
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
				Computed:    true,
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
			"cluster_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the cluster group",
			},
			"labels": labelsSchemaComputed(),
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the AWS Cluster",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Provisioner credential used to create the cluster",
				Computed:    true,
			},
			"credential_name": {
				Type:        schema.TypeString,
				Description: "Kubernetes version of the AWS Cluster",
				Computed:    true,
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Description: "Name of the SSH Keypair used in the AWS Cluster",
				Computed:    true,
			},
			"pod_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Pods",
				Computed:    true,
			},
			"service_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Services",
				Computed:    true,
			},
			"control_plane_spec": {
				Type:        schema.TypeList,
				Description: "Contains information related to the Control Plane of the cluster",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type used to deploy the control plane node",
							Computed:    true,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Description: "Availability zones of the control plane node",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "ID of an existing VPC to be used",
							Computed:    true,
						},
						"vpc_cidrblock": {
							Type:        schema.TypeString,
							Description: "CIDR block used by the Cluster's VPC",
							Computed:    true,
						},
						"private_subnets": {
							Type:        schema.TypeList,
							Description: "IDs of the private subnets in the specified availability zones",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"public_subnets": {
							Type:        schema.TypeList,
							Description: "IDs of the public subnets in the specified availability zones",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAwsClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	var diags diag.Diagnostics

	cluster, err := client.GetCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS cluster",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("description", cluster.Meta.Description)
	d.Set("cluster_group", cluster.Spec.ClusterGroupName)

	if err := d.Set("labels", cluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS cluster",
			Detail:   fmt.Sprintf("Error getting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("region", cluster.Spec.TkgAws.Distribution.Region)
	d.Set("credential_name", cluster.Spec.TkgAws.Distribution.ProvisionerCredentialName)
	d.Set("version", cluster.Spec.TkgAws.Distribution.Version)
	d.Set("ssh_key", cluster.Spec.TkgAws.Settings.Security.SshKey)
	d.Set("pod_cidrblock", cluster.Spec.TkgAws.Settings.Network.ClusterNetwork.Pods[0].CidrBlocks)
	d.Set("service_cidrblock", cluster.Spec.TkgAws.Settings.Network.ClusterNetwork.Services[0].CidrBlocks)

	cp_spec := flatten_aws_control_plane_spec(cluster.Spec)
	spec := make([]map[string]interface{}, 0)
	spec = append(spec, cp_spec)

	if err := d.Set("control_plane_spec", spec); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS cluster",
			Detail:   fmt.Sprintf("Error getting control plane information for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(string(cluster.Meta.UID))

	return diags
}
