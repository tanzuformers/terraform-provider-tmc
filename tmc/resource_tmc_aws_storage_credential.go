package tmc

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcAwsStorageCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceTmcAwsStorageCredentialRead,
		CreateContext: resourceTmcAwsStorageCredentialCreate,
		DeleteContext: resourceTmcAwsStorageCredentialDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Aws Storage Credential",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Tanzu Aws Storage Credential",
			},
			"capability": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Capability of the Tanzu Aws Storage Credential",
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Access Key ID of the Tanzu Aws Storage Credential",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Secret Access Key of the Tanzu Aws Storage Credential",
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

func resourceTmcAwsStorageCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	credential, err := client.GetAwsCredential(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(credential.Meta.UID))
	d.Set("capability", credential.Spec.Capability)
	d.Set("credential_provider", credential.Spec.MetaData.Provider)
	d.Set("status", credential.Status.Phase)

	return diags
}

func resourceTmcAwsStorageCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	credName := d.Get("name").(string)

	if !IsValidTanzuName(credName) {
		return InvalidTanzuNameError("AWS Storage Credential")
	}

	cred := tanzuclient.TmcAwsCredential{
		FullName: &tanzuclient.FullName{
			Name: d.Get("name").(string),
		},
		Spec: &tanzuclient.CredentialSpec{
			MetaData: &tanzuclient.CredentialMetaData{
				Provider: "GENERIC_S3", // Always set to this value for AWS Credentials
			},
			Capability: "DATA_PROTECTION", // Always set to this value for AWS Credentials
			Data: &tanzuclient.CredentialData{
				KeyValue: &tanzuclient.AwsCredentialKey{
					Type: "OPAQUE_SECRET_TYPE", // Always set to this value for AWS Access Keys
					Data: &tanzuclient.AwsAccessKey{
						AccessKeyId:     base64.StdEncoding.EncodeToString([]byte(d.Get("access_key_id").(string))),
						SecretAccessKey: base64.StdEncoding.EncodeToString([]byte(d.Get("secret_access_key").(string))),
					},
				},
			},
		},
	}

	res, err := client.CreateAwsCredential(&cred)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create AWS Storage Credential Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(res.Meta.UID)
	d.Set("status", res.Status.Phase)
	d.Set("capability", res.Spec.Capability)
	d.Set("provider", res.Spec.MetaData.Provider)

	resourceTmcAwsStorageCredentialRead(ctx, d, m)

	return nil
}

func resourceTmcAwsStorageCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	credName := d.Get("name").(string)

	err := client.DeleteAwsCredential(credName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete AWS Storage Credential Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}
