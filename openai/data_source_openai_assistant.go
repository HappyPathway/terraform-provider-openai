package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenAIAssistant() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenAIAssistantRead,
		Schema: map[string]*schema.Schema{
			"assistant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the assistant to retrieve",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the assistant",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the assistant",
			},
			"model": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the model used by the assistant",
			},
			"instructions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "System instructions that the assistant uses",
			},
			"tools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of tool (code_interpreter, retrieval, or function)",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the function (for function type tools)",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of what the function does (for function type tools)",
						},
						"parameters": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameters the function accepts in JSON schema format (for function type tools)",
						},
					},
				},
				Description: "The tools that the assistant can use",
			},
			"file_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of file IDs attached to the assistant",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Metadata associated with the assistant",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp for when the assistant was created",
			},
		},
	}
}

func dataSourceOpenAIAssistantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	assistantID := d.Get("assistant_id").(string)

	assistant, err := client.GetAssistant(ctx, assistantID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(assistant.ID)
	d.Set("name", assistant.Name)
	d.Set("description", assistant.Description)
	d.Set("model", assistant.Model)
	d.Set("instructions", assistant.Instructions)
	d.Set("created_at", assistant.CreatedAt)

	tools := make([]interface{}, len(assistant.Tools))
	for i, tool := range assistant.Tools {
		toolMap := map[string]interface{}{
			"type": tool.Type,
		}

		if tool.Type == "function" && tool.Function != nil {
			toolMap["name"] = tool.Function.Name
			toolMap["description"] = tool.Function.Description
			toolMap["parameters"] = tool.Function.Parameters
		}
		tools[i] = toolMap
	}
	d.Set("tools", tools)

	d.Set("file_ids", assistant.FileIDs)
	d.Set("metadata", assistant.Metadata)

	return nil
}
