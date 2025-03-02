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
	ID              types.String `tfsdk:"id"`
	ThreadID        types.String `tfsdk:"thread_id"`
	Role            types.String `tfsdk:"role"`
	Content         types.String `tfsdk:"content"`
	AssistantID     types.String `tfsdk:"assistant_id"`
	FileIDs         types.List   `tfsdk:"file_ids"`
	Metadata        types.Map    `tfsdk:"metadata"`
	ObjectID        types.String `tfsdk:"object_id"`
	CreatedAt       types.Int64  `tfsdk:"created_at"`
	ResponseContent types.String `tfsdk:"response_content"`
	WaitForResponse types.Bool   `tfsdk:"wait_for_response"`
	RunID           types.String `tfsdk:"run_id"`
}

func (r *MessageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message"
}

func (r *MessageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage messages within OpenAI Assistant threads.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"thread_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the thread to create the message in.",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the message author. Currently only 'user' is supported.",
				Required:            true,
				// In a future version, this could be validated to only allow 'user'
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The content of the message.",
				Required:            true,
			},
			"assistant_id": schema.StringAttribute{
				MarkdownDescription: "ID of the assistant to use for generating a response. If provided, a run will be created.",
				Optional:            true,
			},
			"file_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A list of file IDs to attach to this message. Files can be uploaded using the `openai_file` resource.",
				Optional:            true,
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Metadata in key-value pairs to attach to the message.",
				Optional:            true,
			},
			"object_id": schema.StringAttribute{
				MarkdownDescription: "The OpenAI ID assigned to this message.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the message was created.",
				Computed:            true,
			},
			"wait_for_response": schema.BoolAttribute{
				MarkdownDescription: "Whether to wait for the assistant to respond after sending the message. Requires assistant_id to be set.",
				Optional:            true,
			},
			"response_content": schema.StringAttribute{
				MarkdownDescription: "If wait_for_response is true, this will contain the assistant's response message content.",
				Computed:            true,
			},
			"run_id": schema.StringAttribute{
				MarkdownDescription: "If assistant_id is provided, this will contain the ID of the run created to generate the response.",
				Computed:            true,
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

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	threadID := plan.ThreadID.ValueString()
	if threadID == "" {
		resp.Diagnostics.AddError(
			"Error Creating Message",
			"Thread ID is required to create a message.",
		)
		return
	}

	// Create the message request
	messageReq := openai.MessageRequest{
		Role:    plan.Role.ValueString(),
		Content: plan.Content.ValueString(),
	}

	// Process file IDs if provided
	if !plan.FileIDs.IsNull() {
		var fileIDs []string
		diags := plan.FileIDs.ElementsAs(ctx, &fileIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		messageReq.FileIds = fileIDs // Changed from FileIDs to FileIds
	}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert from map[string]string to map[string]any
		metadataAny := make(map[string]any)
		for k, v := range metadata {
			metadataAny[k] = v
		}
		messageReq.Metadata = metadataAny
	}

	tflog.Debug(ctx, "Creating message", map[string]interface{}{
		"thread_id": threadID,
		"role":      messageReq.Role,
	})

	// Create the message
	message, err := r.client.OpenAI.CreateMessage(ctx, threadID, messageReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Message",
			fmt.Sprintf("Unable to create message: %s", r.client.HandleError(err)),
		)
		return
	}

	// Run the assistant if required
	var runID string
	var responseContent string

	if !plan.AssistantID.IsNull() {
		assistantID := plan.AssistantID.ValueString()
		waitForResponse := false

		if !plan.WaitForResponse.IsNull() {
			waitForResponse = plan.WaitForResponse.ValueBool()
		}

		// Create a run
		runReq := openai.RunRequest{
			AssistantID: assistantID,
		}

		tflog.Debug(ctx, "Creating run with assistant", map[string]interface{}{
			"thread_id":    threadID,
			"assistant_id": assistantID,
		})

		run, err := r.client.OpenAI.CreateRun(ctx, threadID, runReq)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Run",
				fmt.Sprintf("Unable to create run with assistant: %s", r.client.HandleError(err)),
			)
			return
		}

		runID = run.ID

		// If requested, wait for the assistant to respond
		if waitForResponse {
			tflog.Debug(ctx, "Waiting for assistant response", map[string]interface{}{
				"run_id": runID,
			})

			// Poll for run completion
			completedRun, err := r.waitForRunCompletion(ctx, threadID, runID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Waiting for Run Completion",
					fmt.Sprintf("Unable to get assistant response: %s", r.client.HandleError(err)),
				)
				return
			}

			if completedRun.Status == "completed" {
				// Get the assistant's response messages
				listMsgReq := openai.ListMessagesRequest{
					Limit: 1,      // Just get the last message
					Order: "desc", // Get most recent first
				}

				messages, err := r.client.OpenAI.ListMessages(ctx, threadID, &listMsgReq)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Retrieving Assistant Response",
						fmt.Sprintf("Unable to get assistant response message: %s", r.client.HandleError(err)),
					)
					return
				}

				if len(messages.Messages) > 0 && messages.Messages[0].Role == "assistant" && len(messages.Messages[0].Content) > 0 {
					// Get text content from the first content item
					responseContent = messages.Messages[0].Content[0].Text.Value
				}
			} else {
				resp.Diagnostics.AddWarning(
					"Assistant Run Did Not Complete Successfully",
					fmt.Sprintf("Run completed with status: %s", completedRun.Status),
				)
			}
		}
	}

	// Update the state
	plan.ID = types.StringValue(message.ID)
	plan.ObjectID = types.StringValue(message.ID)
	plan.CreatedAt = types.Int64Value(int64(message.CreatedAt))

	if runID != "" {
		plan.RunID = types.StringValue(runID)
	}

	if responseContent != "" {
		plan.ResponseContent = types.StringValue(responseContent)
	}

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MessageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MessageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	threadID := state.ThreadID.ValueString()
	messageID := state.ObjectID.ValueString()
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

	// Update the state with the latest values
	state.Role = types.StringValue(message.Role)
	state.CreatedAt = types.Int64Value(int64(message.CreatedAt))

	// Set content from the first content item if available
	if len(message.Content) > 0 && message.Content[0].Type == "text" {
		state.Content = types.StringValue(message.Content[0].Text.Value)
	}

	// Convert file IDs
	if len(message.FileIds) > 0 { // Changed from FileIDs to FileIds
		fileIDsList, diags := types.ListValueFrom(ctx, types.StringType, message.FileIds)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.FileIDs = fileIDsList
	}

	// Convert metadata
	if message.Metadata != nil {
		// Convert from map[string]any to map[string]string for Terraform state
		metadataStr := make(map[string]string)
		for k, v := range message.Metadata {
			if strValue, ok := v.(string); ok {
				metadataStr[k] = strValue
			} else {
				metadataStr[k] = fmt.Sprintf("%v", v)
			}
		}

		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, metadataStr)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Metadata = metadataMap
	}

	// Save into state
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
	messageID := state.ObjectID.ValueString()
	if threadID == "" || messageID == "" {
		resp.Diagnostics.AddError(
			"Error Updating Message",
			"Thread ID and Message ID are required to update a message.",
		)
		return
	}

	// Create the update request
	messageReq := openai.ModifyMessageRequest{}

	// Process metadata if provided
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert from map[string]string to map[string]any
		metadataAny := make(map[string]any)
		for k, v := range metadata {
			metadataAny[k] = v
		}
		messageReq.Metadata = metadataAny
	}

	tflog.Debug(ctx, "Updating message", map[string]interface{}{
		"thread_id":  threadID,
		"message_id": messageID,
	})

	// Update the message
	message, err := r.client.OpenAI.ModifyMessage(ctx, threadID, messageID, messageReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Message",
			fmt.Sprintf("Unable to update message: %s", r.client.HandleError(err)),
		)
		return
	}

	// Keep content from original request since API doesn't return it in modify response
	plan.Content = state.Content
	plan.Role = state.Role
	plan.ObjectID = types.StringValue(message.ID)
	plan.CreatedAt = types.Int64Value(int64(message.CreatedAt))

	// Preserve response content
	plan.ResponseContent = state.ResponseContent
	plan.RunID = state.RunID

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *MessageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MessageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	threadID := state.ThreadID.ValueString()
	messageID := state.ObjectID.ValueString()
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_id"), idParts[1])...)
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
			return run, nil
		case "queued", "in_progress", "requires_action":
			// Continue waiting
			time.Sleep(pollingInterval)
		default:
			// Unknown status
			return run, fmt.Errorf("run has unknown status: %s", run.Status)
		}
	}

	return nil, fmt.Errorf("run did not complete within the timeout period")
}
