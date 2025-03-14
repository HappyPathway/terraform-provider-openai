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
var _ resource.Resource = &ChatCompletionResource{}
var _ resource.ResourceWithImportState = &ChatCompletionResource{}

func NewChatCompletionResource() resource.Resource {
	return &ChatCompletionResource{}
}

// ChatCompletionResource defines the resource implementation.
type ChatCompletionResource struct {
	client *client.Client
}

// ChatCompletionResourceModel describes the resource data model.
type ChatCompletionResourceModel struct {
	ID              types.String  `tfsdk:"id"`
	Model           types.String  `tfsdk:"model"`
	Messages        types.List    `tfsdk:"messages"`
	Temperature     types.Float64 `tfsdk:"temperature"`
	TopP            types.Float64 `tfsdk:"top_p"`
	N               types.Int64   `tfsdk:"n"`
	MaxTokens       types.Int64   `tfsdk:"max_tokens"`
	ResponseContent types.List    `tfsdk:"response_content"`
	ResponseRole    types.String  `tfsdk:"response_role"`
}

// MessageModel represents a chat message.
type MessageModel struct {
	Role    types.String `tfsdk:"role"`
	Content types.String `tfsdk:"content"`
}

func (r *ChatCompletionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chat_completion"
}

func (r *ChatCompletionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generate a chat completion using OpenAI's GPT models.",
		Blocks: map[string]schema.Block{
			"messages": schema.ListNestedBlock{
				MarkdownDescription: "A list of messages comprising the conversation so far.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							MarkdownDescription: "The role of the message author. Can be 'system', 'user', or 'assistant'.",
							Required:            true,
						},
						"content": schema.StringAttribute{
							MarkdownDescription: "The content of the message.",
							Required:            true,
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "ID of the model to use for completion (e.g., 'gpt-4', 'gpt-3.5-turbo').",
				Required:            true,
			},
			"temperature": schema.Float64Attribute{
				MarkdownDescription: "Sampling temperature between 0 and 2. Higher values like 0.8 make output more random, while lower values like 0.2 make it more focused and deterministic.",
				Optional:            true,
			},
			"top_p": schema.Float64Attribute{
				MarkdownDescription: "An alternative to sampling with temperature, called nucleus sampling. Set this between 0 and 1.",
				Optional:            true,
			},
			"n": schema.Int64Attribute{
				MarkdownDescription: "How many completion choices to generate. Default is 1.",
				Optional:            true,
			},
			"max_tokens": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of tokens to generate.",
				Optional:            true,
			},
			"response_content": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The completion generated by the model. For n=1, this will be a list with a single element.",
				Computed:            true,
			},
			"response_role": schema.StringAttribute{
				MarkdownDescription: "The role of the returned message. Typically 'assistant'.",
				Computed:            true,
			},
		},
	}
}

