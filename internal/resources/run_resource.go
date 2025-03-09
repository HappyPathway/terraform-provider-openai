package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sashabaranov/go-openai"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &RunResource{}
var _ resource.ResourceWithImportState = &RunResource{}

func NewRunResource() resource.Resource {
	return &RunResource{}
}

// RunResource defines the resource implementation.
type RunResource struct {
	client *client.Client
}

// RunResourceModel describes the resource data model.
type RunResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ThreadID            types.String `tfsdk:"thread_id"`
	AssistantID         types.String `tfsdk:"assistant_id"`
	Status              types.String `tfsdk:"status"`
	Model               types.String `tfsdk:"model"`
	Instructions        types.String `tfsdk:"instructions"`
	Tools               types.List   `tfsdk:"tools"`
	WaitForCompletion   types.Bool   `tfsdk:"wait_for_completion"`
	PollingInterval     types.String `tfsdk:"polling_interval"`
	Timeout             types.String `tfsdk:"timeout"`
	CreatedAt           types.Int64  `tfsdk:"created_at"`
	ExpiresAt           types.Int64  `tfsdk:"expires_at"`
	StartedAt           types.Int64  `tfsdk:"started_at"`
	CancelledAt         types.Int64  `tfsdk:"cancelled_at"`
	FailedAt            types.Int64  `tfsdk:"failed_at"`
	CompletedAt         types.Int64  `tfsdk:"completed_at"`
	LastError           types.String `tfsdk:"last_error"`
	Steps               types.List   `tfsdk:"steps"`
	RequiredAction      types.Object `tfsdk:"required_action"`
	Metadata            types.Map    `tfsdk:"metadata"`
	MaxPromptTokens     types.Int64  `tfsdk:"max_prompt_tokens"`
	MaxCompletionTokens types.Int64  `tfsdk:"max_completion_tokens"`
	ResponseContent     types.String `tfsdk:"response_content"`
	IncompleteDetails   types.String `tfsdk:"incomplete_details"`
}

func (r *RunResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "openai_run"
}

