package resources

import (
	"context"
	"fmt"
	"sort"

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
	ID            types.String         `tfsdk:"id"`
	Metadata      types.Map            `tfsdk:"metadata"`
	CreatedAt     types.Int64          `tfsdk:"created_at"`
	Object        types.String         `tfsdk:"object"`
	Tools         types.List           `tfsdk:"tools"`
	ToolResources types.Object         `tfsdk:"tool_resources"`
	Messages      []ThreadMessageModel `tfsdk:"messages"`
}

type ThreadMessageModel struct {
	Role     types.String `tfsdk:"role"`
	Content  types.String `tfsdk:"content"`
	FileIDs  types.List   `tfsdk:"file_ids"`
	Metadata types.Map    `tfsdk:"metadata"`
}

type ThreadToolResourcesModel struct {
	CodeInterpreter *ThreadToolResourcesCodeInterpreterModel `tfsdk:"code_interpreter"`
	FileSearch      *ThreadToolResourcesFileSearchModel      `tfsdk:"file_search"`
}

type ThreadToolResourcesCodeInterpreterModel struct {
	FileIDs types.Set `tfsdk:"file_ids"`
}

type ThreadToolResourcesFileSearchModel struct {
	VectorStoreIDs types.Set `tfsdk:"vector_store_ids"`
}

// Schema returns the schema for this resource.
func (r *ThreadResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages OpenAI Threads, which are conversation contexts that collect and organize messages between users and assistants.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the thread.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of key-value pairs that can be used to store additional information about the thread.",
			},
			"created_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The Unix timestamp (in seconds) for when the thread was created.",
			},
			"object": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The object type, always \"thread\".",
			},
			"tools": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A list of tools enabled for this thread. Valid values are: code_interpreter and file_search.",
			},
		},
		Blocks: map[string]schema.Block{
			"tool_resources": schema.SingleNestedBlock{
				MarkdownDescription: "Resources made available to the thread's tools.",
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
							"vector_store_ids": schema.SetAttribute{
								MarkdownDescription: "Vector store IDs available to the file search tool.",
								ElementType:         types.StringType,
								Optional:            true,
							},
						},
					},
				},
			},
			"messages": schema.ListNestedBlock{
				MarkdownDescription: "Initial messages for the thread.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							MarkdownDescription: "The role of the entity creating the message. Must be either \"user\" or \"assistant\".",
							Required:            true,
						},
						"content": schema.StringAttribute{
							MarkdownDescription: "The content of the message.",
							Required:            true,
						},
						"file_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							DeprecationMessage:  "The file_ids attribute is deprecated in v2 of the Assistants API. Use message attachments instead.",
							MarkdownDescription: "DEPRECATED: A list of file IDs to attach to the message.",
						},
						"metadata": schema.MapAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: "A map of key-value pairs that can be used to store additional information about the message.",
						},
					},
				},
			},
		},
	}
}

