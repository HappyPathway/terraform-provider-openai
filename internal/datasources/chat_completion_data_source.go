package datasources

import (
	"context"
	"fmt"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sashabaranov/go-openai"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ChatCompletionDataSource{}

func NewChatCompletionDataSource() datasource.DataSource {
	return &ChatCompletionDataSource{}
}

// ChatCompletionDataSource defines the data source implementation.
type ChatCompletionDataSource struct {
	client *client.Client
}

// ChatCompletionDataSourceModel describes the data source data model.
type ChatCompletionDataSourceModel struct {
	ID               types.String  `tfsdk:"id"`
	Model            types.String  `tfsdk:"model"`
	Messages         types.List    `tfsdk:"messages"`
	Temperature      types.Float64 `tfsdk:"temperature"`
	TopP             types.Float64 `tfsdk:"top_p"`
	N                types.Int64   `tfsdk:"n"`
	Stream           types.Bool    `tfsdk:"stream"`
	MaxTokens        types.Int64   `tfsdk:"max_tokens"`
	PresencePenalty  types.Float64 `tfsdk:"presence_penalty"`
	FrequencyPenalty types.Float64 `tfsdk:"frequency_penalty"`
	LogitBias        types.Map     `tfsdk:"logit_bias"`
	User             types.String  `tfsdk:"user"`

	// Output
	ResponseContent types.List   `tfsdk:"response_content"`
	Choices         types.List   `tfsdk:"choices"`
	Usage           types.Object `tfsdk:"usage"`
}

// ChatCompletionMessageModel is used to represent a message in a chat completion request.
type ChatCompletionMessageModel struct {
	Role    types.String `tfsdk:"role"`
	Content types.String `tfsdk:"content"`
}

// ChatCompletionChoiceModel describes the choice output from the API.
type ChatCompletionChoiceModel struct {
	Index        types.Int64  `tfsdk:"index"`
	Message      types.Object `tfsdk:"message"`
	FinishReason types.String `tfsdk:"finish_reason"`
}

// ChatCompletionUsageModel describes the token usage in the API response.
type ChatCompletionUsageModel struct {
	PromptTokens     types.Int64 `tfsdk:"prompt_tokens"`
	CompletionTokens types.Int64 `tfsdk:"completion_tokens"`
	TotalTokens      types.Int64 `tfsdk:"total_tokens"`
}

func (d *ChatCompletionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chat_completion"
}

func (d *ChatCompletionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get completions from OpenAI's chat models during the plan phase.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this data source.",
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "ID of the model to use (e.g., 'gpt-4', 'gpt-3.5-turbo').",
				Required:            true,
			},
			"messages": schema.ListNestedAttribute{
				MarkdownDescription: "A list of messages to generate a response for.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							MarkdownDescription: "The role of the message author. One of 'system', 'user', or 'assistant'.",
							Required:            true,
						},
						"content": schema.StringAttribute{
							MarkdownDescription: "The content of the message.",
							Required:            true,
						},
					},
				},
			},
			"temperature": schema.Float64Attribute{
				MarkdownDescription: "Sampling temperature to use. Higher values like 0.8 mean more randomness, lower values like 0.2 mean more deterministic responses. Default is 1.0.",
				Optional:            true,
			},
			"top_p": schema.Float64Attribute{
				MarkdownDescription: "Alternative to temperature, called nucleus sampling. Default is 1.0.",
				Optional:            true,
			},
			"n": schema.Int64Attribute{
				MarkdownDescription: "Number of completions to generate. Default is 1.",
				Optional:            true,
			},
			"stream": schema.BoolAttribute{
				MarkdownDescription: "Whether to stream responses. Not recommended for Terraform. Default is false.",
				Optional:            true,
			},
			"max_tokens": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of tokens to generate. Default is model-specific.",
				Optional:            true,
			},
			"presence_penalty": schema.Float64Attribute{
				MarkdownDescription: "Penalty for new tokens based on their existence in the text so far. Default is 0.",
				Optional:            true,
			},
			"frequency_penalty": schema.Float64Attribute{
				MarkdownDescription: "Penalty for new tokens based on their frequency in the text so far. Default is 0.",
				Optional:            true,
			},
			"logit_bias": schema.MapAttribute{
				ElementType:         types.Float64Type,
				MarkdownDescription: "Modify the likelihood of specific tokens appearing in the completion.",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "A unique identifier representing your end-user for tracking purposes.",
				Optional:            true,
			},

			// Output attributes
			"response_content": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A simplified list of response content strings, one for each generated completion.",
				Computed:            true,
			},
			"choices": schema.ListNestedAttribute{
				MarkdownDescription: "The list of completion choices the model generated.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"index": schema.Int64Attribute{
							MarkdownDescription: "The index of the choice in the list of choices.",
							Computed:            true,
						},
						"message": schema.SingleNestedAttribute{
							MarkdownDescription: "The message generated as a response.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"role": schema.StringAttribute{
									MarkdownDescription: "The role of the message author, always 'assistant' for responses.",
									Computed:            true,
								},
								"content": schema.StringAttribute{
									MarkdownDescription: "The content of the message.",
									Computed:            true,
								},
							},
						},
						"finish_reason": schema.StringAttribute{
							MarkdownDescription: "The reason why the model stopped generating tokens (e.g., 'stop', 'length').",
							Computed:            true,
						},
					},
				},
			},
			"usage": schema.SingleNestedAttribute{
				MarkdownDescription: "Usage statistics for the completion request.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"prompt_tokens": schema.Int64Attribute{
						MarkdownDescription: "Number of tokens in the prompt.",
						Computed:            true,
					},
					"completion_tokens": schema.Int64Attribute{
						MarkdownDescription: "Number of tokens in the completion.",
						Computed:            true,
					},
					"total_tokens": schema.Int64Attribute{
						MarkdownDescription: "Total number of tokens used (prompt + completion).",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *ChatCompletionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ChatCompletionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ChatCompletionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform messages to OpenAI messages
	messages, diags := convertTerraformMessagesToOpenAI(ctx, data.Messages)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request
	request := openai.ChatCompletionRequest{
		Model:    data.Model.ValueString(),
		Messages: messages,
	}

	// Add optional parameters if provided
	if !data.Temperature.IsNull() {
		request.Temperature = float32(data.Temperature.ValueFloat64())
	}
	if !data.TopP.IsNull() {
		request.TopP = float32(data.TopP.ValueFloat64())
	}
	if !data.N.IsNull() {
		request.N = int(data.N.ValueInt64())
	} else {
		request.N = 1 // Default to 1 completion
	}
	if !data.Stream.IsNull() {
		request.Stream = data.Stream.ValueBool()
	}
	if !data.MaxTokens.IsNull() {
		request.MaxTokens = int(data.MaxTokens.ValueInt64())
	}
	if !data.PresencePenalty.IsNull() {
		request.PresencePenalty = float32(data.PresencePenalty.ValueFloat64())
	}
	if !data.FrequencyPenalty.IsNull() {
		request.FrequencyPenalty = float32(data.FrequencyPenalty.ValueFloat64())
	}
	if !data.LogitBias.IsNull() {
		logitBias := make(map[string]int, len(data.LogitBias.Elements()))
		data.LogitBias.ElementsAs(ctx, &logitBias, false)
		request.LogitBias = logitBias
	}
	if !data.User.IsNull() {
		request.User = data.User.ValueString()
	}

	tflog.Debug(ctx, "Generating chat completion", map[string]interface{}{
		"model": request.Model,
	})

	// Call the API
	response, err := d.client.OpenAI.CreateChatCompletion(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Chat Completion",
			fmt.Sprintf("Unable to create chat completion: %s", d.client.HandleError(err)),
		)
		return
	}

	// Process the response
	data.ID = types.StringValue(response.ID)

	// Extract response content for easier access
	responseContent := make([]string, 0, len(response.Choices))
	for _, choice := range response.Choices {
		responseContent = append(responseContent, choice.Message.Content)
	}

	responseContentList, diags := types.ListValueFrom(ctx, types.StringType, responseContent)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ResponseContent = responseContentList

	// Convert choices to Terraform format
	choices, diags := convertOpenAIChoicesToTerraform(ctx, response.Choices)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Choices = choices

	// Convert usage to Terraform format
	usage, diags := convertOpenAIUsageToTerraform(ctx, response.Usage)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Usage = usage

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Helper function to convert Terraform messages to OpenAI format
func convertTerraformMessagesToOpenAI(ctx context.Context, messagesList types.List) ([]openai.ChatCompletionMessage, diag.Diagnostics) {
	var diags diag.Diagnostics

	if messagesList.IsNull() || messagesList.IsUnknown() {
		return nil, diags
	}

	var tfMessages []ChatCompletionMessageModel
	diags.Append(messagesList.ElementsAs(ctx, &tfMessages, false)...)
	if diags.HasError() {
		return nil, diags
	}

	openAIMessages := make([]openai.ChatCompletionMessage, 0, len(tfMessages))

	for i, message := range tfMessages {
		if message.Role.IsNull() || message.Content.IsNull() {
			diags.AddAttributeError(
				path.Root("messages").AtListIndex(i),
				"Invalid Message",
				"Both role and content must be provided for chat completion messages.",
			)
			continue
		}

		openAIMessage := openai.ChatCompletionMessage{
			Role:    message.Role.ValueString(),
			Content: message.Content.ValueString(),
		}

		openAIMessages = append(openAIMessages, openAIMessage)
	}

	return openAIMessages, diags
}

