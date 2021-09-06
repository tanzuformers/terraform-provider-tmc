package tmc

import (
	"context"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {

	// The actual provider
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"TMC_API_TOKEN",
				}, nil),
				Description: descriptions["api_token"],
			},
			"org_url": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"TMC_ORG_URL",
				}, nil),
				Description: descriptions["org_url"],
			},
		},

		// List of Data sources supported by the provider
		DataSourcesMap: map[string]*schema.Resource{
			"tmc_workspace":      dataSourceTmcWorkspace(),
			"tmc_workspaces":     dataSourceTmcWorkspaces(),
			"tmc_cluster_group":  dataSourceClusterGroup(),
			"tmc_cluster_groups": dataSourceClusterGroups(),
			"tmc_provisioners":   dataSourceTmcProvisioners(),
			"tmc_provisioner":    dataSourceTmcProvisioner(),
		},

		// List of Resources supported by the provider
		ResourcesMap: map[string]*schema.Resource{
			"tmc_workspace":     resourceTmcWorkspace(),
			"tmc_cluster_group": resourceClusterGroup(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, provider)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_token": "API_TOKEN generated by the VMware Cloud Services Console. If not set,\n" +
			"defaults to the environment variable TMC_API_TOKEN",
		"org_url": "VMware Cloud Console Service URL unique to your organization. If not set,\n" +
			"defaults to the environment variable TMC_ORG_URL",
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	apiToken := d.Get("api_token").(string)
	orgURL := d.Get("org_url").(string)
	var err error

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (apiToken != "") && (orgURL != "") {
		client, err := tanzuclient.NewClient(&orgURL, &apiToken)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}

	client, err := tanzuclient.NewClient(nil, nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &client, diags
}