func (r *ThreadResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_thread"
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
	plan.Object = types.StringValue(thread.Object)

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
	state.Object = types.StringValue(thread.Object)

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
	plan.Object = types.StringValue(thread.Object)

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
	if val == nil {
		return result, diags
	}

	// Handle code_interpreter
	if codeInterpreter, ok := val["code_interpreter"].(types.Object); ok && !codeInterpreter.IsNull() {
		attrs := codeInterpreter.Attributes()
		if fileIDs, ok := attrs["file_ids"].(types.Set); ok && !fileIDs.IsNull() {
			var ids []string
			diags.Append(fileIDs.ElementsAs(ctx, &ids, false)...)
			if diags.HasError() {
				return result, diags
			}
			// Sort the file IDs for consistency
			sort.Strings(ids)
			result.CodeInterpreter = &openai.CodeInterpreterToolResources{
				FileIDs: ids,
			}
		}
	}

	// Handle file_search
	if fileSearch, ok := val["file_search"].(types.Object); ok && !fileSearch.IsNull() {
		attrs := fileSearch.Attributes()
		if vectorStoreIDs, ok := attrs["vector_store_ids"].(types.Set); ok && !vectorStoreIDs.IsNull() {
			var ids []string
			diags.Append(vectorStoreIDs.ElementsAs(ctx, &ids, false)...)
			if diags.HasError() {
				return result, diags
			}
			// Sort the vector store IDs for consistency
			sort.Strings(ids)
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

	// Define the object type structure for nested objects
	attrTypes := map[string]attr.Type{
		"code_interpreter": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"file_ids": types.SetType{
					ElemType: types.StringType,
				},
			},
		},
		"file_search": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"vector_store_ids": types.SetType{
					ElemType: types.StringType,
				},
			},
		},
	}

	// If no tool resources were present in the API response, return null
	if toolResources == (openai.ToolResources{}) {
		return types.ObjectNull(attrTypes), diags
	}

	// Build the object value starting with empty objects for both tools
	attrs := map[string]attr.Value{
		"code_interpreter": types.ObjectNull(map[string]attr.Type{
			"file_ids": types.SetType{
				ElemType: types.StringType,
			},
		}),
		"file_search": types.ObjectNull(map[string]attr.Type{
			"vector_store_ids": types.SetType{
				ElemType: types.StringType,
			},
		}),
	}

	// Update code_interpreter if present in response
	if toolResources.CodeInterpreter != nil {
		// Sort file IDs for consistency
		sortedFileIDs := make([]string, len(toolResources.CodeInterpreter.FileIDs))
		copy(sortedFileIDs, toolResources.CodeInterpreter.FileIDs)
		sort.Strings(sortedFileIDs)

		fileIDsVal, d := types.SetValueFrom(ctx, types.StringType, sortedFileIDs)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}

		attrs["code_interpreter"], d = types.ObjectValue(
			map[string]attr.Type{
				"file_ids": types.SetType{
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
	}

	// Update file_search if present in response
	if toolResources.FileSearch != nil {
		// Sort vector store IDs for consistency
		sortedVectorStoreIDs := make([]string, len(toolResources.FileSearch.VectorStoreIDs))
		copy(sortedVectorStoreIDs, toolResources.FileSearch.VectorStoreIDs)
		sort.Strings(sortedVectorStoreIDs)

		vectorStoreIDsVal, d := types.SetValueFrom(ctx, types.StringType, sortedVectorStoreIDs)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}

		attrs["file_search"], d = types.ObjectValue(
			map[string]attr.Type{
				"vector_store_ids": types.SetType{
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
	}

	// Create and return the final object
	return types.ObjectValue(attrTypes, attrs)
}

// Helper function to convert from Terraform tool_resources to OpenAI ToolResourcesRequest
func convertTerraformToolResourcesToOpenAIRequest(ctx context.Context, toolResources types.Object) (*openai.ToolResourcesRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// If tool_resources is null or empty, return nil
	if toolResources.IsNull() || toolResources.IsUnknown() {
		return nil, diags
	}

	val := toolResources.Attributes()
	if val == nil || len(val) == 0 {
		return nil, diags
	}

	result := &openai.ToolResourcesRequest{}

	// Handle code_interpreter if present
	if codeInterpreterVal, ok := val["code_interpreter"]; ok {
		if codeInterpreterObj, ok := codeInterpreterVal.(types.Object); ok && !codeInterpreterObj.IsNull() && !codeInterpreterObj.IsUnknown() {
			attrs := codeInterpreterObj.Attributes()
			if attrs != nil {
				if fileIDs, ok := attrs["file_ids"].(types.Set); ok && !fileIDs.IsNull() && !fileIDs.IsUnknown() {
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
		}
	}

	// Handle file_search if present
	if fileSearchVal, ok := val["file_search"]; ok {
		if fileSearchObj, ok := fileSearchVal.(types.Object); ok && !fileSearchObj.IsNull() && !fileSearchObj.IsUnknown() {
			attrs := fileSearchObj.Attributes()
			if attrs != nil {
				if vectorStoreIDs, ok := attrs["vector_store_ids"].(types.Set); ok && !vectorStoreIDs.IsNull() && !vectorStoreIDs.IsUnknown() {
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
		}
	}

	// If no tools were configured, return nil
	if result.CodeInterpreter == nil && result.FileSearch == nil {
		return nil, diags
	}

	return result, diags
}