func (r *RunResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage assistant runs. Runs represent the execution of an assistant on a thread.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the run.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"thread_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the thread to run the assistant on.",
			},
			"assistant_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the assistant to use for this run.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the run (queued, in_progress, completed, requires_action, expired, cancelling, cancelled, failed).",
			},
			"model": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Override the default model used by the assistant.",
			},
			"instructions": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Override the default instructions of the assistant for this run.",
			},
			"tools": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Override the default tools of the assistant for this run.",
			},
			"wait_for_completion": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to wait for the run to complete before marking the resource as created. Defaults to true.",
			},
			"polling_interval": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "How often to poll for run status when wait_for_completion is true. Defaults to 5s.",
			},
			"timeout": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum time to wait for run completion when wait_for_completion is true. Defaults to 10m.",
			},
			"created_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp for when the run was created.",
			},
			"expires_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp for when the run will expire.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"started_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp for when the run was started.",
			},
			"cancelled_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp for when the run was cancelled.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"failed_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp for when the run failed.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"completed_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp for when the run completed.",
			},
			"last_error": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last error message if the run failed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"steps": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The steps taken during the run.",
			},
			"required_action": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"type": types.StringType,
					"submit_tool_outputs": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"tool_calls": types.ListType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"id":   types.StringType,
										"type": types.StringType,
										"function": types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"name":      types.StringType,
												"arguments": types.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
				Computed:            true,
				MarkdownDescription: "Details about any required actions needed to continue the run.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of key-value pairs for run metadata.",
			},
			"max_prompt_tokens": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of tokens to use for prompts in this run. When using File Search tool, recommend setting to at least 50000 for best results.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_completion_tokens": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of tokens to generate in this run. If a completion reaches this limit, the run will terminate with a status of 'incomplete'.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"response_content": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The content of the assistant's response after the run completes.",
			},
			"incomplete_details": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Details about why the run was marked as incomplete, if applicable.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *RunResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert tools to OpenAI format if provided
	var tools []openai.AssistantTool
	if !data.Tools.IsNull() {
		var toolNames []string
		resp.Diagnostics.Append(data.Tools.ElementsAs(ctx, &toolNames, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, toolName := range toolNames {
			tools = append(tools, openai.AssistantTool{Type: openai.AssistantToolType(toolName)})
		}
	}

	// Create run request
	createReq := &client.CreateRunRequest{
		AssistantID: data.AssistantID.ValueString(),
		ThreadID:    data.ThreadID.ValueString(),
	}

	if !data.Model.IsNull() {
		createReq.Model = data.Model.ValueString()
	}

	if !data.Instructions.IsNull() {
		createReq.Instructions = data.Instructions.ValueString()
	}

	if len(tools) > 0 {
		createReq.Tools = tools
	}

	if !data.Metadata.IsNull() {
		metadata := make(map[string]any)
		resp.Diagnostics.Append(data.Metadata.ElementsAs(ctx, &metadata, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Metadata = metadata
	}

	if !data.MaxPromptTokens.IsNull() {
		createReq.MaxPromptTokens = int(data.MaxPromptTokens.ValueInt64())
	}

	if !data.MaxCompletionTokens.IsNull() {
		createReq.MaxCompletionTokens = int(data.MaxCompletionTokens.ValueInt64())
	}

	// Create run
	run, err := r.client.CreateRun(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Run",
			fmt.Sprintf("Unable to create run: %s", err),
		)
		return
	}

	// Store the thread ID and run ID for later use
	threadID := data.ThreadID.ValueString()

	// Update any associated messages with the run ID and assistant ID
	messages, err := r.client.OpenAI.ListMessage(ctx, threadID, nil, nil, nil, nil, nil)
	if err == nil && len(messages.Messages) > 0 {
		// Get the most recent message
		latestMsg := messages.Messages[0]
		metadata := make(map[string]string)
		metadata["run_id"] = run.ID
		metadata["assistant_id"] = data.AssistantID.ValueString()

		_, err = r.client.OpenAI.ModifyMessage(ctx, threadID, latestMsg.ID, metadata)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to update message with run ID: %v", err))
		}
	}

	// Wait for completion if requested
	waitForCompletion := true
	if !data.WaitForCompletion.IsNull() {
		waitForCompletion = data.WaitForCompletion.ValueBool()
	}

	if waitForCompletion {
		pollingInterval := 5 * time.Second
		if !data.PollingInterval.IsNull() {
			d, err := time.ParseDuration(data.PollingInterval.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid Polling Interval",
					fmt.Sprintf("Unable to parse polling interval: %s", err),
				)
				return
			}
			pollingInterval = d
		}

		timeout := 10 * time.Minute
		if !data.Timeout.IsNull() {
			d, err := time.ParseDuration(data.Timeout.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid Timeout",
					fmt.Sprintf("Unable to parse timeout: %s", err),
				)
				return
			}
			timeout = d
		}

		// Poll until completion or timeout
		startTime := time.Now()
		for {
			if time.Since(startTime) > timeout {
				resp.Diagnostics.AddError(
					"Run Timeout",
					fmt.Sprintf("Run did not complete within %s", timeout),
				)
				return
			}

			run, err = r.client.GetRun(ctx, run.ID, threadID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Getting Run",
					fmt.Sprintf("Unable to get run: %s", err),
				)
				return
			}

			switch run.Status {
			case openai.RunStatusCompleted:
				goto DONE
			case openai.RunStatusFailed, openai.RunStatusExpired, openai.RunStatusCancelled:
				resp.Diagnostics.AddError(
					"Run Failed",
					fmt.Sprintf("Run failed with status %s: %v", run.Status, run.LastError),
				)
				return
			case openai.RunStatusRequiresAction:
				resp.Diagnostics.AddError(
					"Run Requires Action",
					"Run requires action but automatic tool outputs are not supported",
				)
				return
			default:
				time.Sleep(pollingInterval)
				continue
			}
		}
	DONE:
	}

	// Update Terraform state
	data.ID = types.StringValue(run.ID)
	data.Status = types.StringValue(string(run.Status))
	data.CreatedAt = types.Int64Value(run.CreatedAt)
	data.ThreadID = types.StringValue(run.ThreadID)
	data.AssistantID = types.StringValue(run.AssistantID)

	// Initialize computed fields with empty/zero values if not set
	if run.ExpiresAt > 0 {
		data.ExpiresAt = types.Int64Value(run.ExpiresAt)
	} else {
		data.ExpiresAt = types.Int64Value(0)
	}

	if run.StartedAt != nil && *run.StartedAt > 0 {
		data.StartedAt = types.Int64Value(*run.StartedAt)
	} else {
		data.StartedAt = types.Int64Value(0)
	}

	if run.CancelledAt != nil && *run.CancelledAt > 0 {
		data.CancelledAt = types.Int64Value(*run.CancelledAt)
	} else {
		data.CancelledAt = types.Int64Value(0)
	}

	if run.FailedAt != nil && *run.FailedAt > 0 {
		data.FailedAt = types.Int64Value(*run.FailedAt)
	} else {
		data.FailedAt = types.Int64Value(0)
	}

	if run.CompletedAt != nil && *run.CompletedAt > 0 {
		data.CompletedAt = types.Int64Value(*run.CompletedAt)
	} else {
		data.CompletedAt = types.Int64Value(0)
	}

	if run.LastError != nil {
		data.LastError = types.StringValue(run.LastError.Message)
	} else {
		data.LastError = types.StringValue("")
	}

	if run.Status == openai.RunStatusIncomplete {
		data.IncompleteDetails = types.StringValue("Run was marked incomplete due to token limit")
	} else {
		data.IncompleteDetails = types.StringValue("")
	}

	// Initialize these fields as empty strings if they're not set
	if data.ResponseContent.IsNull() {
		data.ResponseContent = types.StringValue("")
	}

	// Initialize max tokens fields
	if run.MaxPromptTokens > 0 {
		data.MaxPromptTokens = types.Int64Value(int64(run.MaxPromptTokens))
	} else {
		data.MaxPromptTokens = types.Int64Value(0)
	}

	if run.MaxCompletionTokens > 0 {
		data.MaxCompletionTokens = types.Int64Value(int64(run.MaxCompletionTokens))
	} else {
		data.MaxCompletionTokens = types.Int64Value(0)
	}

	// Initialize required_action as an empty object
	emptyReqAction := make(map[string]attr.Value)
	emptyReqAction["type"] = types.StringValue("")
	emptySubmitOutputs := make(map[string]attr.Value)
	emptyToolCalls, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"type": types.StringType,
			"function": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":      types.StringType,
					"arguments": types.StringType,
				},
			},
		},
	}, []interface{}{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	emptySubmitOutputs["tool_calls"] = emptyToolCalls
	emptyReqAction["submit_tool_outputs"] = types.ObjectValueMust(
		map[string]attr.Type{
			"tool_calls": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":   types.StringType,
						"type": types.StringType,
						"function": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"name":      types.StringType,
								"arguments": types.StringType,
							},
						},
					},
				},
			},
		},
		emptySubmitOutputs,
	)
	data.RequiredAction = types.ObjectValueMust(
		map[string]attr.Type{
			"type": types.StringType,
			"submit_tool_outputs": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"tool_calls": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"id":   types.StringType,
								"type": types.StringType,
								"function": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"name":      types.StringType,
										"arguments": types.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
		emptyReqAction,
	)

	// Get response content from thread messages if run completed successfully
	if run.Status == openai.RunStatusCompleted {
		messages, err := r.client.OpenAI.ListMessage(ctx, data.ThreadID.ValueString(), nil, nil, nil, nil, nil)
		if err == nil && len(messages.Messages) > 0 {
			// Get the latest assistant message
			for _, msg := range messages.Messages {
				if msg.Role == "assistant" {
					if len(msg.Content) > 0 {
						data.ResponseContent = types.StringValue(msg.Content[0].Text.Value)
						break
					}
				}
			}
		}
	} else if run.Status == openai.RunStatusIncomplete {
		data.IncompleteDetails = types.StringValue("Run was marked incomplete due to token limit")
	}

	// Initialize these fields as empty strings if they're not set
	if data.ResponseContent.IsNull() {
		data.ResponseContent = types.StringValue("")
	}
	if data.IncompleteDetails.IsNull() {
		data.IncompleteDetails = types.StringValue("")
	}

	if run.ExpiresAt > 0 {
		data.ExpiresAt = types.Int64Value(run.ExpiresAt)
	}

	if run.StartedAt != nil && *run.StartedAt > 0 {
		data.StartedAt = types.Int64Value(*run.StartedAt)
	}
	if run.CancelledAt != nil && *run.CancelledAt > 0 {
		data.CancelledAt = types.Int64Value(*run.CancelledAt)
	}
	if run.FailedAt != nil && *run.FailedAt > 0 {
		data.FailedAt = types.Int64Value(*run.FailedAt)
	}
	if run.CompletedAt != nil && *run.CompletedAt > 0 {
		data.CompletedAt = types.Int64Value(*run.CompletedAt)
	}

	if run.LastError != nil {
		data.LastError = types.StringValue(run.LastError.Message)
	}

	// For now, we'll set an empty steps list since the API doesn't provide step IDs yet
	stepsList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Steps = stepsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RunResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	run, err := r.client.GetRun(ctx, data.ID.ValueString(), data.ThreadID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Run",
			fmt.Sprintf("Unable to read run: %s", err),
		)
		return
	}

	// Update state with latest data
	data.Status = types.StringValue(string(run.Status))
	data.CreatedAt = types.Int64Value(run.CreatedAt)
	if run.ExpiresAt > 0 {
		data.ExpiresAt = types.Int64Value(run.ExpiresAt)
	} else {
		data.ExpiresAt = types.Int64Value(0)
	}
	if run.StartedAt != nil {
		data.StartedAt = types.Int64Value(*run.StartedAt)
	} else {
		data.StartedAt = types.Int64Value(0)
	}
	if run.CancelledAt != nil {
		data.CancelledAt = types.Int64Value(*run.CancelledAt)
	} else {
		data.CancelledAt = types.Int64Value(0)
	}
	if run.FailedAt != nil {
		data.FailedAt = types.Int64Value(*run.FailedAt)
	} else {
		data.FailedAt = types.Int64Value(0)
	}
	if run.CompletedAt != nil {
		data.CompletedAt = types.Int64Value(*run.CompletedAt)
	} else {
		data.CompletedAt = types.Int64Value(0)
	}
	if run.LastError != nil {
		data.LastError = types.StringValue(run.LastError.Message)
	} else {
		data.LastError = types.StringValue("")
	}

	// Get response content from thread messages if run completed
	if run.Status == openai.RunStatusCompleted {
		messages, err := r.client.OpenAI.ListMessage(ctx, data.ThreadID.ValueString(), nil, nil, nil, nil, nil)
		if err == nil && len(messages.Messages) > 0 {
			// Get the latest assistant message
			for _, msg := range messages.Messages {
				if msg.Role == "assistant" {
					if len(msg.Content) > 0 {
						data.ResponseContent = types.StringValue(msg.Content[0].Text.Value)
						break
					}
				}
			}
		}
	}

	// Set incomplete details if run was incomplete
	if run.Status == openai.RunStatusIncomplete {
		data.IncompleteDetails = types.StringValue("Run was marked incomplete due to token limit")
	}

	// Initialize these fields as empty strings if they're not set
	if data.ResponseContent.IsNull() {
		data.ResponseContent = types.StringValue("")
	}
	if data.IncompleteDetails.IsNull() {
		data.IncompleteDetails = types.StringValue("")
	}

	// Create empty steps list since the API doesn't expose steps
	steps, diags := types.ListValueFrom(ctx, types.StringType, []string{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Steps = steps

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Runs cannot be updated after creation",
	)
}

func (r *RunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RunResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CancelRun(ctx, data.ID.ValueString(), data.ThreadID.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error Cancelling Run",
			fmt.Sprintf("Unable to cancel run: %s", err),
		)
		return
	}
}

func (r *RunResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
