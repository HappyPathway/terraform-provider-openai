package testutil

import (
	"os"

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
	// Store original ConfigureFunc
	originalConfigureFunc := p.ConfigureFunc

	// Override ConfigureFunc to inject mock client
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		// Call original ConfigureFunc first
		if originalConfigureFunc != nil {
			_, err := originalConfigureFunc(d)
			if err != nil {
				return nil, err
			}
		}
		
		// Return mock client
		return NewMockClient(), nil
	}

	return p
}