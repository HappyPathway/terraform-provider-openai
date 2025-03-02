package resources

import (
	"context"
	"fmt"

	"github.com/darnold/terraform-provider-openai/internal/client"
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
	ID        types.String `tfsdk:"id"`
	Metadata  types.Map    `tfsdk:"metadata"`
	ObjectID  types.String `tfsdk:"object_id"`
	CreatedAt types.Int64  `tfsdk:"created_at"`
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
			},
			"object_id": schema.StringAttribute{
				MarkdownDescription: "The OpenAI ID assigned to this thread.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the thread was created.",
				Computed:            true,
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
	threadReq := openai.ThreadRequest{}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Convert metadata to map[string]any
		metadataAny := make(map[string]any)
		for k, v := range metadata {
			metadataAny[k] = v
		}
		threadReq.Metadata = metadataAny
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
	plan.ObjectID = types.StringValue(thread.ID)
	plan.CreatedAt = types.Int64Value(int64(thread.CreatedAt))

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

	threadID := state.ObjectID.ValueString()
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

	threadID := state.ObjectID.ValueString()
	if threadID == "" {
		resp.Diagnostics.AddError(
			"Error Updating Thread",
			"Thread ID is empty. Cannot update thread.",
		)
		return
	}

	// Create the update request
	threadReq := openai.ModifyThreadRequest{}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Convert metadata to map[string]any
		metadataAny := make(map[string]any)
		for k, v := range metadata {
			metadataAny[k] = v
		}
		threadReq.Metadata = metadataAny
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
	plan.ObjectID = types.StringValue(thread.ID)
	plan.CreatedAt = types.Int64Value(int64(thread.CreatedAt))

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

	threadID := state.ObjectID.ValueString()
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
	resource.ImportStatePassthroughID(ctx, path.Root("object_id"), req, resp)
}
