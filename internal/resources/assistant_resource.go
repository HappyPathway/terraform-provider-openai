package resources

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Model        types.String `tfsdk:"model"`
	Instructions types.String `tfsdk:"instructions"`
	Tools        types.Set    `tfsdk:"tools"`
	FilePath     types.String `tfsdk:"file_path"`
	Metadata     types.Map    `tfsdk:"metadata"`
	ObjectID     types.String `tfsdk:"object_id"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
}

// AssistantToolModel represents a tool configuration for an assistant.
type AssistantToolModel struct {
	Type               types.String `tfsdk:"type"`
	FunctionDefinition types.String `tfsdk:"function_definition"`
}

func (r *AssistantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assistant"
}

func (r *AssistantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage OpenAI Assistants to help with tasks using models, tools, and knowledge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the assistant. This is visible in the OpenAI dashboard.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the assistant's purpose and capabilities.",
				Optional:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The model that the assistant will use (e.g., 'gpt-4', 'gpt-3.5-turbo').",
				Required:            true,
			},
			"instructions": schema.StringAttribute{
				MarkdownDescription: "Instructions that the assistant uses to guide its responses.",
				Optional:            true,
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "Path to the file to be used by the assistant. Must be accessible by the provider.",
				Optional:            true,
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Metadata in key-value pairs to attach to the assistant.",
				Optional:            true,
			},
			"object_id": schema.StringAttribute{
				MarkdownDescription: "The OpenAI ID assigned to this assistant.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the assistant was created.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"tools": schema.SetNestedBlock{
				MarkdownDescription: "A list of tools the assistant can use. Can include 'code_interpreter', 'file_search', or function tools.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of tool. Can be 'code_interpreter', 'file_search', or 'function'.",
							Required:            true,
						},
						"function_definition": schema.StringAttribute{
							MarkdownDescription: "JSON string of the function definition when type is 'function'.",
							Optional:            true,
						},
					},
				},
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

	// First, upload the file if provided
	var uploadedFileID string
	if (!plan.FilePath.IsNull()) {
		filePath := plan.FilePath.ValueString()
		
		// Create a file upload
		file, err := os.Open(filePath)
		if (err != nil) {
			resp.Diagnostics.AddError(
				"Error Opening File",
				fmt.Sprintf("Unable to open file %s: %s", filePath, err),
			)
			return
		}
		defer file.Close()

		fileReq := openai.FileRequest{
			FileName: filepath.Base(filePath),
			FilePath: filePath,
			Purpose: "assistants",
		}

		uploadedFile, err := r.client.OpenAI.CreateFile(ctx, fileReq)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Uploading File",
				fmt.Sprintf("Unable to upload file %s: %s", filePath, r.client.HandleError(err)),
			)
			return
		}
		uploadedFileID = uploadedFile.ID
	}

	// Create the assistant request
	assistantReq := openai.AssistantRequest{
		Model: plan.Model.ValueString(),
	}

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
		tools, diags := convertToolsToOpenAI(ctx, plan.Tools)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		assistantReq.Tools = tools
	}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		var stringMetadata map[string]string
		diags := plan.Metadata.ElementsAs(ctx, &stringMetadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Convert string values to interface{} values
		metadata := make(map[string]interface{})
		for k, v := range stringMetadata {
			metadata[k] = v
		}
		assistantReq.Metadata = metadata
	}

	// Set up file resources if we have a file
	if uploadedFileID != "" {
		assistantReq.ToolResources = &openai.AssistantToolResource{
			FileSearch: &openai.AssistantToolFileSearch{
				VectorStoreIDs: []string{uploadedFileID},
			},
			CodeInterpreter: &openai.AssistantToolCodeInterpreter{
				FileIDs: []string{uploadedFileID},
			},
		}
	}

	// Create the assistant with everything included
	assistant, err := r.client.OpenAI.CreateAssistant(ctx, assistantReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Assistant",
			fmt.Sprintf("Unable to create assistant: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state with the response
	plan.ID = types.StringValue(assistant.ID)
	plan.ObjectID = types.StringValue(assistant.ID)
	plan.CreatedAt = types.Int64Value(int64(assistant.CreatedAt))

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AssistantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AssistantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assistantID := state.ObjectID.ValueString()
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

	// Update the state with the latest values
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

	state.Model = types.StringValue(assistant.Model)

	if assistant.Instructions != nil {
		state.Instructions = types.StringValue(*assistant.Instructions)
	} else {
		state.Instructions = types.StringNull()
	}

	state.CreatedAt = types.Int64Value(int64(assistant.CreatedAt))

	// Convert tools
	if len(assistant.Tools) > 0 {
		toolsList, diags := convertOpenAIToolsToTerraform(ctx, assistant.Tools)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Tools = toolsList
	}

	// Keep the file_path from the existing state
	// We don't update it from the API since it's a local path

	// Convert metadata
	if assistant.Metadata != nil {
		// Convert from map[string]any to map[string]string
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
	}

	// Save into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AssistantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AssistantResourceModel
	var state AssistantResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	assistantID := state.ObjectID.ValueString()
	if assistantID == "" {
		resp.Diagnostics.AddError(
			"Error Updating Assistant",
			"Assistant ID is empty. Cannot update assistant.",
		)
		return
	}

	// Create the update request
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
		tools, toolDiags := convertToolsToOpenAI(ctx, plan.Tools)
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

	tflog.Debug(ctx, "Updating assistant", map[string]interface{}{
		"assistant_id": assistantID,
		"model":        assistantReq.Model,
	})

	// Update the assistant
	assistant, err := r.client.OpenAI.ModifyAssistant(ctx, assistantID, assistantReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Assistant",
			fmt.Sprintf("Unable to update assistant: %s", r.client.HandleError(err)),
		)
		return
	}

	// If file_path has changed, handle the file update
	if !plan.FilePath.Equal(state.FilePath) {
		// First, upload the new file if provided
		if !plan.FilePath.IsNull() {
			filePath := plan.FilePath.ValueString()
			
			// Create a file upload
			file, err := os.Open(filePath)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Opening File",
					fmt.Sprintf("Unable to open file %s: %s", filePath, err),
				)
				return
			}
			defer file.Close()

			fileReq := openai.FileRequest{
				FileName: filepath.Base(filePath),
				FilePath: filePath,
				Purpose: "assistants",
			}

			uploadedFile, err := r.client.OpenAI.CreateFile(ctx, fileReq)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Uploading File",
					fmt.Sprintf("Unable to upload file %s: %s", filePath, r.client.HandleError(err)),
				)
				return
			}

			// Then, attach the file to the assistant
			fileAttachReq := openai.AssistantFileRequest{
				FileID: uploadedFile.ID,
			}
			_, err = r.client.OpenAI.CreateAssistantFile(ctx, assistant.ID, fileAttachReq)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Attaching File",
					fmt.Sprintf("Unable to attach file %s: %s", uploadedFile.ID, r.client.HandleError(err)),
				)
				return
			}
		}
	}

	// Update the state
	plan.ObjectID = types.StringValue(assistant.ID)
	plan.CreatedAt = types.Int64Value(int64(assistant.CreatedAt))

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AssistantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AssistantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assistantID := state.ObjectID.ValueString()
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
	resource.ImportStatePassthroughID(ctx, path.Root("object_id"), req, resp)
}

// Helper function to convert from Terraform tools to OpenAI tools
func convertToolsToOpenAI(ctx context.Context, toolsAttr types.Set) ([]openai.AssistantTool, diag.Diagnostics) {
	var diags diag.Diagnostics
	var tools []AssistantToolModel

	if toolsAttr.IsNull() || toolsAttr.IsUnknown() {
		return nil, diags
	}

	diags.Append(toolsAttr.ElementsAs(ctx, &tools, false)...)
	if diags.HasError() {
		return nil, diags
	}

	openaiTools := make([]openai.AssistantTool, 0, len(tools))

	for _, tool := range tools {
		if tool.Type.IsNull() {
			diags.AddAttributeError(
				path.Root("tools"),
				"Invalid Tool",
				"Tool must have a type.",
			)
			continue
		}

		toolType := tool.Type.ValueString()
		switch toolType {
		case "code_interpreter":
			openaiTools = append(openaiTools, openai.AssistantTool{
				Type: openai.AssistantToolTypeCodeInterpreter,
			})
		case "file_search":
			openaiTools = append(openaiTools, openai.AssistantTool{
				Type: "file_search", // Use file_search directly instead of converting to retrieval
			})
		case "function":
			if tool.FunctionDefinition.IsNull() {
				diags.AddAttributeError(
					path.Root("tools").AtName("function_definition"),
					"Missing Function Definition",
					"Function tools must have a function_definition.",
				)
				continue
			}
			openaiTools = append(openaiTools, openai.AssistantTool{
				Type: openai.AssistantToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Description: tool.FunctionDefinition.ValueString(),
				},
			})
		default:
			diags.AddAttributeError(
				path.Root("tools").AtName("type"),
				"Invalid Tool Type",
				fmt.Sprintf("Tool type '%s' is not supported. Must be 'code_interpreter', 'file_search', or 'function'.", toolType),
			)
		}
	}

	return openaiTools, diags
}

// Helper function to convert from OpenAI tools to Terraform tools
func convertOpenAIToolsToTerraform(ctx context.Context, tools []openai.AssistantTool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfTools := make([]attr.Value, 0, len(tools))

	for _, tool := range tools {
		toolMap := make(map[string]attr.Value)
		switch tool.Type {
		case openai.AssistantToolTypeCodeInterpreter:
			toolMap["type"] = types.StringValue("code_interpreter")
			toolMap["function_definition"] = types.StringNull()
		case "file_search", openai.AssistantToolTypeRetrieval: // Handle both file_search and retrieval
			toolMap["type"] = types.StringValue("file_search")
			toolMap["function_definition"] = types.StringNull()
		case openai.AssistantToolTypeFunction:
			toolMap["type"] = types.StringValue("function")
			if tool.Function != nil {
				toolMap["function_definition"] = types.StringValue(tool.Function.Description)
			} else {
				toolMap["function_definition"] = types.StringNull()
			}
		}
		
		toolObj, d := types.ObjectValue(
			map[string]attr.Type{
				"type":                types.StringType,
				"function_definition": types.StringType,
			},
			toolMap,
		)
		diags.Append(d...)
		if diags.HasError() {
			continue
		}
		tfTools = append(tfTools, toolObj)
	}

	toolsSet, d := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":                types.StringType,
				"function_definition": types.StringType,
			},
		},
		tfTools,
	)
	diags.Append(d...)
	return toolsSet, diags
}

// Helper function to convert map[string]string to map[string]interface{}
func convertToMapStringAny(in map[string]string) map[string]interface{} {
	result := make(map[string]interface{}, len(in))
	for k, v := range in {
		result[k] = v
	}
	return result
}
