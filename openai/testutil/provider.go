package testutil

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderFactories returns a map of provider factories configured with the mock client
func ProviderFactories(p *schema.Provider) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"openai": func() (*schema.Provider, error) {
			if os.Getenv("OPENAI_MOCK") != "" {
				// Configure provider with mock client
				return ConfigureProviderWithMockClient(p), nil
			}
			return p, nil
		},
	}
}

// ConfigureProviderWithMockClient returns a provider configured to use the mock client
func ConfigureProviderWithMockClient(p *schema.Provider) *schema.Provider {
	// Store original ConfigureContextFunc
	originalConfigureContextFunc := p.ConfigureContextFunc

	// Override ConfigureContextFunc to inject mock client
	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Call original ConfigureContextFunc first if it exists
		if originalConfigureContextFunc != nil {
			_, diags := originalConfigureContextFunc(ctx, d)
			if diags.HasError() {
				return nil, diags
			}
		}

		// Return mock client
		return NewMockClient(), nil
	}

	// Clear ConfigureFunc to avoid conflict
	p.ConfigureFunc = nil

	return p
}
