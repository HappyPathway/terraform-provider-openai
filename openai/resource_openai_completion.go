package openai

import (
	"context"
	"fmt"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpenAICompletion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAICompletionCreate,
		ReadContext:   resourceOpenAICompletionRead,
		DeleteContext: resourceOpenAICompletionDelete,
		Schema: map[string]*schema.Schema{
			"model": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prompt": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"best_of": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"echo": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"frequency_penalty": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"logit_bias": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeFloat,
				},
			},
			"max_tokens": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"n": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Default:      1,
			},
			"presence_penalty": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"seed": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"stop": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"temperature": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"top_p": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed fields from API response
			"choices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"finish_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"usage": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceOpenAICompletionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	req := &testutil.CreateCompletionRequest{
		Model:  d.Get("model").(string),
		Prompt: d.Get("prompt").(string),
	}

	if v, ok := d.GetOk("best_of"); ok {
		req.BestOf = v.(int)
	}
	if v, ok := d.GetOk("echo"); ok {
		req.Echo = v.(bool)
	}
	if v, ok := d.GetOk("frequency_penalty"); ok {
		req.FrequencyPenalty = v.(float64)
	}
	if v, ok := d.GetOk("max_tokens"); ok {
		req.MaxTokens = v.(int)
	}
	if v, ok := d.GetOk("n"); ok {
		req.N = v.(int)
	}
	if v, ok := d.GetOk("presence_penalty"); ok {
		req.PresencePenalty = v.(float64)
	}
	if v, ok := d.GetOk("seed"); ok {
		req.Seed = v.(int)
	}
	if v, ok := d.GetOk("suffix"); ok {
		req.Suffix = v.(string)
	}
	if v, ok := d.GetOk("temperature"); ok {
		req.Temperature = float32(v.(float64))
	}
	if v, ok := d.GetOk("top_p"); ok {
		req.TopP = v.(float64)
	}
	if v, ok := d.GetOk("user"); ok {
		req.User = v.(string)
	}

	// Handle logit_bias map
	if v, ok := d.GetOk("logit_bias"); ok {
		logitBias := make(map[string]float64)
		for k, v := range v.(map[string]interface{}) {
			logitBias[k] = v.(float64)
		}
		req.LogitBias = logitBias
	}

	// Handle stop list
	if v, ok := d.GetOk("stop"); ok {
		stop := make([]string, 0)
		for _, s := range v.([]interface{}) {
			stop = append(stop, s.(string))
		}
		req.Stop = stop
	}

	completion, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating completion: %v", err))
	}

	// Set ID to a combination of model and prompt
	d.SetId(fmt.Sprintf("%s-%s", d.Get("model").(string), d.Get("prompt").(string)))

	// Set computed values
	choices := make([]interface{}, len(completion.Choices))
	for i, choice := range completion.Choices {
		choiceMap := map[string]interface{}{
			"text":          choice.Text,
			"index":         choice.Index,
			"finish_reason": choice.FinishReason,
		}
		choices[i] = choiceMap
	}
	if err := d.Set("choices", choices); err != nil {
		return diag.FromErr(err)
	}

	usage := map[string]interface{}{
		"completion_tokens": completion.Usage.CompletionTokens,
		"prompt_tokens":     completion.Usage.PromptTokens,
		"total_tokens":      completion.Usage.TotalTokens,
	}
	if err := d.Set("usage", usage); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOpenAICompletionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Completions are stateless and can't be retrieved after creation
	return nil
}

func resourceOpenAICompletionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Completions are stateless and don't need explicit deletion
	return nil
}
