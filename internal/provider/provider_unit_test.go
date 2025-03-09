package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Provider_Impl is a test helper for ensuring the provider implements the expected interfaces.
func TestProvider_Impl(t *testing.T) {
	var _ provider.Provider = &OpenAIProvider{}
}

func TestProviderSchema(t *testing.T) {
	// Create provider instance
	p := &OpenAIProvider{}

	// Get schema
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, resp)

	// Verify schema
	if len(resp.Schema.Attributes) == 0 {
		t.Error("expected schema to have attributes")
	}
}

func TestProviderConfigure(t *testing.T) {
	testCases := map[string]struct {
		config        map[string]any
		expectedToken string
		expectedError bool
	}{
		"valid-config": {
			config: map[string]any{
				"api_token": "test-token",
			},
			expectedToken: "test-token",
			expectedError: false,
		},
		"missing-token": {
			config:        map[string]any{},
			expectedError: true,
		},
		"with-base-url": {
			config: map[string]any{
				"api_token": "test-token",
				"base_url":  "https://api.example.com",
			},
			expectedToken: "test-token",
			expectedError: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			p := &OpenAIProvider{}

			// Get schema for configuration
			schemaResp := &provider.SchemaResponse{}
			p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			// Create configuration
			var data tfsdk.Config
			err := tfsdk.ValueFrom(context.Background(), tc.config, schemaResp.Schema.Type(), &data)
			if err != nil {
				t.Fatalf("error creating config: %v", err)
			}

			// Configure provider
			resp := &provider.ConfigureResponse{}
			p.Configure(context.Background(), provider.ConfigureRequest{
				Config: data,
			}, resp)

			if tc.expectedError {
				if !resp.Diagnostics.HasError() {
					t.Error("expected error")
				}
				return
			}

			if resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}

			// Validate client configuration
			if p.client == nil {
				t.Fatal("client is nil")
			}
		})
	}
}

func TestProviderMetadata(t *testing.T) {
	t.Parallel()
	p := &OpenAIProvider{
		version: "test",
	}
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), provider.MetadataRequest{}, resp)
	if resp.TypeName != "openai" {
		t.Errorf("expected type name to be openai, got: %s", resp.TypeName)
	}
	if resp.Version != "test" {
		t.Errorf("expected version to be test, got: %s", resp.Version)
	}
}