// Helper function to convert OpenAI choices to Terraform format
func convertOpenAIChoicesToTerraform(ctx context.Context, choices []openai.ChatCompletionChoice) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfChoices := make([]attr.Value, 0, len(choices))

	for _, choice := range choices {
		// Convert the message
		messageAttrValues := map[string]attr.Value{
			"role":    types.StringValue(choice.Message.Role),
			"content": types.StringValue(choice.Message.Content),
		}

		messageObj, d := types.ObjectValue(
			map[string]attr.Type{
				"role":    types.StringType,
				"content": types.StringType,
			},
			messageAttrValues,
		)
		diags.Append(d...)
		if diags.HasError() {
			continue
		}

		// Create the choice object
		choiceAttrValues := map[string]attr.Value{
			"index":         types.Int64Value(int64(choice.Index)),
			"message":       messageObj,
			"finish_reason": types.StringValue(choice.FinishReason),
		}

		choiceObj, d := types.ObjectValue(
			map[string]attr.Type{
				"index":         types.Int64Type,
				"message":       types.ObjectType{AttrTypes: map[string]attr.Type{"role": types.StringType, "content": types.StringType}},
				"finish_reason": types.StringType,
			},
			choiceAttrValues,
		)
		diags.Append(d...)
		if diags.HasError() {
			continue
		}

		tfChoices = append(tfChoices, choiceObj)
	}

	choicesList, d := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"index":         types.Int64Type,
				"message":       types.ObjectType{AttrTypes: map[string]attr.Type{"role": types.StringType, "content": types.StringType}},
				"finish_reason": types.StringType,
			},
		},
		tfChoices,
	)

	diags.Append(d...)
	return choicesList, diags
}

// Helper function to convert OpenAI usage to Terraform format
func convertOpenAIUsageToTerraform(ctx context.Context, usage openai.Usage) (types.Object, diag.Diagnostics) {
	usageAttrValues := map[string]attr.Value{
		"prompt_tokens":     types.Int64Value(int64(usage.PromptTokens)),
		"completion_tokens": types.Int64Value(int64(usage.CompletionTokens)),
		"total_tokens":      types.Int64Value(int64(usage.TotalTokens)),
	}

	return types.ObjectValue(
		map[string]attr.Type{
			"prompt_tokens":     types.Int64Type,
			"completion_tokens": types.Int64Type,
			"total_tokens":      types.Int64Type,
		},
		usageAttrValues,
	)
}
