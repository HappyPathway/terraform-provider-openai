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
	ID          types.String `tfsdk:"id"`
	ThreadID    types.String `tfsdk:"thread_id"`
	Role        types.String `tfsdk:"role"`
	Content     types.String `tfsdk:"content"`
	FileIDs     types.List   `tfsdk:"file_ids"`
	Metadata    types.Map    `tfsdk:"metadata"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	Object      types.String `tfsdk:"object"`
	AssistantID types.String `tfsdk:"assistant_id"`
	Attachments types.List   `tfsdk:"attachments"`
}

// AttachmentModel represents a file attachment in a message
type AttachmentModel struct {
	FileID types.String `tfsdk:"file_id"`
	Tools  types.List   `tfsdk:"tools"`
}

func (r *MessageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages individual messages within OpenAI Threads. Messages are the building blocks of conversations between users and assistants.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The OpenAI-assigned ID for this message.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"thread_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the thread to add the message to.",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the message author. Can be \"user\" or \"assistant\".",
				Required:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The content of the message.",
				Required:            true,
			},
			"file_ids": schema.ListAttribute{
				MarkdownDescription: "A list of file IDs to attach to the message. These files must already be uploaded with purpose \"assistants\".",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"metadata": schema.MapAttribute{
				MarkdownDescription: "Set of key-value pairs that can be used to store additional information about the message.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the message was created.",
				Computed:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object type, always \"thread.message\".",
				Computed:            true,
			},
			"assistant_id": schema.StringAttribute{
				MarkdownDescription: "If applicable, the ID of the assistant that created the message.",
				Computed:            true,
			},
			"attachments": schema.ListNestedAttribute{
				MarkdownDescription: "A list of files to attach to the message with their associated tools.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"file_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The ID of the file to attach.",
						},
						"tools": schema.ListAttribute{
							Required:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of tools that can use this file. Can be \"code_interpreter\" and/or \"retrieve\".",
						},
					},
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

	// Add file IDs if specified
	if !plan.FileIDs.IsNull() {
		var fileIDs []string
		resp.Diagnostics.Append(plan.FileIDs.ElementsAs(ctx, &fileIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		messageReq.FileIds = fileIDs
	}

	// Add metadata if specified
	if !plan.Metadata.IsNull() {
		metadata := make(map[string]string)
		resp.Diagnostics.Append(plan.Metadata.ElementsAs(ctx, &metadata, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert from map[string]string to map[string]interface{}
		metadataAny := make(map[string]interface{})
		for k, v := range metadata {
			metadataAny[k] = v
		}
		messageReq.Metadata = metadataAny
	}

	// Add attachments if specified
	if !plan.Attachments.IsNull() {
		var attachments []AttachmentModel
		resp.Diagnostics.Append(plan.Attachments.ElementsAs(ctx, &attachments, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		messageAttachments := make([]openai.ThreadAttachment, 0, len(attachments))
		for _, attachment := range attachments {
			var tools []string
			resp.Diagnostics.Append(attachment.Tools.ElementsAs(ctx, &tools, false)...)
			if resp.Diagnostics.HasError() {
				return
			}

			threadTools := make([]openai.ThreadAttachmentTool, 0, len(tools))
			for _, tool := range tools {
				threadTools = append(threadTools, openai.ThreadAttachmentTool{
					Type: tool,
				})
			}

			messageAttachments = append(messageAttachments, openai.ThreadAttachment{
				FileID: attachment.FileID.ValueString(),
				Tools:  threadTools,
			})
		}
		messageReq.Attachments = messageAttachments
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

	// Update the state
	plan.ID = types.StringValue(message.ID)
	plan.Object = types.StringValue(message.Object)
	plan.CreatedAt = types.Int64Value(int64(message.CreatedAt))

	// Ensure assistant_id is always set, even if null
	if message.AssistantID != nil {
		plan.AssistantID = types.StringValue(*message.AssistantID)
	} else {
		plan.AssistantID = types.StringNull()
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
	state.Role = types.StringValue(message.Role)
	state.CreatedAt = types.Int64Value(int64(message.CreatedAt))
	state.Object = types.StringValue(message.Object)

	// Always set assistant_id, even if null
	if message.AssistantID != nil {
		state.AssistantID = types.StringValue(*message.AssistantID)
	} else {
		state.AssistantID = types.StringNull()
	}

	// Set content from the first content item if available
	if len(message.Content) > 0 && message.Content[0].Type == "text" {
		state.Content = types.StringValue(message.Content[0].Text.Value)
	}

	// Convert file IDs
	if len(message.FileIds) > 0 {
		fileIDsList, diags := types.ListValueFrom(ctx, types.StringType, message.FileIds)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.FileIDs = fileIDsList
	} else {
		state.FileIDs = types.ListNull(types.StringType)
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

		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, metadataStr)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Metadata = metadataMap
	} else {
		state.Metadata = types.MapNull(types.StringType)
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

	// Create the update request
	messageReq := openai.Message{
		Content: []openai.MessageContent{{
			Type: "text",
			Text: &openai.MessageText{
				Value: plan.Content.ValueString(),
			},
		}},
	}

	// Update file IDs if changed
	if !plan.FileIDs.IsNull() && !plan.FileIDs.Equal(state.FileIDs) {
		var fileIDs []string
		diags := plan.FileIDs.ElementsAs(ctx, &fileIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		messageReq.FileIds = fileIDs
	}

	// Update metadata if changed
	if !plan.Metadata.IsNull() && !plan.Metadata.Equal(state.Metadata) {
		metadataStr := make(map[string]string)
		diags := plan.Metadata.ElementsAs(ctx, &metadataStr, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert from map[string]string to map[string]interface{}
		metadata := make(map[string]interface{})
		for k, v := range metadataStr {
			metadata[k] = v
		}
		messageReq.Metadata = metadata
	}

	tflog.Debug(ctx, "Updating message", map[string]interface{}{
		"thread_id":  threadID,
		"message_id": messageID,
	})

	// Convert the request to string map as required by ModifyMessage
	modifyRequest := make(map[string]string)

	// Add metadata if present
	for k, v := range messageReq.Metadata {
		if strValue, ok := v.(string); ok {
			modifyRequest[k] = strValue
		} else {
			modifyRequest[k] = fmt.Sprintf("%v", v)
		}
	}

	// Add content if changed
	if len(messageReq.Content) > 0 {
		modifyRequest["content"] = messageReq.Content[0].Text.Value
	}

	// Don't add file_ids as a comma-separated string as it's not supported by the API
	// File IDs should be handled through the messageReq structure directly

	// Update the message
	message, err := r.client.OpenAI.ModifyMessage(ctx, threadID, messageID, modifyRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Message",
			fmt.Sprintf("Unable to update message: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state
	if message.Role != "" {
		plan.Role = types.StringValue(message.Role)
	} else {
		plan.Role = state.Role // Preserve existing role if not returned
	}

	// Update content if it was returned, otherwise keep the plan value
	if len(message.Content) > 0 && message.Content[0].Type == "text" {
		plan.Content = types.StringValue(message.Content[0].Text.Value)
	}

	plan.Object = types.StringValue(message.Object)
	plan.CreatedAt = types.Int64Value(int64(message.CreatedAt))

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

func (r *MessageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message"
}
