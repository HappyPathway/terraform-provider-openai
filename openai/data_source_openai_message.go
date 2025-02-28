package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openai/openai-go"
)

func dataSourceOpenAIMessage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenAIMessageRead,
		Schema: map[string]*schema.Schema{
			"thread_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the thread containing the message",
			},
			"message_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the message to retrieve",
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The role of the entity that created the message",
			},
			"content": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of content: text, image_file, or image_url",
						},
						"text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The text content when type is text",
						},
						"image_file": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"detail": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_url": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"detail": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"assistant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the assistant that authored this message",
			},
			"run_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the run associated with message creation",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp for when the message was created",
			},
			"completed_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp for when the message was completed",
			},
			"incomplete_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp for when the message was marked as incomplete",
			},
			"incomplete_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of key-value pairs attached to the message",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the message: in_progress, incomplete, or completed",
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tool": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOpenAIMessageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	threadID := d.Get("thread_id").(string)
	messageID := d.Get("message_id").(string)

	message, err := client.Beta.Threads.Messages.Get(ctx, threadID, messageID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the resource ID to the message ID
	d.SetId(message.ID)

	if err := d.Set("role", message.Role); err != nil {
		return diag.FromErr(err)
	}

	// Process message content
	var content []interface{}
	for _, c := range message.Content {
		contentMap := map[string]interface{}{
			"type": c.Type,
		}

		switch v := c.AsUnion().(type) {
		case *openai.TextContentBlock:
			contentMap["text"] = v.Text.Value
		case *openai.ImageFileContentBlock:
			imageFile := []interface{}{
				map[string]interface{}{
					"file_id": v.ImageFile.FileID,
					"detail":  v.ImageFile.Detail,
				},
			}
			contentMap["image_file"] = imageFile
		case *openai.ImageURLContentBlock:
			imageURL := []interface{}{
				map[string]interface{}{
					"url":    v.ImageURL.URL,
					"detail": v.ImageURL.Detail,
				},
			}
			contentMap["image_url"] = imageURL
		}

		content = append(content, contentMap)
	}

	if err := d.Set("content", content); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("assistant_id", message.AssistantID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("run_id", message.RunID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", message.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("completed_at", message.CompletedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("incomplete_at", message.IncompleteAt); err != nil {
		return diag.FromErr(err)
	}

	if message.IncompleteDetails.Reason != "" {
		incompleteDetails := []interface{}{
			map[string]interface{}{
				"reason": message.IncompleteDetails.Reason,
			},
		}
		if err := d.Set("incomplete_details", incompleteDetails); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("metadata", message.Metadata); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", message.Status); err != nil {
		return diag.FromErr(err)
	}

	// Process message attachments
	var attachments []interface{}
	for _, a := range message.Attachments {
		tools := make([]string, len(a.Tools))
		for i, t := range a.Tools {
			tools[i] = string(t.Type)
		}

		for _, tool := range tools {
			attachment := map[string]interface{}{
				"file_id": a.FileID,
				"tool":    tool,
			}
			attachments = append(attachments, attachment)
		}
	}
	if err := d.Set("attachments", attachments); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
