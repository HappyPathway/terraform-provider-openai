package openai

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
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

// NewClientFromConfig creates a new OpenAI client from the configuration
func NewClientFromConfig(config *Config) *openaiapi.Client {
	opts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
		option.WithMaxRetries(config.RetryMax),
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
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The content of the message.",
						},
						"role": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "assistant"}, false),
							Description:  "The role of the entity that is creating the message. Either 'user' or 'assistant'.",
						},
						"file_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of File IDs to attach to the message.",
						},
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Set of key-value pairs that can be attached to the message.",
						},
					},
				},
				Description: "A list of messages to start the thread with.",
			},
			"tool_resources": {
				Type:     schema.TypeList,
				Optional: true,
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
				Description: "Configuration for the tools available in this thread.",
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

	if v, ok := d.GetOk("messages"); ok {
		messages := make([]openaiapi.BetaThreadNewParamsMessage, len(v.([]interface{})))
		messagesList := v.([]interface{})

		for i, msg := range messagesList {
			msgMap := msg.(map[string]interface{})
			
			content := []openaiapi.MessageContentPartParamUnion{
				openaiapi.TextContentBlockParam{
					Text: openaiapi.String(msgMap["content"].(string)),
					Type: openaiapi.F(openaiapi.TextContentBlockParamTypeText),
				},
			}

			role := openaiapi.BetaThreadNewParamsMessagesRoleUser
			if r, ok := msgMap["role"].(string); ok {
				if r == "assistant" {
					role = openaiapi.BetaThreadNewParamsMessagesRoleAssistant
				}
			}

			message := openaiapi.BetaThreadNewParamsMessage{
				Content: openaiapi.F(content),
				Role:    openaiapi.F(role),
			}

			if fileIDs, ok := msgMap["file_ids"].([]interface{}); ok && len(fileIDs) > 0 {
				ids := make([]string, len(fileIDs))
				for j, id := range fileIDs {
					ids[j] = id.(string)
				}
				attachments := make([]openaiapi.BetaThreadNewParamsMessagesAttachment, len(ids))
				for j, id := range ids {
					attachments[j] = openaiapi.BetaThreadNewParamsMessagesAttachment{
						FileID: openaiapi.F(id),
						Tools: openaiapi.F([]openaiapi.BetaThreadNewParamsMessagesAttachmentsToolUnion{
							openaiapi.CodeInterpreterToolParam{
								Type: openaiapi.F(openaiapi.CodeInterpreterToolTypeCodeInterpreter),
							},
						}),
					}
				}
				message.Attachments = openaiapi.F(attachments)
			}

			if metadata, ok := msgMap["metadata"].(map[string]interface{}); ok {
				meta := shared.MetadataParam{}
				for key, value := range metadata {
					meta[key] = value.(string)
				}
				message.Metadata = openaiapi.F(meta)
			}

			messages[i] = message
		}

		params.Messages = openaiapi.F(messages)
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := shared.MetadataParam{}
		metaMap := v.(map[string]interface{})
		for key, value := range metaMap {
			metadata[key] = value.(string)
		}
		params.Metadata = openaiapi.F(metadata)
	}

	if v, ok := d.GetOk("tool_resources"); ok {
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

	if d.HasChange("metadata") {
		metadata := shared.MetadataParam{}
		metaMap := d.Get("metadata").(map[string]interface{})
		for key, value := range metaMap {
			metadata[key] = value.(string)
		}
		params.Metadata = openaiapi.F(metadata)
	}

	if d.HasChange("tool_resources") {
		toolResources := d.Get("tool_resources").([]interface{})[0].(map[string]interface{})
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