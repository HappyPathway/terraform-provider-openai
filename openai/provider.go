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

// Config holds the provider configuration
type Config struct {
	APIKey     string
	OrgID      string
	RetryMax   int
	RetryDelay time.Duration
	Timeout    time.Duration
	Client     *openaiapi.Client
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENAI_API_KEY", nil),
			},
			"org_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENAI_ORG_ID", nil),
			},
			"retry_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"openai_thread":    resourceOpenAIThread(),
			"openai_assistant": resourceOpenAIAssistant(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"openai_model":            dataSourceOpenAIModel(),
			"openai_file":             dataSourceOpenAIFile(),
			"openai_models":           dataSourceOpenAIModels(),
			"openai_fine_tuned_model": dataSourceOpenAIFineTunedModel(),
			"openai_assistant":        dataSourceOpenAIAssistant(),
			"openai_message":          dataSourceOpenAIMessage(),
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

	retryMax := d.Get("retry_max").(int)
	timeout := time.Duration(d.Get("timeout").(int)) * time.Second
	orgID := d.Get("org_id").(string)

	config := &Config{
		APIKey:   apiKey,
		OrgID:    orgID,
		RetryMax: retryMax,
		Timeout:  timeout,
	}

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithMaxRetries(retryMax),
		option.WithRequestTimeout(timeout),
		option.WithHeader("OpenAI-Beta", "assistants=v2"),
	}

	if orgID != "" {
		opts = append(opts, option.WithOrganization(orgID))
	}

	config.Client = openaiapi.NewClient(opts...)

	return config, nil
}
