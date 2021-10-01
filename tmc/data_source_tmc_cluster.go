package tmc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterRead,
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
			"cluster_group_name": {
				Type:        schema.TypeString,
				Computed: true,
				Description: "Name of the cluster group",
			},
			"installer_link": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"labels": labelsSchemaComputed(),
			"tkg_aws": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of Cluster hosted on AWS",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"availability_zones": {
							Type:        schema.TypeList,
							Description: "Availability zones of the control plane node",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type used to deploy the control plane node",
							Computed:    true,
						},
						"vpc_cidrblock": {
							Type:        schema.TypeString,
							Description: "CIDR block used by the Cluster's VPC",
							Computed:    true,
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Description: "Name of the SSH Keypair used in the AWS Cluster",
							Computed:    true,
						},
						// "pods_cidrblocks": {
						// 	Type:        schema.TypeList,
						// 	Description: "CIDR blocks allocated to the pods in the cluster",
						// 	Computed:    true,
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

func dataSourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func flattenAwsData(data *tanzuclient.Cluster) map[string]interface{} {
	aws := make(map[string]interface{})

	// pods_cidrs := make([]interface{}, 0)
	// pods_cidrs = append(pods_cidrs, data.Spec.TkgAws.Settings.Network.ClusterNetwork.Pods)

	aws["availability_zones"] = data.Spec.TkgAws.Topology.ControlPlane.AvailabilityZones
	aws["instance_type"] = data.Spec.TkgAws.Topology.ControlPlane.InstanceType
	aws["vpc_cidrblock"] = data.Spec.TkgAws.Settings.Network.Provider.Vpc.CidrBlock
	aws["region"] = data.Spec.TkgAws.Distribution.Region
	aws["credential_name"] = data.Spec.TkgAws.Distribution.ProvisionerCredentialName
	aws["version"] = data.Spec.TkgAws.Distribution.Version
	aws["ssh_key"] = data.Spec.TkgAws.Settings.Security.SshKey
	// aws["pods_cidrblocks"] = pods_cidrs

	return aws
}

// func sliceToStrMap(elements []string) map[string]string {
// 	elementMap := make(map[string]string)
// 	for _, s := range elements {
// 		elementMap[s] = s
// 	}
// 	return elementMap
// }
