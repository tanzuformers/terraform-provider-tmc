package tmc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func dataSourceTmcClusterBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcClusterBackupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Cluster Backup",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tanzu Cluster Backup Name",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Cluster",
			},
			"management_cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Management Cluster",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Management Cluster Provisioner",
			},
			"labels": labelsSchemaComputed(),
			"included_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The namespaces to be included for backup from. If empty, all namespaces are included.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"excluded_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The namespaces to be excluded in the backup.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"included_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The name list for the resources to be included into backup. If empty, all resources are included.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"excluded_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The name list for the resources to be excluded in backup.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"label_selector": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Label query over a set of resources to be included in backup.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_labels": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						},
						"match_expressions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"retention_period": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The backup retention period.",
			},
			"storage_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of a BackupStorageLocation where the backup should be stored.",
			},
			"snapshot_volumes": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A flag which specifies whether to take cloud snapshots of any PV's referenced in the set of objects included in the Backup.",
			},
			"volume_snapshot_locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "A list containing names of VolumeSnapshotLocations associated with this backup.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"include_cluster_scoped_resources": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A flag which specifies whether cluster-scoped resources should be included for consideration in the backup",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTmcClusterBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	backup, err := client.GetClusterBackup(d.Get("name").(string), d.Get("management_cluster_name").(string), d.Get("cluster_name").(string), d.Get("provisioner_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(backup.Meta.UID))
	d.Set("included_namespaces", backup.Spec.IncludedNamespaces)
	d.Set("excluded_namespaces", backup.Spec.ExcludedNamespaces)
	d.Set("included_resources", backup.Spec.IncludedResources)
	d.Set("excluded_resources", backup.Spec.ExcludedResources)
	d.Set("retention_period", backup.Spec.TTL)
	if err := d.Set("label_selector", flattenLabelSelector(backup.Spec.LabelSelector)); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse backup resource data",
			Detail:   fmt.Sprintf("Error parsing label selectors for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.Set("storage_location", backup.Spec.StorageLocation)
	d.Set("snapshot_volumes", backup.Spec.SnapshotVolumes)
	d.Set("include_cluster_scoped_resources", backup.Spec.IncludeClusterResources)
	d.Set("volume_snapshot_locations", backup.Spec.VolumeSnapshotLocations)
	if err := d.Set("labels", backup.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse backup resource data",
			Detail:   fmt.Sprintf("Error setting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.Set("status", backup.Status.Phase)

	return diags
}

func flattenLabelSelector(labelSelector *tanzuclient.LabelSelector) []interface{} {
	ls := make(map[string]interface{})

	if labelSelector != nil {
		if labelSelector.MatchLabels != nil {
			ls["match_labels"] = labelSelector.MatchLabels

		}
		ls["match_expressions"] = flattenMatchExpressions(labelSelector.MatchExpressions)
	}

	return []interface{}{ls}
}

func flattenMatchExpressions(matchExpressions []tanzuclient.MatchExpressions) []interface{} {
	if matchExpressions != nil {
		mes := make([]interface{}, len(matchExpressions))

		for i, matchExpression := range matchExpressions {
			me := make(map[string]interface{})

			me["key"] = matchExpression.Key
			me["operator"] = matchExpression.Operator
			me["values"] = matchExpression.Values

			mes[i] = me
		}

		return mes
	}

	return make([]interface{}, 0)
}
