package openai

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The OpenAI API key to use for API requests.",
				DefaultFunc: schema.EnvDefaultFunc("OPENAI_API_KEY", nil),
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OpenAI organization ID to use for API requests.",
				DefaultFunc: schema.EnvDefaultFunc("OPENAI_ORGANIZATION_ID", nil),
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The base URL to use for API requests. Defaults to https://api.openai.com/v1",
			},
			"retry_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     defaultRetryMax,
				Description: "Maximum number of retries for API requests",
			},
			"retry_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     int(defaultRetryDelay.Seconds()),
				Description: "Delay between retries in seconds",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     int(defaultHTTPTimeout.Seconds()),
				Description: "Timeout for API requests in seconds",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"openai_model":  dataSourceOpenAIModel(),
			"openai_models": dataSourceOpenAIModels(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"openai_file":              resourceOpenAIFile(),
			"openai_assistant":         resourceOpenAIAssistant(),
			"openai_fine_tuning_job":   resourceOpenAIFineTuningJob(),
			"openai_completion":        resourceOpenAICompletion(),
			"openai_embedding":         resourceOpenAIEmbedding(),
			"openai_image_generation":  resourceOpenAIImageGeneration(),
			"openai_moderation":        resourceOpenAIModeration(),
			"openai_speech":            resourceOpenAISpeech(),
			"openai_transcription":     resourceOpenAITranscription(),
			"openai_translation":       resourceOpenAITranslation(),
			"openai_content_generator": ResourceOpenAIContentGenerator(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiKey := d.Get("api_key")
	if apiKey == nil {
		return nil, diag.Errorf("api_key is required")
	}

	config := ClientConfig{
		APIKey:     apiKey.(string),
		RetryMax:   d.Get("retry_max").(int),
		RetryDelay: time.Duration(d.Get("retry_delay").(int)) * time.Second,
		Timeout:    time.Duration(d.Get("timeout").(int)) * time.Second,
	}

	if v, ok := d.GetOk("organization"); ok {
		config.Organization = v.(string)
	}

	if v, ok := d.GetOk("base_url"); ok {
		config.BaseURL = v.(string)
	}

	client := NewClientWithConfig(config)
	return client, diags
}
