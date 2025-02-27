package openai

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Default:     2,
				Description: "Maximum number of retries for API requests",
			},
			"retry_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Delay between retries in seconds",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Timeout for API requests in seconds",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"openai_model":     dataSourceOpenAIModel(),
			"openai_models":    dataSourceOpenAIModels(),
			"openai_assistant": dataSourceOpenAIAssistant(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"openai_file":              resourceOpenAIFile(),
			"openai_assistant":         resourceOpenAIAssistant(),
			"openai_fine_tuning_job":   resourceOpenAIFineTuningJob(),
			"openai_content_generator": ResourceOpenAIContentGenerator(),
			"openai_thread":            resourceOpenAIThread(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)

	// If api_key is not set in provider config, check environment variable
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	// If still no API key, return error
	if apiKey == "" {
		return nil, diag.Errorf("api_key must be provided via configuration or OPENAI_API_KEY environment variable")
	}

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithMaxRetries(d.Get("retry_max").(int)),
		option.WithRequestTimeout(time.Duration(d.Get("timeout").(int)) * time.Second),
	}

	if orgID, ok := d.GetOk("organization_id"); ok {
		opts = append(opts, option.WithOrganization(orgID.(string)))
	}

	client := openaiapi.NewClient(opts...)
	return &Config{Client: client}, nil
}
