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
				Description: "OpenAI API Key",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OpenAI Organization ID",
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
			"openai_content_generator": ResourceOpenAIContentGenerator(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := ClientConfig{
		APIKey:     d.Get("api_key").(string),
		RetryMax:   d.Get("retry_max").(int),
		RetryDelay: time.Duration(d.Get("retry_delay").(int)) * time.Second,
		Timeout:    time.Duration(d.Get("timeout").(int)) * time.Second,
	}

	if v, ok := d.GetOk("organization_id"); ok {
		// Organization ID will be used in a future implementation
		_ = v.(string)
	}

	client := NewClientWithConfig(config)
	return client, nil
}
