package resources

import (
	"context"
	"fmt"

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
var _ resource.Resource = &ThreadResource{}
var _ resource.ResourceWithImportState = &ThreadResource{}

func NewThreadResource() resource.Resource {
	return &ThreadResource{}
}

// ThreadResource defines the resource implementation.
type ThreadResource struct {
	client *client.Client
}

// ThreadResourceModel describes the resource data model.
type ThreadResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Metadata      types.Map    `tfsdk:"metadata"`
	ToolResources types.Object `tfsdk:"tool_resources"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
}

func (r *ThreadResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_thread"
}

func (r *ThreadResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage conversation threads for OpenAI Assistants.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Metadata in key-value pairs to attach to the thread.",
				Optional:            true,
				Computed:            false,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the thread was created.",
				Computed:            true,
			},
			"tool_resources": schema.SingleNestedAttribute{
				MarkdownDescription: "Resources available to the tools used in the thread.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"code_interpreter": schema.SingleNestedAttribute{
						MarkdownDescription: "Code interpreter tool resources",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"file_ids": schema.ListAttribute{
								ElementType:         types.StringType,
								MarkdownDescription: "List of file IDs to use with code interpreter",
								Optional:            true,
							},
						},
					},
					"file_search": schema.SingleNestedAttribute{
						MarkdownDescription: "File search tool resources",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"vector_store_ids": schema.ListAttribute{
								ElementType:         types.StringType,
								MarkdownDescription: "List of vector store IDs to use with file search",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *ThreadResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ThreadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ThreadResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the thread request
	threadReq := openai.ThreadRequest{
		Metadata: make(map[string]interface{}),
	}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Convert metadata to map[string]any
		for k, v := range metadata {
			threadReq.Metadata[k] = v
		}
	}

	// Handle tool_resources if provided
	if !plan.ToolResources.IsNull() {
		toolResourcesState, diags := convertTerraformToolResourcesToOpenAIRequest(ctx, plan.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		threadReq.ToolResources = toolResourcesState
	}

	tflog.Debug(ctx, "Creating thread")

	// Create the thread
	thread, err := r.client.OpenAI.CreateThread(ctx, threadReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Thread",
			fmt.Sprintf("Unable to create thread: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state
	plan.ID = types.StringValue(thread.ID)
	plan.CreatedAt = types.Int64Value(int64(thread.CreatedAt))

	// Handle metadata in response
	if thread.Metadata != nil {
		metadataMap := make(map[string]string)
		for k, v := range thread.Metadata {
			if strVal, ok := v.(string); ok {
				metadataMap[k] = strVal
			}
		}
		metadataValue, diags := types.MapValueFrom(ctx, types.StringType, metadataMap)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Metadata = metadataValue
	}

	// Handle tool_resources in response
	if thread.ToolResources != (openai.ToolResources{}) {
		toolResourcesState, diags := convertOpenAIToolResourcesToState(ctx, thread.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.ToolResources = toolResourcesState
	}

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ThreadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ThreadResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	threadID := state.ID.ValueString()
	if threadID == "" {
		resp.Diagnostics.AddError(
			"Error Reading Thread",
			"Thread ID is empty. Cannot retrieve thread details.",
		)
		return
	}

	tflog.Debug(ctx, "Reading thread", map[string]interface{}{
		"thread_id": threadID,
	})

	// Retrieve thread information
	thread, err := r.client.OpenAI.RetrieveThread(ctx, threadID)
	if err != nil {
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			// Thread doesn't exist anymore, remove from state
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Thread",
			fmt.Sprintf("Unable to read thread details: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state with the latest values
	state.CreatedAt = types.Int64Value(int64(thread.CreatedAt))

	// Convert metadata
	if thread.Metadata != nil {
		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, thread.Metadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Metadata = metadataMap
	}

	// Handle tool_resources in response
	if thread.ToolResources != (openai.ToolResources{}) {
		toolResourcesState, diags := convertOpenAIToolResourcesToState(ctx, thread.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.ToolResources = toolResourcesState
	}

	// Save into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ThreadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ThreadResourceModel
	var state ThreadResourceModel

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

	threadID := state.ID.ValueString()
	if threadID == "" {
		resp.Diagnostics.AddError(
			"Error Updating Thread",
			"Thread ID is empty. Cannot update thread.",
		)
		return
	}

	// Create the update request
	threadReq := openai.ModifyThreadRequest{
		Metadata: make(map[string]interface{}),
	}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Convert metadata to map[string]any
		for k, v := range metadata {
			threadReq.Metadata[k] = v
		}
	}

	// Handle tool_resources if provided
	if !plan.ToolResources.IsNull() {
		toolResources, diags := convertTerraformToolResourcesToOpenAI(ctx, plan.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		threadReq.ToolResources = &toolResources
	}

	tflog.Debug(ctx, "Updating thread", map[string]interface{}{
		"thread_id": threadID,
	})

	// Update the thread
	thread, err := r.client.OpenAI.ModifyThread(ctx, threadID, threadReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Thread",
			fmt.Sprintf("Unable to update thread: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state
	plan.CreatedAt = types.Int64Value(int64(thread.CreatedAt))

	// Handle metadata in response
	if thread.Metadata != nil {
		metadataMap := make(map[string]string)
		for k, v := range thread.Metadata {
			if strVal, ok := v.(string); ok {
				metadataMap[k] = strVal
			}
		}
		metadataValue, diags := types.MapValueFrom(ctx, types.StringType, metadataMap)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Metadata = metadataValue
	}

	// Handle tool_resources in response
	if thread.ToolResources != (openai.ToolResources{}) {
		toolResourcesState, diags := convertOpenAIToolResourcesToState(ctx, thread.ToolResources)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.ToolResources = toolResourcesState
	}

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ThreadResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ThreadResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	threadID := state.ID.ValueString()
	if threadID == "" {
		// Nothing to delete
		return
	}

	tflog.Debug(ctx, "Deleting thread", map[string]interface{}{
		"thread_id": threadID,
	})

	_, err := r.client.OpenAI.DeleteThread(ctx, threadID)
	if err != nil {
		// If thread doesn't exist, don't return an error
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting Thread",
			fmt.Sprintf("Unable to delete thread: %s", r.client.HandleError(err)),
		)
		return
	}
}

func (r *ThreadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to convert from Terraform tool_resources to OpenAI ToolResources
func convertTerraformToolResourcesToOpenAI(ctx context.Context, toolResources types.Object) (openai.ToolResources, diag.Diagnostics) {
	var result openai.ToolResources
	var diags diag.Diagnostics

	val := toolResources.Attributes()

	// Handle code_interpreter
	if codeInterpreter, ok := val["code_interpreter"].(types.Object); ok && !codeInterpreter.IsNull() {
		attrs := codeInterpreter.Attributes()
		if fileIDs, ok := attrs["file_ids"].(types.List); ok && !fileIDs.IsNull() {
			var ids []string
			diags.Append(fileIDs.ElementsAs(ctx, &ids, false)...)
			if diags.HasError() {
				return result, diags
			}
			result.CodeInterpreter = &openai.CodeInterpreterToolResources{
				FileIDs: ids,
			}
		}
	}

	// Handle file_search
	if fileSearch, ok := val["file_search"].(types.Object); ok && !fileSearch.IsNull() {
		attrs := fileSearch.Attributes()
		if vectorStoreIDs, ok := attrs["vector_store_ids"].(types.List); ok && !vectorStoreIDs.IsNull() {
			var ids []string
			diags.Append(vectorStoreIDs.ElementsAs(ctx, &ids, false)...)
			if diags.HasError() {
				return result, diags
			}
			result.FileSearch = &openai.FileSearchToolResources{
				VectorStoreIDs: ids,
			}
		}
	}

	return result, diags
}

// Helper function to convert from OpenAI ToolResources to Terraform state
func convertOpenAIToolResourcesToState(ctx context.Context, toolResources openai.ToolResources) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Define the object type structure
	attrTypes := map[string]attr.Type{
		"code_interpreter": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"file_ids": types.ListType{
					ElemType: types.StringType,
				},
			},
		},
		"file_search": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"vector_store_ids": types.ListType{
					ElemType: types.StringType,
				},
			},
		},
	}

	// Build the object value
	attrs := make(map[string]attr.Value)

	// Handle code_interpreter
	if toolResources.CodeInterpreter != nil {
		fileIDsVal, d := types.ListValueFrom(ctx, types.StringType, toolResources.CodeInterpreter.FileIDs)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}

		codeInterpreterVal, d := types.ObjectValue(
			map[string]attr.Type{
				"file_ids": types.ListType{
					ElemType: types.StringType,
				},
			},
			map[string]attr.Value{
				"file_ids": fileIDsVal,
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}
		attrs["code_interpreter"] = codeInterpreterVal
	}

	// Handle file_search
	if toolResources.FileSearch != nil {
		vectorStoreIDsVal, d := types.ListValueFrom(ctx, types.StringType, toolResources.FileSearch.VectorStoreIDs)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}

		fileSearchVal, d := types.ObjectValue(
			map[string]attr.Type{
				"vector_store_ids": types.ListType{
					ElemType: types.StringType,
				},
			},
			map[string]attr.Value{
				"vector_store_ids": vectorStoreIDsVal,
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}
		attrs["file_search"] = fileSearchVal
	}

	return types.ObjectValue(attrTypes, attrs)
}

// Helper function to convert from Terraform tool_resources to OpenAI ToolResourcesRequest
func convertTerraformToolResourcesToOpenAIRequest(ctx context.Context, toolResources types.Object) (*openai.ToolResourcesRequest, diag.Diagnostics) {
	var result openai.ToolResourcesRequest
	var diags diag.Diagnostics

	val := toolResources.Attributes()

	// Handle code_interpreter
	if codeInterpreter, ok := val["code_interpreter"].(types.Object); ok && !codeInterpreter.IsNull() {
		attrs := codeInterpreter.Attributes()
		if fileIDs, ok := attrs["file_ids"].(types.List); ok && !fileIDs.IsNull() {
			var ids []string
			diags.Append(fileIDs.ElementsAs(ctx, &ids, false)...)
			if diags.HasError() {
				return nil, diags
			}
			result.CodeInterpreter = &openai.CodeInterpreterToolResourcesRequest{
				FileIDs: ids,
			}
		}
	}

	// Handle file_search
	if fileSearch, ok := val["file_search"].(types.Object); ok && !fileSearch.IsNull() {
		attrs := fileSearch.Attributes()
		if vectorStoreIDs, ok := attrs["vector_store_ids"].(types.List); ok && !vectorStoreIDs.IsNull() {
			var ids []string
			diags.Append(vectorStoreIDs.ElementsAs(ctx, &ids, false)...)
			if diags.HasError() {
				return nil, diags
			}
			result.FileSearch = &openai.FileSearchToolResourcesRequest{
				VectorStoreIDs: ids,
			}
		}
	}

	return &result, diags
}