func (r *ChatCompletionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ChatCompletionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ChatCompletionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan to OpenAI request
	messages, messageDiags := convertMessagesToOpenAI(ctx, plan.Messages)
	resp.Diagnostics.Append(messageDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	chatReq := openai.ChatCompletionRequest{
		Model:    plan.Model.ValueString(),
		Messages: messages,
	}

	// Set optional parameters if they're provided
	if !plan.Temperature.IsNull() {
		chatReq.Temperature = float32(plan.Temperature.ValueFloat64())
	}

	if !plan.TopP.IsNull() {
		chatReq.TopP = float32(plan.TopP.ValueFloat64())
	}

	if !plan.N.IsNull() {
		chatReq.N = int(plan.N.ValueInt64())
	} else {
		chatReq.N = 1 // Default to 1
	}

	if !plan.MaxTokens.IsNull() {
		chatReq.MaxTokens = int(plan.MaxTokens.ValueInt64())
	}

	tflog.Debug(ctx, "Creating chat completion", map[string]interface{}{
		"model": chatReq.Model,
		"n":     chatReq.N,
	})

	// Use ExecuteWithRetry for API call
	result, err := r.client.ExecuteWithRetry(func() (interface{}, error) {
		return r.client.OpenAI.CreateChatCompletion(ctx, chatReq)
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Chat Completion",
			fmt.Sprintf("Unable to create chat completion: %s", r.client.HandleError(err)),
		)
		return
	}

	chatCompletion, ok := result.(openai.ChatCompletionResponse)
	if !ok {
		resp.Diagnostics.AddError(
			"Error Creating Chat Completion",
			"Unable to cast API response to ChatCompletionResponse",
		)
		return
	}

	// Store response content as list
	responseContents := []attr.Value{}
	for _, choice := range chatCompletion.Choices {
		responseContents = append(responseContents, types.StringValue(choice.Message.Content))
	}

	responseContentsList, diags := types.ListValue(types.StringType, responseContents)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set computed values
	plan.ID = types.StringValue(fmt.Sprintf("chat-%s-%d", plan.Model.ValueString(), chatCompletion.Created))
	plan.ResponseContent = responseContentsList

	if len(chatCompletion.Choices) > 0 {
		plan.ResponseRole = types.StringValue(chatCompletion.Choices[0].Message.Role)
	} else {
		plan.ResponseRole = types.StringValue("assistant")
	}

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ChatCompletionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ChatCompletionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ChatCompletion is stateless - we just keep whatever is in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ChatCompletionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Same as create, since we regenerate the completion
	var plan ChatCompletionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan to OpenAI request
	messages, messageDiags := convertMessagesToOpenAI(ctx, plan.Messages)
	resp.Diagnostics.Append(messageDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	chatReq := openai.ChatCompletionRequest{
		Model:    plan.Model.ValueString(),
		Messages: messages,
	}

	// Set optional parameters if they're provided
	if !plan.Temperature.IsNull() {
		chatReq.Temperature = float32(plan.Temperature.ValueFloat64())
	}

	if !plan.TopP.IsNull() {
		chatReq.TopP = float32(plan.TopP.ValueFloat64())
	}

	if !plan.N.IsNull() {
		chatReq.N = int(plan.N.ValueInt64())
	} else {
		chatReq.N = 1 // Default to 1
	}

	if !plan.MaxTokens.IsNull() {
		chatReq.MaxTokens = int(plan.MaxTokens.ValueInt64())
	}

	tflog.Debug(ctx, "Updating chat completion", map[string]interface{}{
		"model": chatReq.Model,
		"n":     chatReq.N,
	})

	// Use ExecuteWithRetry for API call
	result, err := r.client.ExecuteWithRetry(func() (interface{}, error) {
		return r.client.OpenAI.CreateChatCompletion(ctx, chatReq)
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Chat Completion",
			fmt.Sprintf("Unable to create chat completion: %s", r.client.HandleError(err)),
		)
		return
	}

	chatCompletion, ok := result.(openai.ChatCompletionResponse)
	if !ok {
		resp.Diagnostics.AddError(
			"Error Creating Chat Completion",
			"Unable to cast API response to ChatCompletionResponse",
		)
		return
	}

	// Store response content as list
	responseContents := []attr.Value{}
	for _, choice := range chatCompletion.Choices {
		responseContents = append(responseContents, types.StringValue(choice.Message.Content))
	}

	responseContentsList, diags := types.ListValue(types.StringType, responseContents)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update computed values
	plan.ResponseContent = responseContentsList

	if len(chatCompletion.Choices) > 0 {
		plan.ResponseRole = types.StringValue(chatCompletion.Choices[0].Message.Role)
	} else {
		plan.ResponseRole = types.StringValue("assistant")
	}

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ChatCompletionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No API call needed for deletion since chat completions are stateless
	tflog.Info(ctx, "Deleted chat completion resource", map[string]interface{}{
		"id": req.State.GetAttribute(ctx, path.Root("id"), nil),
	})
}

func (r *ChatCompletionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to convert from Terraform messages to OpenAI messages
func convertMessagesToOpenAI(ctx context.Context, messagesAttr types.List) ([]openai.ChatCompletionMessage, diag.Diagnostics) {
	var diags diag.Diagnostics
	var messages []MessageModel

	if messagesAttr.IsNull() || messagesAttr.IsUnknown() {
		return nil, diags
	}

	diags.Append(messagesAttr.ElementsAs(ctx, &messages, false)...)
	if diags.HasError() {
		return nil, diags
	}

	openaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for i, msg := range messages {
		if msg.Role.IsNull() || msg.Content.IsNull() {
			diags.AddAttributeError(
				path.Root("messages").AtListIndex(i),
				"Invalid Message",
				"Message must have both role and content.",
			)
			continue
		}

		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    msg.Role.ValueString(),
			Content: msg.Content.ValueString(),
		})
	}

	return openaiMessages, diags
}
