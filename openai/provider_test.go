package openai

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var providerFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"openai": testAccProvider,
	}
	providerFactories = map[string]func() (*schema.Provider, error){
		"openai": func() (*schema.Provider, error) {
			return Provider(), nil
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
	if v := os.Getenv("OPENAI_API_KEY"); v == "" {
		t.Fatal("OPENAI_API_KEY must be set for acceptance tests")
	}
}

func TestProvider_configure(t *testing.T) {
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
			configFn, diags := tc.P.ConfigureContextFunc(context.Background(), &schema.ResourceData{})
			if diags.HasError() && !tc.Err {
				t.Fatalf("unexpected error: %v", diags)
			}
			if tc.Err {
				if !diags.HasError() {
					t.Fatalf("expected error")
				}
				return
			}

			d := &schema.ResourceData{}
			for k, v := range tc.Config {
				if err := d.Set(k, v); err != nil {
					t.Fatalf("err: %s", err)
				}
			}

			// Call the configuration function directly
			configFunc := configFn.(func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics))
			_, diags = configFunc(context.Background(), d)

			if tc.Err {
				if !diags.HasError() {
					t.Fatalf("expected error")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
		})
	}
}
