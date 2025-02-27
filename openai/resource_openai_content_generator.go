package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceOpenAIContentGenerator returns the resource for content generation
func ResourceOpenAIContentGenerator() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOpenAIContentGeneratorCreate,
		ReadContext:   ResourceOpenAIContentGeneratorRead,
		DeleteContext: ResourceOpenAIContentGeneratorDelete,
		Schema: map[string]*schema.Schema{
			"model": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The model to use for content generation (e.g. gpt-4)",
			},
			"messages": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"system", "user", "assistant", "function"}, false),
							Description:  "The role of the message author (system, user, assistant, or function)",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The content of the message",
						},
					},
				},
				Description: "The messages to send to the model",
			},
			"response_format": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"text", "json_object"}, false),
							Description:  "The format to return the response in (text or json_object)",
						},
						"schema": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "JSON schema that response should conform to when type is json_object",
						},
					},
				},
			},
			"temperature": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Default:     1.0,
				ForceNew:    true,
				Description: "Sampling temperature between 0 and 2",
			},
			// Computed values
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated content",
			},
			"raw_response": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full raw response from the API",
			},
			"usage": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Usage statistics for the completion request",
			},
		},
	}
}

// ResourceOpenAIContentGeneratorCreate creates a new content generation
func ResourceOpenAIContentGeneratorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)
	messages := make([]testutil.ChatCompletionMessage, 0)
	for _, msg := range d.Get("messages").([]interface{}) {
		msgMap := msg.(map[string]interface{})
		messages = append(messages, testutil.ChatCompletionMessage{
			Role:    msgMap["role"].(string),
			Content: msgMap["content"].(string),
		})
	}

	req := &testutil.CreateChatCompletionRequest{
		Model:    d.Get("model").(string),
		Messages: messages,
	}

	if v, ok := d.GetOk("temperature"); ok {
		temp := float32(v.(float64))
		req.Temperature = temp
	}

	if v, ok := d.GetOk("response_format"); ok {
		formats := v.([]interface{})
		if len(formats) > 0 {
			format := formats[0].(map[string]interface{})
			responseFormat := make(map[string]string)
			responseFormat["type"] = format["type"].(string)
			if schema, ok := format["schema"]; ok && schema.(string) != "" {
				responseFormat["schema"] = schema.(string)
			}
			req.ResponseFormat = responseFormat
		}
	}

	completion, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating chat completion: %v", err))
	}

	// Set ID to a hash of the inputs
	d.SetId(fmt.Sprintf("%s-%s", req.Model, messages[0].Content))

	if len(completion.Choices) > 0 {
		d.Set("content", completion.Choices[0].Message.Content)
	}

	rawResponse, err := json.Marshal(completion)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error marshaling response: %v", err))
	}
	d.Set("raw_response", string(rawResponse))

	usage := map[string]interface{}{
		"completion_tokens": completion.Usage.CompletionTokens,
		"prompt_tokens":     completion.Usage.PromptTokens,
		"total_tokens":      completion.Usage.TotalTokens,
	}
	d.Set("usage", usage)

	return nil
}

// ResourceOpenAIContentGeneratorRead reads the content generation state
func ResourceOpenAIContentGeneratorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Content generations are stateless and can't be retrieved after creation
	return nil
}

// ResourceOpenAIContentGeneratorDelete deletes the content generation state
func ResourceOpenAIContentGeneratorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Content generations are stateless and don't need explicit deletion
	return nil
}
