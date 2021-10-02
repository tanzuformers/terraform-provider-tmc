package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceTmcWorkspaceRead,
		CreateContext: resourceTmcWorkspaceCreate,
		UpdateContext: resourceTmcWorkspaceUpdate,
		DeleteContext: resourceTmcWorkspaceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Workspace",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Workspace",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Tanzu Workspace",
			},
			"labels": labelsSchema(),
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceTmcWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workspace, err := client.GetWorkspace(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", workspace.Meta.Description)
	if err := d.Set("labels", workspace.Meta.SimpleMetaData.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to perform specified operation",
			Detail:   fmt.Sprintf("Error setting tags for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.SetId(string(workspace.Meta.SimpleMetaData.UID))

	return diags
}

func resourceTmcWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	workspace, err := client.CreateWorkspace(d.Get("name").(string), d.Get("description").(string), d.Get("labels").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(workspace.Meta.SimpleMetaData.UID))

	return resourceTmcWorkspaceRead(ctx, d, meta)
}

func resourceTmcWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workspaceName := d.Get("name").(string)

	if d.HasChange("description") || d.HasChange("labels") {
		description := d.Get("description").(string)
		labels := d.Get("labels").(map[string]interface{})

		_, err := client.UpdateWorkspace(workspaceName, description, labels)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to update workspace",
				Detail:   fmt.Sprintf("Cannot update the workspace %s with the new values: %s", d.Get("name").(string), err),
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceTmcWorkspaceRead(ctx, d, meta)
}

func resourceTmcWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := client.DeleteWorkspace(d.Get("name").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete workspace",
			Detail:   fmt.Sprintf("Cannot delete given workspace %s: %s", d.Get("name").(string), err),
		})
		return diags
	}

	d.SetId("")

	return diags
}
