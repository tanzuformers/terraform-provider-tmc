package tmc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tanzuformers/terraform-provider-tmc/tanzuclient"
)

func resourceTmcAwsDataProtectionCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceTmcAwsDataProtectionCredentialRead,
		CreateContext: resourceTmcAwsDataProtectionCredentialCreate,
		DeleteContext: resourceTmcAwsDataProtectionCredentialDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Aws Data Protection Credential",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Tanzu Aws Data Protection Credential",
			},
			"capability": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Capability of the Tanzu Aws Data Protection Credential",
			},
			"iam_role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM Role Arn of the Tanzu Aws Data Protection Credential",
			},
			"credential_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Credential Provider of the Tanzu Aws Data Protection Credential",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTmcAwsDataProtectionCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	credential, err := client.GetAwsCredential(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(credential.Meta.UID))
	d.Set("capability", credential.Spec.Capability)
	d.Set("iam_role_arn", credential.Spec.Data.AwsCredential.IamRole.Arn)
	d.Set("credential_provider", credential.Spec.MetaData.Provider)
	d.Set("status", credential.Status.Phase)

	return diags
}

func resourceTmcAwsDataProtectionCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	credName := d.Get("name").(string)

	if !IsValidTanzuName(credName) {
		return InvalidTanzuNameError("AWS Data Protection Credential")
	}

	cred := tanzuclient.TmcAwsCredential{
		FullName: &tanzuclient.FullName{
			Name: d.Get("name").(string),
		},
		Spec: &tanzuclient.CredentialSpec{
			MetaData: &tanzuclient.CredentialMetaData{
				Provider: "AWS_EC2", // Always set to this value for AWS Credentials
			},
			Capability: "DATA_PROTECTION", // Always set to this value for AWS Credentials
			Data: &tanzuclient.CredentialData{
				AwsCredential: &tanzuclient.AwsCredential{
					IamRole: &tanzuclient.IamRole{
						Arn: d.Get("iam_role_arn").(string),
					},
				},
			},
		},
	}

	res, err := client.CreateAwsCredential(&cred)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create AWS Data Protection Credential Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(res.Meta.UID)
	d.Set("status", res.Status.Phase)
	d.Set("capability", res.Spec.Capability)
	d.Set("provider", res.Spec.MetaData.Provider)

	resourceTmcAwsDataProtectionCredentialRead(ctx, d, m)

	return nil
}

func resourceTmcAwsDataProtectionCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	credName := d.Get("name").(string)

	err := client.DeleteAwsCredential(credName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete AWS Data Protection Credential Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}
