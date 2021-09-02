const (
	// Provider name for single configuration testing
	ProviderNameAws = "tmc"
)

// testAccProviderFactories is a static map containing only the main provider instance
//
// Use other testAccProviderFactories functions, such as testAccProviderFactoriesAlternate,
// for tests requiring special provider configurations.
var testAccProviderFactories map[string]func() (*schema.Provider, error)

// testAccProvider is the "main" provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// testAccPreCheck(t) must be called before using this provider instance.
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()

	testAccProviders = map[string]*schema.Provider{
		ProviderNameTmc: testAccProvider,
	}

	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderNameTmc: func() (*schema.Provider, error) { return Provider(), nil },
	}
}

func TestProvider(t *testing.T) {
	provider := Provider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ schema.Provider = *Provider()
}

// testAccPreCheck verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to be validated to ensure necessary provider attributes are present beforehand.
func testAccPreCheck(t *testing.T) {
	ctx := context.TODO()

	hasTmcCredentials = os.Getenv("TMC_ORG_URL") != "" && os.Getenv("TMC_API_TOKEN") != ""

	if !hasUserCredentials {
		t.Fatalf("Credentials for TMC must be configured as the environment variables" +
			"TMC_ORG_URL and TMC_API_TOKEN for running acceptance tests")
	}

	diags := testAccProvider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if diags.HasError() {
		t.Fatal(diags[0].Summary)
	}
	return
}
