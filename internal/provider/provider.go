package provider

import (
	"context"
	"os"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/darnold/terraform-provider-openai/internal/datasources"
	"github.com/darnold/terraform-provider-openai/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &OpenAIProvider{}
)

// OpenAIProvider is the provider implementation.
type OpenAIProvider struct {
	// version is set to the provider version on release.
	version string
	client  *client.Client
}

// OpenAIProviderModel describes the provider data model.
type OpenAIProviderModel struct {
	APIKey             types.String `tfsdk:"api_key"`
	Organization       types.String `tfsdk:"organization"`
	BaseURL            types.String `tfsdk:"base_url"`
	EnableDebugLogging types.Bool   `tfsdk:"enable_debug_logging"`
}

// New creates a new provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenAIProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *OpenAIProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openai"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *OpenAIProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The OpenAI provider provides resources to interact with the OpenAI API.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "OpenAI API Key. Can also be specified with the `OPENAI_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "OpenAI Organization ID. Can also be specified with the `OPENAI_ORGANIZATION` environment variable.",
				Optional:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "OpenAI Base URL. Can also be specified with the `OPENAI_BASE_URL` environment variable.",
				Optional:            true,
			},
			"enable_debug_logging": schema.BoolAttribute{
				MarkdownDescription: "Enable debug logging. Defaults to false.",
				Optional:            true,
			},
		},
	}
}

// Configure prepares a OpenAI API client for data sources and resources.
func (p *OpenAIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config OpenAIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as API key",
		)
		return
	}

	if config.APIKey.IsNull() {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			resp.Diagnostics.AddError(
				"Missing API Key Configuration",
				"While configuring the provider, the API key was not found. "+
					"Either set the api_key argument in the provider configuration, "+
					"or set the OPENAI_API_KEY environment variable.",
			)
			return
		}
		config.APIKey = types.StringValue(apiKey)
	}

	// Initialize client configuration
	clientConfig := client.Config{
		APIKey:       config.APIKey.ValueString(),
		BaseURL:      config.BaseURL.ValueString(),
		Organization: config.Organization.ValueString(),
	}

	// Create new OpenAI client
	c, err := client.NewClient(ctx, clientConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OpenAI API Client",
			"An unexpected error occurred when creating the OpenAI API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OpenAI Client Error: "+err.Error(),
		)
		return
	}

	p.client = c
	resp.DataSourceData = c
	resp.ResourceData = c
}

// DataSources defines the data sources implemented in the provider.
func (p *OpenAIProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewModelDataSource,
		datasources.NewAssistantDataSource,
		datasources.NewChatCompletionDataSource,
		datasources.NewVectorStoreDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *OpenAIProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewAssistantResource,
		resources.NewChatCompletionResource,
		resources.NewEmbeddingResource,
		resources.NewFileResource,
		resources.NewFineTuneResource,
		resources.NewMessageResource,
		resources.NewRunResource,
		resources.NewThreadResource,
		resources.NewVectorStoreResource,
		resources.NewVectorStoreFileResource,
	}
}
