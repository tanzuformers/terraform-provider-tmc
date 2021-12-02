package tmc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func dataSourceTmcAwsStorageCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcAwsStorageCredentialRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Aws Storage Credential",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Aws Storage Credential",
			},
			"capability": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Capability of the Tanzu Aws Storage Credential",
			},
			"credential_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Credential Provider of the Tanzu Aws Storage Credential",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTmcAwsStorageCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	credential, err := client.GetAwsCredential(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(credential.Meta.UID))
	d.Set("status", credential.Status.Phase)
	d.Set("capability", credential.Spec.Capability)
	d.Set("credential_provider", credential.Spec.MetaData.Provider)

	return diags
}
