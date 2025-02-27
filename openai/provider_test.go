package openai

import (
	"context"
	"os"
	"testing"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var providerFactories map[string]func() (*schema.Provider, error)
var useMockClient bool

func init() {
	useMockClient = os.Getenv("OPENAI_USE_MOCK") != ""

	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"openai": testAccProvider,
	}

	providerFactories = map[string]func() (*schema.Provider, error){
		"openai": func() (*schema.Provider, error) {
			p := Provider()
			if useMockClient {
				p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
					return testutil.NewMockClient(), nil
				}
			}
			return p, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if !useMockClient {
		if v := os.Getenv("OPENAI_API_KEY"); v == "" {
			t.Fatal("OPENAI_API_KEY must be set for acceptance tests when not using mock client")
		}
	}
}

func TestProvider_configure(t *testing.T) {
	// Use mock client to avoid real API calls in provider configuration tests
	useMockClient = true
	defer func() { useMockClient = os.Getenv("OPENAI_USE_MOCK") != "" }()

	cases := map[string]struct {
		P      *schema.Provider
		Config map[string]interface{}
		Err    bool
	}{
		"valid minimal config": {
			P: Provider(),
			Config: map[string]interface{}{
				"api_key": "test-api-key",
			},
			Err: false,
		},
		"missing api_key": {
			P:      Provider(),
			Config: map[string]interface{}{},
			Err:    true,
		},
		"valid full config": {
			P: Provider(),
			Config: map[string]interface{}{
				"api_key":         "test-api-key",
				"organization_id": "test-org",
				"retry_max":       5,
				"retry_delay":     2,
				"timeout":         60,
			},
			Err: false,
		},
		"invalid retry_max": {
			P: Provider(),
			Config: map[string]interface{}{
				"api_key":   "test-api-key",
				"retry_max": -1,
			},
			Err: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, tc.P.Schema, tc.Config)
			_, diags := tc.P.ConfigureContextFunc(context.Background(), d)

			if diags.HasError() != tc.Err {
				t.Fatalf("expected error: %t, got error: %t - errors: %v", tc.Err, diags.HasError(), diags)
			}
		})
	}
}
