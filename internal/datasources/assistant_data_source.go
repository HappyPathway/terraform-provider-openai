package datasources

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sashabaranov/go-openai"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &AssistantDataSource{}

func NewAssistantDataSource() datasource.DataSource {
	return &AssistantDataSource{}
}

// AssistantDataSource defines the data source implementation.
type AssistantDataSource struct {
	client *client.Client
}

// AssistantDataSourceModel describes the data source data model.
type AssistantDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	AssistantID  types.String `tfsdk:"assistant_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Model        types.String `tfsdk:"model"`
	Instructions types.String `tfsdk:"instructions"`
	Tools        types.List   `tfsdk:"tools"`
	FileIDs      types.List   `tfsdk:"file_ids"`
	Metadata     types.Map    `tfsdk:"metadata"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
}

// AssistantToolModel represents a tool configuration for an assistant.
type AssistantToolModel struct {
	Type               types.String `tfsdk:"type"`
	FunctionDefinition types.String `tfsdk:"function_definition"`
}

func (d *AssistantDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assistant"
}

func (d *AssistantDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve information about an OpenAI Assistant.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this data source.",
			},
			"assistant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the assistant to retrieve.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the assistant.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the assistant.",
				Computed:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The model used by the assistant.",
				Computed:            true,
			},
			"instructions": schema.StringAttribute{
				MarkdownDescription: "The instructions that the assistant uses to guide its behavior.",
				Computed:            true,
			},
			"tools": schema.ListNestedAttribute{
				MarkdownDescription: "The tools used by the assistant.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of tool. Can be 'code_interpreter', 'retrieval', or 'function'.",
							Computed:            true,
						},
						"function_definition": schema.StringAttribute{
							MarkdownDescription: "JSON string of the function definition when type is 'function'.",
							Computed:            true,
						},
					},
				},
			},
			"file_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A list of file IDs attached to this assistant.",
				Computed:            true,
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Metadata in key-value pairs attached to the assistant.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the assistant was created.",
				Computed:            true,
			},
		},
	}
}

func (d *AssistantDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *AssistantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AssistantDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assistantID := data.AssistantID.ValueString()
	if assistantID == "" {
		resp.Diagnostics.AddError(
			"Missing Assistant ID",
			"The assistant_id attribute is required to retrieve assistant details.",
		)
		return
	}

	tflog.Debug(ctx, "Reading assistant details", map[string]interface{}{
		"assistant_id": assistantID,
	})

	// Retrieve assistant information from API
	assistant, err := d.client.OpenAI.RetrieveAssistant(ctx, assistantID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading OpenAI Assistant",
			fmt.Sprintf("Unable to read assistant %s: %s", assistantID, d.client.HandleError(err)),
		)
		return
	}

	// Map response body to model
	data.ID = types.StringValue(assistant.ID)
	if assistant.Name != nil {
		data.Name = types.StringValue(*assistant.Name)
	} else {
		data.Name = types.StringNull()
	}
	if assistant.Description != nil {
		data.Description = types.StringValue(*assistant.Description)
	} else {
		data.Description = types.StringNull()
	}
	data.Model = types.StringValue(assistant.Model)
	if assistant.Instructions != nil {
		data.Instructions = types.StringValue(*assistant.Instructions)
	} else {
		data.Instructions = types.StringNull()
	}
	data.CreatedAt = types.Int64Value(int64(assistant.CreatedAt))

	// Convert tools
	if len(assistant.Tools) > 0 {
		toolsList, diags := convertOpenAIToolsToTerraform(ctx, assistant.Tools)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tools = toolsList
	}

	// Convert file IDs
	if len(assistant.FileIDs) > 0 {
		fileIDsList, diags := types.ListValueFrom(ctx, types.StringType, assistant.FileIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.FileIDs = fileIDsList
	}

	// Convert metadata
	if assistant.Metadata != nil {
		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, assistant.Metadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Metadata = metadataMap
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Helper function to convert from OpenAI tools to Terraform tools
func convertOpenAIToolsToTerraform(ctx context.Context, tools []openai.AssistantTool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfTools := make([]attr.Value, 0, len(tools))
	for _, tool := range tools {
		toolMap := make(map[string]attr.Value)
		toolMap["type"] = types.StringValue(string(tool.Type))
		if tool.Function != nil {
			functionJSON, err := json.Marshal(tool.Function)
			if err != nil {
				diags.AddError(
					"Error Converting Function Definition",
					fmt.Sprintf("Unable to convert function definition to JSON: %s", err),
				)
				continue
			}
			toolMap["function_definition"] = types.StringValue(string(functionJSON))
		} else {
			toolMap["function_definition"] = types.StringNull()
		}

		toolObj, err := types.ObjectValue(
			map[string]attr.Type{
				"type":                types.StringType,
				"function_definition": types.StringType,
			},
			toolMap,
		)
		if err != nil {
			diags.AddError(
				"Error Converting Tool",
				fmt.Sprintf("Unable to convert tool to Terraform object: %s", err),
			)
			continue
		}
		tfTools = append(tfTools, toolObj)
	}
	toolsList, err := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":                types.StringType,
				"function_definition": types.StringType,
			},
		},
		tfTools,
	)
	if err != nil {
		diags.AddError(
			"Error Creating Tools List",
			fmt.Sprintf("Unable to create tools list: %s", err),
		)
	}
	return toolsList, diags
}
