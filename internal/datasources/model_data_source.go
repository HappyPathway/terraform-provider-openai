package datasources

import (
	"context"
	"fmt"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ModelDataSource{}
	_ datasource.DataSourceWithConfigure = &ModelDataSource{}
)

// NewModelDataSource is a helper function to simplify the provider implementation.
func NewModelDataSource() datasource.DataSource {
	return &ModelDataSource{}
}

// ModelDataSource is the data source implementation.
type ModelDataSource struct {
	client *client.Client
}

// ModelDataSourceModel maps the data source schema data.
type ModelDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	ModelID     types.String `tfsdk:"model_id"`
	Created     types.Int64  `tfsdk:"created"`
	OwnedBy     types.String `tfsdk:"owned_by"`
	Object      types.String `tfsdk:"object"`
	FilterOwner types.String `tfsdk:"filter_owner"`
}

// Metadata returns the data source type name.
func (d *ModelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

// Schema defines the schema for the data source.
func (d *ModelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an OpenAI model.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for this data source. Matches the model_id.",
				Computed:            true,
			},
			"model_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the model to retrieve (e.g., 'gpt-4', 'gpt-3.5-turbo', 'text-embedding-ada-002').",
				Required:            true,
			},
			"created": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) when the model was created.",
				Computed:            true,
			},
			"owned_by": schema.StringAttribute{
				MarkdownDescription: "The organization that owns the model.",
				Computed:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object type, which is always 'model'.",
				Computed:            true,
			},
			"filter_owner": schema.StringAttribute{
				MarkdownDescription: "Filter models by owner (e.g., 'openai', 'user'). Optional.",
				Optional:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ModelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ModelDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	modelID := state.ModelID.ValueString()
	filterOwner := state.FilterOwner.ValueString()

	tflog.Info(ctx, "Reading OpenAI Model", map[string]interface{}{
		"model_id": modelID,
	})

	// First check if we can get the specific model by ID
	model, err := d.client.OpenAI.GetModel(ctx, modelID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading OpenAI Model",
			fmt.Sprintf("Unable to read model %s: %s", modelID, d.client.HandleError(err)),
		)
		return
	}

	// If filter_owner is specified, verify that the model belongs to the specified owner
	if filterOwner != "" && model.OwnedBy != filterOwner {
		resp.Diagnostics.AddError(
			"OpenAI Model Ownership Mismatch",
			fmt.Sprintf("Model %s is owned by '%s', but filter specified '%s'",
				modelID, model.OwnedBy, filterOwner),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(model.ID)
	state.ModelID = types.StringValue(model.ID)
	state.Created = types.Int64Value(int64(model.CreatedAt))
	state.OwnedBy = types.StringValue(model.OwnedBy)
	state.Object = types.StringValue(model.Object)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *ModelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}
