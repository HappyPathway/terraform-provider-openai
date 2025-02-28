package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

// NewClientFromConfig creates a new OpenAI client from the configuration
func NewClientFromConfig(config *Config) *openaiapi.Client {
	opts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
		option.WithMaxRetries(config.RetryMax),
		option.WithRequestTimeout(config.Timeout),
	}

	if config.OrgID != "" {
		opts = append(opts, option.WithOrganization(config.OrgID))
	}

	return openaiapi.NewClient(opts...)
}

func resourceOpenAIThread() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIThreadCreate,
		ReadContext:   resourceOpenAIThreadRead,
		UpdateContext: resourceOpenAIThreadUpdate,
		DeleteContext: resourceOpenAIThreadDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of key-value pairs that can be attached to the thread.",
			},
			"messages": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "assistant"}, false),
							Description:  "The role of the entity that is creating the message.",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The content of the message.",
						},
					},
				},
				Description: "Initial messages to create with the thread. Can only be set during creation.",
			},
			"tool_resources": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code_interpreter": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "List of File IDs for the code interpreter tool.",
									},
								},
							},
							Description: "Configuration for the code interpreter tool.",
						},
						"file_search": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vector_store_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "List of vector store IDs for the file search tool.",
									},
								},
							},
							Description: "Configuration for the file search tool.",
						},
					},
				},
				Description: "Configuration for the tools available in this thread. Can only be set during creation.",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp (in seconds) for when the thread was created.",
			},
		},
	}
}

func resourceOpenAIThreadCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	params := &openaiapi.BetaThreadNewParams{}

	// Convert metadata
	if v, ok := d.GetOk("metadata"); ok {
		metadata := shared.MetadataParam{}
		for key, value := range v.(map[string]interface{}) {
			metadata[key] = value.(string)
		}
		params.Metadata = openaiapi.F(metadata)
	}

	// Convert messages for thread creation
	if v, ok := d.GetOk("messages"); ok {
		messagesList := v.([]interface{})
		messages := make([]openaiapi.BetaThreadNewParamsMessage, len(messagesList))

		for i, msg := range messagesList {
			msgMap := msg.(map[string]interface{})

			content := []openaiapi.MessageContentPartParamUnion{
				openaiapi.TextContentBlockParam{
					Type: openaiapi.F(openaiapi.TextContentBlockParamTypeText),
					Text: openaiapi.F(msgMap["content"].(string)),
				},
			}

			role := openaiapi.BetaThreadNewParamsMessagesRole(msgMap["role"].(string))

			messages[i] = openaiapi.BetaThreadNewParamsMessage{
				Role:    openaiapi.F(role),
				Content: openaiapi.F(content),
			}
		}

		params.Messages = openaiapi.F(messages)
	}

	// Convert tool resources
	if v, ok := d.GetOk("tool_resources"); ok {
		if len(v.([]interface{})) > 0 {
			toolResources := v.([]interface{})[0].(map[string]interface{})
			tr := openaiapi.BetaThreadNewParamsToolResources{}

			if ci, ok := toolResources["code_interpreter"]; ok && len(ci.([]interface{})) > 0 {
				ciConfig := ci.([]interface{})[0].(map[string]interface{})
				if fileIDs, ok := ciConfig["file_ids"]; ok {
					fileIDsList := make([]string, len(fileIDs.([]interface{})))
					for i, id := range fileIDs.([]interface{}) {
						fileIDsList[i] = id.(string)
					}
					tr.CodeInterpreter = openaiapi.F(openaiapi.BetaThreadNewParamsToolResourcesCodeInterpreter{
						FileIDs: openaiapi.F(fileIDsList),
					})
				}
			}

			if fs, ok := toolResources["file_search"]; ok && len(fs.([]interface{})) > 0 {
				fsConfig := fs.([]interface{})[0].(map[string]interface{})
				if vectorStoreIDs, ok := fsConfig["vector_store_ids"]; ok {
					vectorStoreIDsList := make([]string, len(vectorStoreIDs.([]interface{})))
					for i, id := range vectorStoreIDs.([]interface{}) {
						vectorStoreIDsList[i] = id.(string)
					}
					tr.FileSearch = openaiapi.F(openaiapi.BetaThreadNewParamsToolResourcesFileSearch{
						VectorStoreIDs: openaiapi.F(vectorStoreIDsList),
					})
				}
			}

			params.ToolResources = openaiapi.F(tr)
		}
	}

	thread, err := client.Beta.Threads.New(ctx, *params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(thread.ID)
	return resourceOpenAIThreadRead(ctx, d, m)
}

func resourceOpenAIThreadRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	thread, err := client.Beta.Threads.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("created_at", thread.CreatedAt)

	if thread.Metadata != nil {
		d.Set("metadata", thread.Metadata)
	}

	// Set tool resources if present
	tr := make([]interface{}, 1)
	trMap := make(map[string]interface{})

	if thread.ToolResources.CodeInterpreter.FileIDs != nil {
		ci := make([]interface{}, 1)
		ciMap := make(map[string]interface{})
		ciMap["file_ids"] = thread.ToolResources.CodeInterpreter.FileIDs
		ci[0] = ciMap
		trMap["code_interpreter"] = ci
	}

	if thread.ToolResources.FileSearch.VectorStoreIDs != nil {
		fs := make([]interface{}, 1)
		fsMap := make(map[string]interface{})
		fsMap["vector_store_ids"] = thread.ToolResources.FileSearch.VectorStoreIDs
		fs[0] = fsMap
		trMap["file_search"] = fs
	}

	if len(trMap) > 0 {
		tr[0] = trMap
		d.Set("tool_resources", tr)
	}

	return nil
}

func resourceOpenAIThreadUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	params := &openaiapi.BetaThreadUpdateParams{}

	// Only update metadata and tool resources
	if d.HasChange("metadata") {
		metadata := shared.MetadataParam{}
		if v, ok := d.GetOk("metadata"); ok {
			for key, value := range v.(map[string]interface{}) {
				metadata[key] = value.(string)
			}
		}
		params.Metadata = openaiapi.F(metadata)
	}

	if d.HasChange("tool_resources") {
		if v, ok := d.GetOk("tool_resources"); ok && len(v.([]interface{})) > 0 {
			toolResources := v.([]interface{})[0].(map[string]interface{})
			tr := openaiapi.BetaThreadUpdateParamsToolResources{}

			if ci, ok := toolResources["code_interpreter"]; ok && len(ci.([]interface{})) > 0 {
				ciConfig := ci.([]interface{})[0].(map[string]interface{})
				if fileIDs, ok := ciConfig["file_ids"]; ok {
					fileIDsList := make([]string, len(fileIDs.([]interface{})))
					for i, id := range fileIDs.([]interface{}) {
						fileIDsList[i] = id.(string)
					}
					tr.CodeInterpreter = openaiapi.F(openaiapi.BetaThreadUpdateParamsToolResourcesCodeInterpreter{
						FileIDs: openaiapi.F(fileIDsList),
					})
				}
			}

			if fs, ok := toolResources["file_search"]; ok && len(fs.([]interface{})) > 0 {
				fsConfig := fs.([]interface{})[0].(map[string]interface{})
				if vectorStoreIDs, ok := fsConfig["vector_store_ids"]; ok {
					vectorStoreIDsList := make([]string, len(vectorStoreIDs.([]interface{})))
					for i, id := range vectorStoreIDs.([]interface{}) {
						vectorStoreIDsList[i] = id.(string)
					}
					tr.FileSearch = openaiapi.F(openaiapi.BetaThreadUpdateParamsToolResourcesFileSearch{
						VectorStoreIDs: openaiapi.F(vectorStoreIDsList),
					})
				}
			}

			params.ToolResources = openaiapi.F(tr)
		}
	}

	_, err := client.Beta.Threads.Update(ctx, d.Id(), *params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOpenAIThreadRead(ctx, d, m)
}

func resourceOpenAIThreadDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	_, err := client.Beta.Threads.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
