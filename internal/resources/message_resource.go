package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sashabaranov/go-openai"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MessageResource{}
var _ resource.ResourceWithImportState = &MessageResource{}

func NewMessageResource() resource.Resource {
	return &MessageResource{}
}

// MessageResource defines the resource implementation.
type MessageResource struct {
	client *client.Client
}

// MessageResourceModel describes the resource data model.
type MessageResourceModel struct {
	ID          types.String      `tfsdk:"id"`
	Object      types.String      `tfsdk:"object"`
	ThreadID    types.String      `tfsdk:"thread_id"`
	Role        types.String      `tfsdk:"role"`
	Content     types.String      `tfsdk:"content"`
	FileIDs     []string          `tfsdk:"file_ids"`
	Metadata    map[string]string `tfsdk:"metadata"`
	AssistantID types.String      `tfsdk:"assistant_id"`
	RunID       types.String      `tfsdk:"run_id"`
	CreatedAt   types.Int64       `tfsdk:"created_at"`
}

func (r *MessageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message"
}

func (r *MessageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage messages within OpenAI Threads. Messages are the building blocks of conversations between users and assistants.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the message.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
				"object": schema.StringAttribute{
				MarkdownDescription: "The object type, always 'thread.message'.",
				Computed:            true,
			},
			"thread_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the thread this message belongs to.",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the entity that created this message. Currently supported values are 'user' or 'assistant'.",
				Required:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The content of the message.",
				Required:            true,
			},
			"file_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A list of File IDs that the message should use.",
				Optional:            true,
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of key-value pairs that can be attached to an object. This can be useful for storing additional information about the object in a structured format.",
				Optional:            true,
			},
			"assistant_id": schema.StringAttribute{
				MarkdownDescription: "If set, the ID of the assistant that authored this message.",
				Optional:            true,
				Computed:            true,
			},
			"run_id": schema.StringAttribute{
				MarkdownDescription: "If set, the ID of the run associated with this message.",
				Optional:            true,
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) when the message was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *MessageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MessageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MessageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the message request
	messageReq := openai.MessageRequest{
		Role:    plan.Role.ValueString(),
		Content: plan.Content.ValueString(),
	}

	// Add optional fields if specified
	if len(plan.FileIDs) > 0 {
		messageReq.FileIds = plan.FileIDs
	}

	if len(plan.Metadata) > 0 {
		// Convert from map[string]string to map[string]interface{}
		metadataAny := make(map[string]interface{})
		for k, v := range plan.Metadata {
			metadataAny[k] = v
		}
		messageReq.Metadata = metadataAny
	}

	// Create the message
	message, err := r.client.OpenAI.CreateMessage(ctx, plan.ThreadID.ValueString(), messageReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Message",
			fmt.Sprintf("Unable to create message: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the plan with values from the response
	plan.ID = types.StringValue(message.ID)
	plan.Object = types.StringValue(message.Object)
	plan.CreatedAt = types.Int64Value(int64(message.CreatedAt))
	plan.Role = types.StringValue(message.Role)

	// Handle file IDs in response
	if len(message.FileIds) > 0 {
		plan.FileIDs = message.FileIds
	} else {
		plan.FileIDs = nil
	}

	// Ensure content is properly set from response
	if len(message.Content) > 0 {
		for _, content := range message.Content {
			if content.Type == "text" {
				plan.Content = types.StringValue(content.Text.Value)
				break
			}
		}
	}

	// Handle optional fields
	if message.AssistantID != nil {
		plan.AssistantID = types.StringValue(*message.AssistantID)
	} else {
		plan.AssistantID = types.StringNull()
	}

	if message.RunID != nil {
		plan.RunID = types.StringValue(*message.RunID)
	} else {
		plan.RunID = types.StringNull()
	}

	// Save the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MessageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MessageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	threadID := state.ThreadID.ValueString()
	messageID := state.ID.ValueString()
	if threadID == "" || messageID == "" {
		resp.Diagnostics.AddError(
			"Error Reading Message",
			"Thread ID and Message ID are required to read a message.",
		)
		return
	}

	tflog.Debug(ctx, "Reading message", map[string]interface{}{
		"thread_id":  threadID,
		"message_id": messageID,
	})

	// Retrieve message information
	message, err := r.client.OpenAI.RetrieveMessage(ctx, threadID, messageID)
	if err != nil {
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			// Message doesn't exist anymore, remove from state
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Message",
			fmt.Sprintf("Unable to read message details: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update state with latest values
	state.Object = types.StringValue(message.Object) // Add this line
	state.Role = types.StringValue(message.Role)
	state.CreatedAt = types.Int64Value(int64(message.CreatedAt))

	// Set content from the first content item if available
	if len(message.Content) > 0 && message.Content[0].Type == "text" {
		state.Content = types.StringValue(message.Content[0].Text.Value)
	}

	// Convert file IDs
	if len(message.FileIds) > 0 {
		state.FileIDs = message.FileIds
	} else {
		state.FileIDs = nil
	}

	// Handle optional fields
	if message.AssistantID != nil {
		state.AssistantID = types.StringValue(*message.AssistantID)
	} else {
		state.AssistantID = types.StringNull()
	}

	if message.RunID != nil {
		state.RunID = types.StringValue(*message.RunID)
	} else {
		state.RunID = types.StringNull()
	}

	// Convert metadata
	if message.Metadata != nil {
		// Convert from map[string]interface{} to map[string]string for Terraform state
		metadataStr := make(map[string]string)
		for k, v := range message.Metadata {
			if strValue, ok := v.(string); ok {
				metadataStr[k] = strValue
			} else {
				metadataStr[k] = fmt.Sprintf("%v", v)
			}
		}
		state.Metadata = metadataStr
	} else {
		state.Metadata = nil
	}

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MessageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MessageResourceModel
	var state MessageResourceModel

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

	threadID := state.ThreadID.ValueString()
	messageID := state.ID.ValueString()
	if threadID == "" || messageID == "" {
		resp.Diagnostics.AddError(
			"Error Updating Message",
			"Thread ID and Message ID are required to update a message.",
		)
		return
	}

	// The v2 API only allows updating metadata
	if !equalMetadata(plan.Metadata, state.Metadata) {
		// Convert metadata to the format expected by ModifyMessage
		metadata := plan.Metadata
		if metadata == nil {
			metadata = make(map[string]string)
		}

		// Update the message
		message, err := r.client.OpenAI.ModifyMessage(ctx, threadID, messageID, metadata)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Message",
				fmt.Sprintf("Unable to update message metadata: %s", r.client.HandleError(err)),
			)
			return
		}

		// Update state from response
		plan.Object = types.StringValue(message.Object) // Add this line
		if message.AssistantID != nil {
			plan.AssistantID = types.StringValue(*message.AssistantID)
		} else {
			plan.AssistantID = types.StringNull()
		}

		if message.RunID != nil {
			plan.RunID = types.StringValue(*message.RunID)
		} else {
			plan.RunID = types.StringNull()
		}
	} else {
		// If non-metadata fields changed, return an error as they cannot be modified
		if !plan.Content.Equal(state.Content) ||
			!plan.Role.Equal(state.Role) ||
			!equalStringSlice(plan.FileIDs, state.FileIDs) {
			resp.Diagnostics.AddError(
				"Cannot Update Message Fields",
				"Message content, role, and file IDs cannot be modified after creation. Only metadata can be updated.",
			)
			return
		}
	}

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MessageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MessageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	threadID := state.ThreadID.ValueString()
	messageID := state.ID.ValueString()
	if threadID == "" || messageID == "" {
		// Nothing to delete
		return
	}

	tflog.Debug(ctx, "Deleting message", map[string]interface{}{
		"thread_id":  threadID,
		"message_id": messageID,
	})

	// Note: The OpenAI API currently does not support deleting individual messages
	// from a thread. This is a placeholder for when that functionality is added.
	//
	// For now, we'll just remove it from state
	tflog.Info(ctx, "OpenAI API does not currently support deleting individual messages. Resource will be removed from state only.")
}

func (r *MessageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: thread_id:message_id
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Error Importing Message",
			"Invalid import ID format. Expected format: thread_id:message_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("thread_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

// Helper function to wait for a run to complete
func (r *MessageResource) waitForRunCompletion(ctx context.Context, threadID, runID string) (*openai.Run, error) {
	maxAttempts := 30
	pollingInterval := 2 * time.Second

	for i := 0; i < maxAttempts; i++ {
		run, err := r.client.OpenAI.RetrieveRun(ctx, threadID, runID)
		if err != nil {
			return nil, err
		}

		// Check if run status is terminal
		switch run.Status {
		case "completed", "failed", "cancelled", "expired":
			return &run, nil
		case "queued", "in_progress", "requires_action":
			// Continue waiting
			time.Sleep(pollingInterval)
		default:
			// Unknown status
			return &run, fmt.Errorf("run has unknown status: %s", run.Status)
		}
	}

	return nil, fmt.Errorf("run did not complete within the timeout period")
}

func equalMetadata(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
