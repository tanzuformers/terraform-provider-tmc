package tmc

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcClusterBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceTmcClusterBackupRead,
		CreateContext: resourceTmcClusterBackupCreate,
		DeleteContext: resourceTmcClusterBackupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Cluster Backup",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique Name of the Tanzu Cluster Backup",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`).MatchString(v) {
						errs = append(errs, fmt.Errorf("invalid resource name: name must start and end with a letter or number, and can contain only lowercase letters, numbers, and hyphens"))
					}
					return
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Tanzu Cluster",
			},
			"management_cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Tanzu Management Cluster",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Tanzu Management Cluster Provisioner",
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"included_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The namespaces to be included for backup from. If empty, all namespaces are included.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"excluded_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The namespaces to be excluded in the backup.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"included_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The name list for the resources to be included into backup. If empty, all resources are included.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"excluded_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The name list for the resources to be excluded in backup.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"label_selector": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Label query over a set of resources to be included in backup.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_labels": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
						},
						"match_expressions": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
										ForceNew: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
										ForceNew: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										ForceNew: true,
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
				Required:    true,
				ForceNew:    true,
				Description: "The backup retention period.",
			},
			"storage_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of a BackupStorageLocation where the backup should be stored.",
			},
			"snapshot_volumes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "A flag which specifies whether to take cloud snapshots of any PV's referenced in the set of objects included in the Backup.",
			},
			"volume_snapshot_locations": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"snapshot_volumes", "volume_snapshot_locations"},
				Description:  "A list containing names of VolumeSnapshotLocations associated with this backup.",
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"include_cluster_scoped_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "A flag which specifies whether cluster-scoped resources should be included for consideration in the backup",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTmcClusterBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("snapshot_volumes", backup.Spec.SnapshotVolumes)
	d.Set("storage_location", backup.Spec.StorageLocation)
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

func resourceTmcClusterBackupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	backupName := d.Get("name").(string)

	if !IsValidTanzuName(backupName) {
		return InvalidTanzuNameError("Cluster Backup")
	}

	labelSelector, err := expandLabelSelector(d.Get("label_selector").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	backup := tanzuclient.TmcClusterBackup{
		FullName: &tanzuclient.FullName{
			Name:                  d.Get("name").(string),
			ManagementClusterName: d.Get("management_cluster_name").(string),
			ProvisionerName:       d.Get("provisioner_name").(string),
			ClusterName:           d.Get("cluster_name").(string),
		},
		Meta: &tanzuclient.MetaData{
			Labels: d.Get("labels").(map[string]interface{}),
		},
		Spec: &tanzuclient.ClusterBackupSpec{
			LabelSelector:   labelSelector,
			SnapshotVolumes: d.Get("snapshot_volumes").(bool),
			TTL:             d.Get("retention_period").(string),

			StorageLocation: d.Get("storage_location").(string),
		},
	}

	if v, ok := d.Get("included_namespaces").([]string); ok && len(v) > 0 {
		backup.Spec.IncludedNamespaces = v
	}
	if v, ok := d.Get("excluded_namespaces").([]string); ok && len(v) > 0 {
		backup.Spec.ExcludedNamespaces = v
	}
	if v, ok := d.Get("included_resources").([]string); ok && len(v) > 0 {
		backup.Spec.IncludedResources = v
	}
	if v, ok := d.Get("excluded_resources").([]string); ok && len(v) > 0 {
		backup.Spec.ExcludedResources = v
	}
	if v, ok := d.Get("include_cluster_scoped_resources").(bool); ok {
		backup.Spec.IncludeClusterResources = v
	}
	if v, ok := d.Get("volume_snapshot_locations").([]string); ok && len(v) > 0 {
		backup.Spec.VolumeSnapshotLocations = v
	}

	res, err := client.CreateClusterBackup(d.Get("cluster_name").(string), &backup)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Creating Cluster backup Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(res.Meta.UID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"PENDING",
			"INPROGRESS",
		},
		Target: []string{
			"COMPLETED",
		},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.GetClusterBackup(backupName, d.Get("management_cluster_name").(string), d.Get("cluster_name").(string), d.Get("provisioner_name").(string))
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status.Phase, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Creating Cluster Backup Failed",
			Detail:   fmt.Sprintf("Error waiting for Cluster Backup (%s) to be created: %s", d.Get("name").(string), err),
		})
		return diags
	}

	return resourceTmcClusterBackupRead(ctx, d, m)
}

func resourceTmcClusterBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*tanzuclient.Client)

	err := client.DeleteClusterBackup(d.Get("name").(string), d.Get("management_cluster_name").(string), d.Get("cluster_name").(string), d.Get("provisioner_name").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete Backup Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}

func expandLabelSelector(in []interface{}) (*tanzuclient.LabelSelector, error) {
	ls := &tanzuclient.LabelSelector{}

	if len(in) < 1 {
		return ls, nil
	}
	l := in[0].(map[string]interface{})

	if v, ok := l["match_labels"].(map[string]interface{}); ok && len(v) > 0 {
		ls.MatchLabels = expandStringMap(l["match_labels"].(map[string]interface{}))
	}

	if v, ok := l["match_expressions"].([]interface{}); ok && len(v) > 0 {
		exp, err := expandMatchExpressions(v)
		if err != nil {
			return ls, err
		}
		ls.MatchExpressions = exp
	}

	return ls, nil
}

func expandStringMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func expandMatchExpressions(exprs []interface{}) ([]tanzuclient.MatchExpressions, error) {
	if len(exprs) == 0 {
		return []tanzuclient.MatchExpressions{}, nil
	}
	ex := make([]tanzuclient.MatchExpressions, len(exprs))
	for i, c := range exprs {
		expr := c.(map[string]interface{})
		if key, ok := expr["key"]; ok {
			ex[i].Key = key.(string)
		}
		if operator, ok := expr["operator"]; ok {
			ex[i].Operator = operator.(string)
		}
		if values, ok := expr["values"]; ok {
			ex[i].Values = values.([]string)
		}
	}
	return ex, nil
}
