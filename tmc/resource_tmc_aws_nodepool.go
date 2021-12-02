package tmc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceAwsNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsNodePoolCreate,
		ReadContext:   resourceAwsNodePoolRead,
		UpdateContext: resourceAwsNodePoolUpdate,
		DeleteContext: resourceAwsNodePoolDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Nodepool",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Nodepool in the cluster",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !IsValidTanzuName(v) {
						errs = append(errs, fmt.Errorf("name should contain only lowercase letters, numbers or hyphens and should begin with either an alphabet or number"))
					}
					return
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the cluster in which the nodepool is present",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique ID of the cluster in which the nodepool is present",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the Nodepool",
			},
			"management_cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the management cluster used",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the provisioner",
			},
			"node_labels":  labelsSchemaImmutable(),
			"cloud_labels": labelsSchemaImmutable(),
			"worker_node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of worker nodes in the nodepool",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Availability zone of the worker node",
				Required:    true,
			},
			"instance_type": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Instance type used to deploy the worker node",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Kubernetes version to be used",
				Required:    true,
			},
		},
	}
}

func resourceAwsNodePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*tanzuclient.Client)

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	cloud_labels := d.Get("cloud_labels").(map[string]interface{})
	node_labels := d.Get("node_labels").(map[string]interface{})
	cluster_name := d.Get("cluster_name").(string)
	worker_node_count := d.Get("worker_node_count").(int)

	awsNodeSpec := &tanzuclient.AwsNodeSpec{
		AvailabilityZone: d.Get("availability_zone").(string),
		Version:          d.Get("version").(string),
		InstanceType:     d.Get("instance_type").(string),
	}

	nodepool, err := client.CreateNodePool(npName, managementClusterName, provisionerName, cluster_name, description, cloud_labels, node_labels, worker_node_count, awsNodeSpec)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create AWS nodepool",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"CREATING",
			"WAITING",
			"UPGRADING",
		},
		Target: []string{
			"READY",
		},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.GetNodePool(npName, cluster_name, managementClusterName, provisionerName)
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

	d.SetId(nodepool.Meta.UID)

	resourceAwsNodePoolRead(ctx, d, m)

	return diags
}

func resourceAwsNodePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	cluster_name := d.Get("cluster_name").(string)

	nodepool, err := client.GetNodePool(npName, cluster_name, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS nodepool",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	cluster, err := client.GetCluster(cluster_name, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS nodepool",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("cluster_id", cluster.Meta.UID)

	nodeCount, _ := strconv.Atoi(nodepool.Spec.WorkerNodeCount)

	d.Set("description", nodepool.Meta.Description)
	d.Set("worker_node_count", nodeCount)

	if err := d.Set("cloud_labels", nodepool.Spec.CloudLabels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS nodepool",
			Detail:   fmt.Sprintf("Error getting cloud labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	if err := d.Set("node_labels", nodepool.Spec.NodeLabels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS nodepool",
			Detail:   fmt.Sprintf("Error getting node labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("availability_zone", nodepool.Spec.NodeTkgAws.AvailabilityZone)
	d.Set("instance_type", nodepool.Spec.NodeTkgAws.InstanceType)
	d.Set("version", nodepool.Spec.NodeTkgAws.Version)

	return diags
}

func resourceAwsNodePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*tanzuclient.Client)

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	cloud_labels := d.Get("cloud_labels").(map[string]interface{})
	node_labels := d.Get("node_labels").(map[string]interface{})
	cluster_name := d.Get("cluster_name").(string)
	worker_node_count := d.Get("worker_node_count").(int)

	awsNodeSpec := &tanzuclient.AwsNodeSpec{
		AvailabilityZone: d.Get("availability_zone").(string),
		Version:          d.Get("version").(string),
		InstanceType:     d.Get("instance_type").(string),
	}

	if d.HasChange("worker_node_count") {
		_, err := client.UpdateNodePool(npName, managementClusterName, provisionerName, cluster_name, description, cloud_labels, node_labels, worker_node_count, awsNodeSpec)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to update AWS nodepool",
				Detail:   fmt.Sprintf("Error updating resource %s: %s", d.Get("name"), err),
			})
			return diags
		}

		createStateConf := &resource.StateChangeConf{
			Pending: []string{
				"CREATING",
				"RESIZING",
				"WAITING",
				"UPGRADING",
			},
			Target: []string{
				"READY",
			},
			Refresh: func() (interface{}, string, error) {
				resp, err := client.GetNodePool(npName, cluster_name, managementClusterName, provisionerName)
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

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceAwsNodePoolRead(ctx, d, m)
}

func resourceAwsNodePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	cluster_name := d.Get("cluster_name").(string)

	err := client.DeleteNodePool(npName, cluster_name, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete AWS nodepool",
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
			resp, err := client.DescribeNodePool(npName, cluster_name, managementClusterName, provisionerName)
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
