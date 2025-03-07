package resources

import (
	"context"
	"fmt"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sashabaranov/go-openai"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &AssistantResource{}
var _ resource.ResourceWithImportState = &AssistantResource{}

func NewAssistantResource() resource.Resource {
	return &AssistantResource{}
}

// AssistantResource defines the resource implementation.
type AssistantResource struct {
	client *client.Client
}

// AssistantResourceModel describes the resource data model.
type AssistantResourceModel struct {
	ID            types.String                 `tfsdk:"id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Model         types.String                 `tfsdk:"model"`
	Instructions  types.String                 `tfsdk:"instructions"`
	Tools         types.List                   `tfsdk:"tools"`
	ToolResources *AssistantToolResourcesModel `tfsdk:"tool_resources"`
	Metadata      types.Map                    `tfsdk:"metadata"`
	CreatedAt     types.Int64                  `tfsdk:"created_at"`
}

type AssistantToolModel struct {
	Type     types.String `tfsdk:"type"`
	Function types.String `tfsdk:"function"`
}

type AssistantToolResourcesModel struct {
	CodeInterpreter *AssistantToolResourcesCodeInterpreterModel `tfsdk:"code_interpreter"`
	FileSearch      *AssistantToolResourcesFileSearchModel      `tfsdk:"file_search"`
}

type AssistantToolResourcesCodeInterpreterModel struct {
	FileIDs []string `tfsdk:"file_ids"`
}

type AssistantToolResourcesFileSearchModel struct {
	VectorStoreIDs []string `tfsdk:"vector_store_ids"`
}

func (r *AssistantResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assistant"
}

func (r *AssistantResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages an OpenAI Assistant, which can use various tools and capabilities to help with tasks.",

		Blocks: map[string]schema.Block{
			"tool_resources": schema.SingleNestedBlock{
				MarkdownDescription: "Resources made available to the assistant's tools.",
				Blocks: map[string]schema.Block{
					"code_interpreter": schema.SingleNestedBlock{
						MarkdownDescription: "Resources for the code interpreter tool.",
						Attributes: map[string]schema.Attribute{
							"file_ids": schema.SetAttribute{
								MarkdownDescription: "File IDs that the code interpreter can use.",
								ElementType:         types.StringType,
								Optional:            true,
							},
						},
					},
					"file_search": schema.SingleNestedBlock{
						MarkdownDescription: "Resources for the file search tool.",
						Attributes: map[string]schema.Attribute{
							"vector_store_ids": schema.ListAttribute{
								MarkdownDescription: "Vector store IDs for the file search tool.",
								ElementType:         types.StringType,
								Optional:            true,
							},
						},
					},
				},
			},
		},

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the assistant.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the assistant.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the assistant.",
				Optional:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "ID of the model to use for the assistant.",
				Required:            true,
			},
			"instructions": schema.StringAttribute{
				MarkdownDescription: "The system instructions that the assistant uses for tasks.",
				Optional:            true,
			},
			"tools": schema.ListAttribute{
				MarkdownDescription: "A list of tools enabled for the assistant. Valid values are: code_interpreter, file_search, and function.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Metadata in key-value pairs for the assistant.",
				Optional:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the assistant was created.",
				Computed:            true,
			},
		},
	}
}

func (r *AssistantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *AssistantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AssistantResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the assistant request
	assistantReq := openai.AssistantRequest{
		Model: plan.Model.ValueString(),
	}

	// Set optional parameters
	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		assistantReq.Name = &name
	}

	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		assistantReq.Description = &description
	}

	if !plan.Instructions.IsNull() {
		instructions := plan.Instructions.ValueString()
		assistantReq.Instructions = &instructions
	}

	// Process tools if provided
	if !plan.Tools.IsNull() {
		var toolStrings []string
		diags = plan.Tools.ElementsAs(ctx, &toolStrings, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tools, toolDiags := convertToolsToOpenAI(ctx, toolStrings)
		resp.Diagnostics.Append(toolDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		assistantReq.Tools = tools
	}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		assistantReq.Metadata = convertToMapStringAny(metadata)
	}

	// Handle tool_resources if provided
	if plan.ToolResources != nil {
		toolResources, diags := convertToolResourcesToOpenAI(ctx, plan.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		assistantReq.ToolResources = toolResources
	}

	// Create the assistant
	assistant, err := r.client.OpenAI.CreateAssistant(ctx, assistantReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Assistant",
			fmt.Sprintf("Unable to create assistant: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update state with response
	plan.ID = types.StringValue(assistant.ID)
	plan.CreatedAt = types.Int64Value(int64(assistant.CreatedAt))

	// Convert tools back to state
	if len(assistant.Tools) > 0 {
		toolStrings, diags := convertOpenAIToolsToTerraform(ctx, assistant.Tools)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		toolsList, diags := types.ListValueFrom(ctx, types.StringType, toolStrings)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Tools = toolsList
	} else {
		plan.Tools = types.ListNull(types.StringType)
	}

	// Convert tool resources back to state
	if assistant.ToolResources != nil {
		toolResources, diags := convertOpenAIToolResourcesToTerraform(ctx, assistant.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.ToolResources = toolResources
	}

	// Save the updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *AssistantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AssistantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assistantID := state.ID.ValueString()
	if assistantID == "" {
		resp.Diagnostics.AddError(
			"Error Reading Assistant",
			"Assistant ID is empty. Cannot retrieve assistant details.",
		)
		return
	}

	tflog.Debug(ctx, "Reading assistant", map[string]interface{}{
		"assistant_id": assistantID,
	})

	// Retrieve assistant information
	assistant, err := r.client.OpenAI.RetrieveAssistant(ctx, assistantID)
	if err != nil {
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			// Assistant doesn't exist anymore, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Assistant",
			fmt.Sprintf("Unable to read assistant details: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update state with the latest values from the API
	state.Model = types.StringValue(assistant.Model)
	state.CreatedAt = types.Int64Value(int64(assistant.CreatedAt))

	if assistant.Name != nil {
		state.Name = types.StringValue(*assistant.Name)
	} else {
		state.Name = types.StringNull()
	}

	if assistant.Description != nil {
		state.Description = types.StringValue(*assistant.Description)
	} else {
		state.Description = types.StringNull()
	}

	if assistant.Instructions != nil {
		state.Instructions = types.StringValue(*assistant.Instructions)
	} else {
		state.Instructions = types.StringNull()
	}

	// Convert tools
	if len(assistant.Tools) > 0 {
		toolStrings, diags := convertOpenAIToolsToTerraform(ctx, assistant.Tools)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		toolsList, diags := types.ListValueFrom(ctx, types.StringType, toolStrings)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Tools = toolsList
	} else {
		state.Tools = types.ListNull(types.StringType)
	}

	// Convert metadata
	if assistant.Metadata != nil {
		stringMetadata := make(map[string]string)
		for k, v := range assistant.Metadata {
			if strVal, ok := v.(string); ok {
				stringMetadata[k] = strVal
			} else {
				stringMetadata[k] = fmt.Sprintf("%v", v)
			}
		}

		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, stringMetadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Metadata = metadataMap
	} else {
		state.Metadata = types.MapNull(types.StringType)
	}

	// Convert tool resources
	if assistant.ToolResources != nil {
		toolResources, diags := convertOpenAIToolResourcesToTerraform(ctx, assistant.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.ToolResources = toolResources
	}

	// Save the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AssistantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AssistantResourceModel
	var state AssistantResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assistantID := state.ID.ValueString()
	if assistantID == "" {
		resp.Diagnostics.AddError(
			"Error Updating Assistant",
			"Assistant ID is empty. Cannot update assistant.",
		)
		return
	}

	// Create update request
	assistantReq := openai.AssistantRequest{
		Model: plan.Model.ValueString(),
	}

	// Set optional parameters
	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		assistantReq.Name = &name
	}

	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		assistantReq.Description = &description
	}

	if !plan.Instructions.IsNull() {
		instructions := plan.Instructions.ValueString()
		assistantReq.Instructions = &instructions
	}

	// Process tools if provided
	if !plan.Tools.IsNull() {
		var toolStrings []string
		diags := plan.Tools.ElementsAs(ctx, &toolStrings, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tools, toolDiags := convertToolsToOpenAI(ctx, toolStrings)
		resp.Diagnostics.Append(toolDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		assistantReq.Tools = tools
	}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		assistantReq.Metadata = convertToMapStringAny(metadata)
	}

	// Handle tool_resources if provided
	if plan.ToolResources != nil {
		toolResources, diags := convertToolResourcesToOpenAI(ctx, plan.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		assistantReq.ToolResources = toolResources
	}

	// Update the assistant
	assistant, err := r.client.OpenAI.ModifyAssistant(ctx, assistantID, assistantReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Assistant",
			fmt.Sprintf("Unable to update assistant: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update plan with response data
	plan.CreatedAt = types.Int64Value(int64(assistant.CreatedAt))

	// Convert tools back to state
	if len(assistant.Tools) > 0 {
		toolStrings, diags := convertOpenAIToolsToTerraform(ctx, assistant.Tools)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		toolsList, diags := types.ListValueFrom(ctx, types.StringType, toolStrings)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Tools = toolsList
	} else {
		plan.Tools = types.ListNull(types.StringType)
	}

	// Convert tool resources back to state
	if assistant.ToolResources != nil {
		toolResources, diags := convertOpenAIToolResourcesToTerraform(ctx, assistant.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.ToolResources = toolResources
	}

	// Save the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AssistantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AssistantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assistantID := state.ID.ValueString()
	if assistantID == "" {
		// Nothing to delete
		return
	}

	tflog.Debug(ctx, "Deleting assistant", map[string]interface{}{
		"assistant_id": assistantID,
	})

	_, err := r.client.OpenAI.DeleteAssistant(ctx, assistantID)
	if err != nil {
		// If assistant doesn't exist, don't return an error
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Assistant",
			fmt.Sprintf("Unable to delete assistant: %s", r.client.HandleError(err)),
		)
		return
	}
}

func (r *AssistantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to convert from Terraform tools to OpenAI tools
func convertToolsToOpenAI(ctx context.Context, toolNames []string) ([]openai.AssistantTool, diag.Diagnostics) {
	var diags diag.Diagnostics
	openaiTools := make([]openai.AssistantTool, 0, len(toolNames))

	for _, toolName := range toolNames {
		switch toolName {
		case "code_interpreter":
			openaiTools = append(openaiTools, openai.AssistantTool{
				Type: openai.AssistantToolTypeCodeInterpreter,
			})
		case "file_search":
			openaiTools = append(openaiTools, openai.AssistantTool{
				Type: "file_search",
			})
		case "function":
			diags.AddAttributeError(
				path.Root("tools"),
				"Invalid Tool Configuration",
				"Function tool type requires function definition. Use tool_resources to configure functions.",
			)
		default:
			diags.AddAttributeError(
				path.Root("tools"),
				"Invalid Tool Type",
				fmt.Sprintf("Tool type '%s' is not supported. Must be 'code_interpreter', 'file_search', or 'function'.", toolName),
			)
		}
	}

	return openaiTools, diags
}

// Helper function to convert from OpenAI tools to Terraform tools
func convertOpenAIToolsToTerraform(ctx context.Context, tools []openai.AssistantTool) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfTools := make([]string, 0, len(tools))

	for _, tool := range tools {
		switch tool.Type {
		case openai.AssistantToolTypeCodeInterpreter:
			tfTools = append(tfTools, "code_interpreter")
		case "file_search", openai.AssistantToolTypeRetrieval:
			tfTools = append(tfTools, "file_search")
		case openai.AssistantToolTypeFunction:
			tfTools = append(tfTools, "function")
		}
	}

	return tfTools, diags
}

// Helper function to convert map[string]string to map[string]interface{}
func convertToMapStringAny(in map[string]string) map[string]interface{} {
	result := make(map[string]interface{}, len(in))
	for k, v := range in {
		result[k] = v
	}
	return result
}

// Helper function to convert from Terraform tool resources to OpenAI tool resources
func convertToolResourcesToOpenAI(ctx context.Context, toolResourcesAttr *AssistantToolResourcesModel) (*openai.AssistantToolResource, diag.Diagnostics) {
	if toolResourcesAttr == nil {
		return nil, nil
	}

	var diags diag.Diagnostics
	var hasResources bool
	toolResources := &openai.AssistantToolResource{}

	// Only include file_search if it has vector store IDs
	if toolResourcesAttr.FileSearch != nil && len(toolResourcesAttr.FileSearch.VectorStoreIDs) > 0 {
		hasResources = true
		toolResources.FileSearch = &openai.AssistantToolFileSearch{
			VectorStoreIDs: toolResourcesAttr.FileSearch.VectorStoreIDs,
		}
	}

	// Only include code_interpreter if it has file IDs
	if toolResourcesAttr.CodeInterpreter != nil && len(toolResourcesAttr.CodeInterpreter.FileIDs) > 0 {
		hasResources = true
		toolResources.CodeInterpreter = &openai.AssistantToolCodeInterpreter{
			FileIDs: toolResourcesAttr.CodeInterpreter.FileIDs,
		}
	}

	if !hasResources {
		return nil, diags
	}

	return toolResources, diags
}

// Helper function to convert from OpenAI tool resources to Terraform tool resources
func convertOpenAIToolResourcesToTerraform(ctx context.Context, toolResources *openai.AssistantToolResource) (*AssistantToolResourcesModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if toolResources == nil {
		return nil, diags
	}

	var hasResources bool
	tfToolResources := &AssistantToolResourcesModel{}

	// Only convert file_search if it has vector store IDs
	if toolResources.FileSearch != nil && len(toolResources.FileSearch.VectorStoreIDs) > 0 {
		hasResources = true
		tfToolResources.FileSearch = &AssistantToolResourcesFileSearchModel{
			VectorStoreIDs: toolResources.FileSearch.VectorStoreIDs,
		}
	}

	// Only convert code_interpreter if it has file IDs
	if toolResources.CodeInterpreter != nil && len(toolResources.CodeInterpreter.FileIDs) > 0 {
		hasResources = true
		tfToolResources.CodeInterpreter = &AssistantToolResourcesCodeInterpreterModel{
			FileIDs: toolResources.CodeInterpreter.FileIDs,
		}
	}

	if !hasResources {
		return nil, diags
	}

	return tfToolResources, diags
}
