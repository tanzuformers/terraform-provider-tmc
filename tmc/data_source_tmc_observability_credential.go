package tmc

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func dataSourceTmcObservabilityCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcObservabilityCredentialRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Observability Credential",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Observability Credential",
			},
			"capability": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Capability of the Tanzu Observability Credential",
			},
			"observability_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the Tanzu Observability Instance",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTmcObservabilityCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	credential, err := client.GetObservabilityCredential(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(credential.Meta.UID)
	d.Set("status", credential.Status.Phase)
	d.Set("capability", credential.Spec.Capability)

	if url, err := flattenObservabilityCredentialAnnotations(credential.Meta.Annotations); err != nil {
		return diag.FromErr(err)
	} else {
		d.Set("observability_url", url)
	}

	return diags
}

func flattenObservabilityCredentialAnnotations(annotations map[string]string) (*string, error) {
	if url, ok := annotations["wavefront.url"]; ok {
		return &url, nil
	}
	return nil, errors.New("key does not exist")
}
