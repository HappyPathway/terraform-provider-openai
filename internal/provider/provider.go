package provider

import (
	"context"
	"os"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/darnold/terraform-provider-openai/internal/datasources"
	"github.com/darnold/terraform-provider-openai/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &OpenAIProvider{}
)

// OpenAIProvider is the provider implementation.
type OpenAIProvider struct {
	// version is set to the provider version on release.
	version string
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
	tflog.Info(ctx, "Configuring OpenAI client")

	var config OpenAIProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override with Terraform configuration if set
	apiKey := os.Getenv("OPENAI_API_KEY")
	organization := os.Getenv("OPENAI_ORGANIZATION")
	baseURL := os.Getenv("OPENAI_BASE_URL")
	enableDebug := false

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if !config.Organization.IsNull() {
		organization = config.Organization.ValueString()
	}

	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	if !config.EnableDebugLogging.IsNull() {
		enableDebug = config.EnableDebugLogging.ValueBool()
	}

	// Validate required configuration
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing OpenAI API Key",
			"The provider cannot create the OpenAI API client as there is a missing or empty value for the OpenAI API Key. "+
				"Set the api_key value in the configuration or use the OPENAI_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	// Create a new OpenAI client using the configuration values
	openaiClient, err := client.NewClient(apiKey, organization, baseURL, enableDebug)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OpenAI API Client",
			"An error occurred when creating the OpenAI API client: "+
				err.Error(),
		)
		return
	}

	if enableDebug {
		tflog.Debug(ctx, "Enabling debug logging for OpenAI provider")
	}

	// Make the client available during DataSource and Resource Configure methods
	resp.DataSourceData = openaiClient
	resp.ResourceData = openaiClient

	tflog.Info(ctx, "Configured OpenAI client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *OpenAIProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewModelDataSource,
		datasources.NewAssistantDataSource,
		datasources.NewChatCompletionDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *OpenAIProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewChatCompletionResource,
		resources.NewEmbeddingResource,
		resources.NewFileResource,
		resources.NewFineTuneResource,
		resources.NewAssistantResource,
		resources.NewThreadResource,
		resources.NewMessageResource,
	}
}
