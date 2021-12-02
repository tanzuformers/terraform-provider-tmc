package tmc

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcObservabilityCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceTmcObservabilityCredentialRead,
		CreateContext: resourceTmcObservabilityCredentialCreate,
		DeleteContext: resourceTmcObservabilityCredentialDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Observability Credential",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Tanzu Observability Credential",
			},
			"capability": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Capability of the Tanzu Observability Credential",
			},
			"observability_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "URL of the Tanzu Observability Credential",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "API Token of the Tanzu Observability Credential",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTmcObservabilityCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	credential, err := client.GetAwsCredential(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(credential.Meta.UID))
	d.Set("capability", credential.Spec.Capability)
	d.Set("status", credential.Status.Phase)

	return diags
}

func resourceTmcObservabilityCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	credName := d.Get("name").(string)

	if !IsValidTanzuName(credName) {
		return InvalidTanzuNameError("Observability Credential")
	}

	cred := tanzuclient.TmcObservabilityCredential{
		FullName: &tanzuclient.FullName{
			Name: d.Get("name").(string),
		},
		Meta: &tanzuclient.MetaData{
			Annotations: map[string]string{
				"wavefront.url": d.Get("observability_url").(string),
			},
		},
		Spec: &tanzuclient.ObservabilityCredentialSpec{
			Capability: "TANZU_OBSERVABILITY", // Always set to this value for Observability Credentials
			Data: tanzuclient.ObservabilityCredentialData{
				KeyValue: tanzuclient.ObservabilityKey{
					Data: tanzuclient.WaveFrontData{
						Token: base64.StdEncoding.EncodeToString([]byte(d.Get("api_token").(string))),
					},
				},
			},
		},
	}

	res, err := client.CreateObservabilityCredential(&cred)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create Observability Credential Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(res.Meta.UID)
	d.Set("status", res.Status.Phase)
	d.Set("capability", res.Spec.Capability)

	resourceTmcObservabilityCredentialRead(ctx, d, m)

	return nil
}

func resourceTmcObservabilityCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	credName := d.Get("name").(string)

	err := client.DeleteObservabilityCredential(credName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete Observability Credential Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}
