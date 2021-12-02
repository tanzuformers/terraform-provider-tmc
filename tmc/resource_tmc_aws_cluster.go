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

func resourceAwsCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsClusterCreate,
		ReadContext:   resourceAwsClusterRead,
		UpdateContext: resourceAwsClusterUpdate,
		DeleteContext: resourceAwsClusterDelete,
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
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the AWS Cluster",
				ForceNew:    true,
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Kubernetes version to be used",
				ForceNew:    true,
				Required:    true,
			},
			"credential_name": {
				Type:        schema.TypeString,
				Description: "Kubernetes version of the AWS Cluster",
				ForceNew:    true,
				Required:    true,
			},
			"pod_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Pods",
				Optional:    true,
				ForceNew:    true,
				Default:     "192.168.0.0/16",
			},
			"service_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Services",
				Optional:    true,
				ForceNew:    true,
				Default:     "10.96.0.0/12",
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Description: "Name of the SSH Keypair used in the AWS Cluster",
				ForceNew:    true,
				Required:    true,
			},
			"control_plane_spec": {
				Type:        schema.TypeList,
				Description: "Contains information related to the Control Plane of the cluster",
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type used to deploy the control plane node",
							ForceNew:    true,
							Required:    true,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Description: "Availability zones of the control plane node",
							Required:    true,
							ForceNew:    true,
							MinItems:    1,
							MaxItems:    3,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": {
							Type:         schema.TypeString,
							Description:  "ID of an existing VPC to be used",
							Optional:     true,
							ForceNew:     true,
							RequiredWith: []string{"control_plane_spec.0.private_subnets", "control_plane_spec.0.public_subnets"},
							ExactlyOneOf: []string{"control_plane_spec.0.vpc_cidrblock", "control_plane_spec.0.vpc_id"},
						},
						"vpc_cidrblock": {
							Type:         schema.TypeString,
							Description:  "CIDR block used by the Cluster's VPC",
							Optional:     true,
							ForceNew:     true,
							ExactlyOneOf: []string{"control_plane_spec.0.vpc_cidrblock", "control_plane_spec.0.vpc_id"},
						},
						"private_subnets": {
							Type:         schema.TypeList,
							Description:  "IDs of the private subnets in the specified availability zones",
							Optional:     true,
							ForceNew:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							RequiredWith: []string{"control_plane_spec.0.vpc_id", "control_plane_spec.0.public_subnets"},
						},
						"public_subnets": {
							Type:         schema.TypeList,
							Description:  "IDs of the public subnets in the specified availability zones",
							Optional:     true,
							ForceNew:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							RequiredWith: []string{"control_plane_spec.0.vpc_id", "control_plane_spec.0.private_subnets"},
						},
					},
				},
			},
		},
	}
}

func resourceAwsClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})
	cluster_group := d.Get("cluster_group").(string)
	control_plane_spec := d.Get("control_plane_spec").([]interface{})[0].(map[string]interface{})
	azs := control_plane_spec["availability_zones"].([]interface{})

	if len(azs) == 2 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create AWS Cluster",
			Detail:   "number of availability zones must be either 1 for a development cluster or 3 for a highly available cluster",
		})
		return diags
	}

	if _, ok := d.GetOk("control_plane_spec.0.vpc_id"); ok {
		pvt_subnets := control_plane_spec["private_subnets"].([]interface{})
		pub_subnets := control_plane_spec["public_subnets"].([]interface{})
		if len(pvt_subnets) != len(azs) || len(pub_subnets) != len(azs) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create AWS Cluster",
				Detail:   "number of private subnets and public subnets must be equal to the number of availability zones specified",
			})
			return diags
		}
	}

	opts := &tanzuclient.ClusterOpts{
		Region:           d.Get("region").(string),
		Version:          d.Get("version").(string),
		CredentialName:   d.Get("credential_name").(string),
		ControlPlaneSpec: control_plane_spec,
		PodCidrBlock:     d.Get("pod_cidrblock").(string),
		ServiceCidrBlock: d.Get("service_cidrblock").(string),
		SshKey:           d.Get("ssh_key").(string),
	}

	cluster, err := client.CreateCluster(clusterName, managementClusterName, provisionerName, cluster_group, description, labels, opts)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create AWS cluster",
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
			Summary:  "Failed to create AWS cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(cluster.Meta.UID)

	resourceAwsClusterRead(ctx, d, m)

	return diags
}

func resourceAwsClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	cluster, err := client.GetCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS cluster",
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

	return diags
}

func resourceAwsClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})
	cluster_group := d.Get("cluster_group").(string)
	resourceVersion := d.Get("resource_version").(string)
	control_plane_spec := d.Get("control_plane_spec").([]interface{})[0].(map[string]interface{})
	azs := control_plane_spec["availability_zones"].([]interface{})

	if control_plane_spec["availability_zones"] != nil {
		if len(azs) == 2 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create AWS Cluster",
				Detail:   "number of availability zones must be either 1 for a development cluster or 3 for a highly available cluster",
			})
			return diags
		}
	}

	if _, ok := d.GetOk("control_plane_spec.0.vpc_id"); ok {
		pvt_subnets := control_plane_spec["private_subnets"].([]interface{})
		pub_subnets := control_plane_spec["public_subnets"].([]interface{})
		if len(pvt_subnets) != len(azs) || len(pub_subnets) != len(azs) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create AWS Cluster",
				Detail:   "number of private subnets and public subnets must be equal to the number of availability zones specified",
			})
			return diags
		}
	}

	opts := &tanzuclient.ClusterOpts{
		Region:           d.Get("region").(string),
		Version:          d.Get("version").(string),
		CredentialName:   d.Get("credential_name").(string),
		ControlPlaneSpec: control_plane_spec,
		PodCidrBlock:     d.Get("pod_cidrblock").(string),
		ServiceCidrBlock: d.Get("service_cidrblock").(string),
		SshKey:           d.Get("ssh_key").(string),
	}

	if d.HasChange("labels") || d.HasChange("cluster_group") {
		_, err := client.UpdateCluster(clusterName, managementClusterName, provisionerName, cluster_group, description, resourceVersion, labels, opts)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to update AWS cluster",
				Detail:   fmt.Sprintf("Error updating resource %s: %s", d.Get("name"), err),
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceAwsClusterRead(ctx, d, m)

}

func resourceAwsClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	err := client.DeleteCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete AWS cluster",
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
			Summary:  "Failed to delete AWS cluster",
			Detail:   fmt.Sprintf("Error waiting to delete resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId("")

	return diags
}

func flatten_aws_control_plane_spec(s *tanzuclient.ClusterSpec) map[string]interface{} {
	cp_spec := make(map[string]interface{})

	cp_spec["instance_type"] = s.TkgAws.Topology.ControlPlane.InstanceType
	cp_spec["availability_zones"] = s.TkgAws.Topology.ControlPlane.AvailabilityZones
	cp_spec["vpc_cidrblock"] = s.TkgAws.Settings.Network.Provider.Vpc.CidrBlock
	cp_spec["vpc_id"] = s.TkgAws.Settings.Network.Provider.Vpc.Id

	var pvt_subnets []string
	var pub_subnets []string
	for _, subnet := range s.TkgAws.Settings.Network.Provider.Subnets {
		if subnet.IsPublic {
			pub_subnets = append(pub_subnets, subnet.Id)
		} else {
			pvt_subnets = append(pvt_subnets, subnet.Id)
		}
	}
	cp_spec["public_subnets"] = pub_subnets
	cp_spec["private_subnets"] = pvt_subnets

	return cp_spec
}
