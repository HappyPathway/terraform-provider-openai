package openai

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

func resourceOpenAIMessage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIMessageCreate,
		ReadContext:   resourceOpenAIMessageRead,
		UpdateContext: resourceOpenAIMessageUpdate,
		DeleteContext: resourceOpenAIMessageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// Try to split the ID - if it contains a slash, it's in the old format
				threadID, messageID := extractIDs(d.Id())
				if threadID != "" && messageID != "" {
					d.Set("thread_id", threadID)
					d.SetId(messageID)
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceOpenAIMessageResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceOpenAIMessageUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"thread_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the thread to create a message for.",
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The role of the message creator. Either 'user' or 'assistant'.",
				ValidateFunc: validation.StringInSlice([]string{"user", "assistant"}, false),
			},
			"content": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The type of content. Currently only supports 'text'.",
							ValidateFunc: validation.StringInSlice([]string{"text"}, false),
						},
						"text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The text content of the message.",
						},
					},
				},
			},
			"file_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of File IDs to attach to the message.",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of key-value pairs that can be attached to the message.",
			},
			"assistant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the assistant that authored this message, if applicable.",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp (in seconds) for when the message was created.",
			},
			"run_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the run associated with the message creation, if applicable.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the message: in_progress, incomplete, or completed.",
			},
		},
	}
}

func resourceOpenAIMessageResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"thread_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Only include fields that existed in v0
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"text": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"file_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceOpenAIMessageUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState["id"] == nil {
		return rawState, nil
	}

	// Check if the ID is in the old format (contains a slash)
	oldID := rawState["id"].(string)
	threadID, messageID := extractIDs(oldID)
	if threadID != "" && messageID != "" {
		// Store just the message ID in the id field
		rawState["id"] = messageID
	}

	return rawState, nil
}

func resourceOpenAIMessageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	threadID := d.Get("thread_id").(string)

	params := &openaiapi.BetaThreadMessageNewParams{
		Role: openaiapi.F(openaiapi.BetaThreadMessageNewParamsRole(d.Get("role").(string))),
	}

	// Convert content
	if v, ok := d.GetOk("content"); ok {
		contentList := v.([]interface{})
		contentParams := make([]openaiapi.MessageContentPartParamUnion, len(contentList))

		for i, content := range contentList {
			contentMap := content.(map[string]interface{})
			if contentMap["type"].(string) == "text" {
				contentParams[i] = openaiapi.TextContentBlockParam{
					Type: openaiapi.F(openaiapi.TextContentBlockParamTypeText),
					Text: openaiapi.F(contentMap["text"].(string)),
				}
			}
		}

		params.Content = openaiapi.F(contentParams)
	}

	// Add file attachments if present
	if v, ok := d.GetOk("file_ids"); ok {
		fileIDs := make([]openaiapi.BetaThreadMessageNewParamsAttachment, 0)
		for _, id := range v.([]interface{}) {
			fileIDs = append(fileIDs, openaiapi.BetaThreadMessageNewParamsAttachment{
				FileID: openaiapi.F(id.(string)),
				Tools: openaiapi.F([]openaiapi.BetaThreadMessageNewParamsAttachmentsToolUnion{
					openaiapi.CodeInterpreterToolParam{
						Type: openaiapi.F(openaiapi.CodeInterpreterToolTypeCodeInterpreter),
					},
				}),
			})
		}
		params.Attachments = openaiapi.F(fileIDs)
	}

	// Add metadata if present
	if v, ok := d.GetOk("metadata"); ok {
		metadata := shared.MetadataParam{}
		for key, value := range v.(map[string]interface{}) {
			metadata[key] = value.(string)
		}
		params.Metadata = openaiapi.F(metadata)
	}

	message, err := client.Beta.Threads.Messages.New(ctx, threadID, *params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Just store the message ID, since thread_id is already stored separately
	d.SetId(message.ID)
	return resourceOpenAIMessageRead(ctx, d, m)
}

func resourceOpenAIMessageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	threadID := d.Get("thread_id").(string)
	messageID := d.Id() // This is just the message ID now

	message, err := client.Beta.Threads.Messages.Get(ctx, threadID, messageID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the computed fields
	d.Set("thread_id", threadID)
	d.Set("role", message.Role)
	d.Set("assistant_id", message.AssistantID)
	d.Set("created_at", message.CreatedAt)
	d.Set("run_id", message.RunID)
	d.Set("status", message.Status)

	if message.Metadata != nil {
		d.Set("metadata", message.Metadata)
	}

	// Convert message content back to schema format
	content := make([]map[string]interface{}, 0)
	for _, part := range message.Content {
		contentMap := map[string]interface{}{
			"type": part.Type,
		}

		switch v := part.AsUnion().(type) {
		case *openaiapi.TextContentBlock:
			contentMap["text"] = v.Text.Value
		}

		content = append(content, contentMap)
	}
	d.Set("content", content)

	return nil
}

func resourceOpenAIMessageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	threadID := d.Get("thread_id").(string)
	messageID := d.Id()

	// Only metadata can be updated
	if d.HasChange("metadata") {
		metadata := shared.MetadataParam{}
		if v, ok := d.GetOk("metadata"); ok {
			for key, value := range v.(map[string]interface{}) {
				metadata[key] = value.(string)
			}
		}

		_, err := client.Beta.Threads.Messages.Update(
			ctx,
			threadID,
			messageID,
			openaiapi.BetaThreadMessageUpdateParams{
				Metadata: openaiapi.F(metadata),
			},
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceOpenAIMessageRead(ctx, d, m)
}

func resourceOpenAIMessageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	threadID := d.Get("thread_id").(string)
	messageID := d.Id()

	_, err := client.Beta.Threads.Messages.Delete(ctx, threadID, messageID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// Helper function to extract thread_id and message_id from the resource ID
func extractIDs(id string) (threadID, messageID string) {
	parts := strings.Split(id, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
