package openai

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

func resourceOpenAIAssistant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIAssistantCreate,
		ReadContext:   resourceOpenAIAssistantRead,
		UpdateContext: resourceOpenAIAssistantUpdate,
		DeleteContext: resourceOpenAIAssistantDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the assistant",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the assistant",
			},
			"model": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the model to use",
			},
			"instructions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "System instructions that the assistant uses",
			},
			"tools": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of tool (code_interpreter, retrieval, or function)",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the function (required when type is function)",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A description of what the function does (required when type is function)",
						},
						"parameters": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The parameters the function accepts in JSON schema format (required when type is function)",
						},
					},
					Description: "The tools that the assistant can use",
				},
			},
			"file_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of file IDs attached to the assistant",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Metadata to associate with the assistant",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp for when the assistant was created",
			},
		},
	}
}

func resourceOpenAIAssistantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	tools := []openaiapi.AssistantToolUnionParam{}
	if v, ok := d.GetOk("tools"); ok {
		toolsList := v.([]interface{})
		for _, tool := range toolsList {
			toolMap := tool.(map[string]interface{})
			toolType := toolMap["type"].(string)

			switch toolType {
			case "code_interpreter":
				tools = append(tools, openaiapi.CodeInterpreterToolParam{
					Type: openaiapi.F(openaiapi.CodeInterpreterToolTypeCodeInterpreter),
				})
			case "retrieval":
				tools = append(tools, openaiapi.FileSearchToolParam{
					Type: openaiapi.F(openaiapi.FileSearchToolTypeFileSearch),
				})
			case "function":
				var params shared.FunctionParameters
				if err := json.Unmarshal([]byte(toolMap["parameters"].(string)), &params); err != nil {
					return diag.FromErr(err)
				}
				tools = append(tools, openaiapi.FunctionToolParam{
					Type: openaiapi.F(openaiapi.FunctionToolTypeFunction),
					Function: openaiapi.F(openaiapi.FunctionDefinitionParam{
						Name:        openaiapi.F(toolMap["name"].(string)),
						Description: openaiapi.F(toolMap["description"].(string)),
						Parameters:  openaiapi.F(params),
					}),
				})
			}
		}
	}

	fileIDs := []string{}
	if v, ok := d.GetOk("file_ids"); ok {
		fileIDsRaw := v.([]interface{})
		for _, id := range fileIDsRaw {
			fileIDs = append(fileIDs, id.(string))
		}
	}

	metadata := shared.MetadataParam{}
	if v, ok := d.GetOk("metadata"); ok {
		metaMap := v.(map[string]interface{})
		for key, value := range metaMap {
			metadata[key] = value.(string)
		}
	}

	toolResources := openaiapi.BetaAssistantNewParamsToolResources{}
	if len(fileIDs) > 0 {
		toolResources.CodeInterpreter = openaiapi.F(openaiapi.BetaAssistantNewParamsToolResourcesCodeInterpreter{
			FileIDs: openaiapi.F(fileIDs),
		})
	}

	assistant, err := client.Beta.Assistants.New(ctx, openaiapi.BetaAssistantNewParams{
		Name:          openaiapi.F(d.Get("name").(string)),
		Model:         openaiapi.F(openaiapi.ChatModel(d.Get("model").(string))),
		Description:   openaiapi.F(d.Get("description").(string)),
		Instructions:  openaiapi.F(d.Get("instructions").(string)),
		Tools:         openaiapi.F(tools),
		Metadata:      openaiapi.F(metadata),
		ToolResources: openaiapi.F(toolResources),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(assistant.ID)
	return resourceOpenAIAssistantRead(ctx, d, m)
}

func resourceOpenAIAssistantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	assistant, err := client.Beta.Assistants.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", assistant.Name)
	d.Set("description", assistant.Description)
	d.Set("model", assistant.Model)
	d.Set("instructions", assistant.Instructions)
	d.Set("created_at", assistant.CreatedAt)

	tools := make([]interface{}, len(assistant.Tools))
	for i, tool := range assistant.Tools {
		toolMap := map[string]interface{}{
			"type": string(tool.Type),
		}

		// Handle each tool type appropriately
		switch unionTool := tool.AsUnion().(type) {
		case *openaiapi.FunctionTool:
			if unionTool.Function.Name != "" {
				toolMap["name"] = unionTool.Function.Name
				toolMap["description"] = unionTool.Function.Description
				if unionTool.Function.Parameters != nil {
					paramsJson, err := json.Marshal(unionTool.Function.Parameters)
					if err != nil {
						return diag.FromErr(err)
					}
					toolMap["parameters"] = string(paramsJson)
				}
			}
		case *openaiapi.CodeInterpreterTool:
			// code_interpreter tools don't have additional fields
		case *openaiapi.FileSearchTool:
			// retrieval/file_search tools don't have additional fields
		}
		tools[i] = toolMap
	}
	d.Set("tools", tools)

	if assistant.ToolResources.CodeInterpreter.FileIDs != nil {
		d.Set("file_ids", assistant.ToolResources.CodeInterpreter.FileIDs)
	}
	d.Set("metadata", assistant.Metadata)

	return nil
}

func resourceOpenAIAssistantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	tools := []openaiapi.AssistantToolUnionParam{}
	if v, ok := d.GetOk("tools"); ok {
		toolsList := v.([]interface{})
		for _, tool := range toolsList {
			toolMap := tool.(map[string]interface{})
			toolType := toolMap["type"].(string)

			switch toolType {
			case "code_interpreter":
				tools = append(tools, openaiapi.CodeInterpreterToolParam{
					Type: openaiapi.F(openaiapi.CodeInterpreterToolTypeCodeInterpreter),
				})
			case "retrieval":
				tools = append(tools, openaiapi.FileSearchToolParam{
					Type: openaiapi.F(openaiapi.FileSearchToolTypeFileSearch),
				})
			case "function":
				var params shared.FunctionParameters
				if err := json.Unmarshal([]byte(toolMap["parameters"].(string)), &params); err != nil {
					return diag.FromErr(err)
				}
				tools = append(tools, openaiapi.FunctionToolParam{
					Type: openaiapi.F(openaiapi.FunctionToolTypeFunction),
					Function: openaiapi.F(openaiapi.FunctionDefinitionParam{
						Name:        openaiapi.F(toolMap["name"].(string)),
						Description: openaiapi.F(toolMap["description"].(string)),
						Parameters:  openaiapi.F(params),
					}),
				})
			}
		}
	}

	fileIDs := []string{}
	if v, ok := d.GetOk("file_ids"); ok {
		fileIDsRaw := v.([]interface{})
		for _, id := range fileIDsRaw {
			fileIDs = append(fileIDs, id.(string))
		}
	}

	metadata := shared.MetadataParam{}
	if v, ok := d.GetOk("metadata"); ok {
		metaMap := v.(map[string]interface{})
		for key, value := range metaMap {
			metadata[key] = value.(string)
		}
	}

	toolResources := openaiapi.BetaAssistantUpdateParamsToolResources{}
	if len(fileIDs) > 0 {
		toolResources.CodeInterpreter = openaiapi.F(openaiapi.BetaAssistantUpdateParamsToolResourcesCodeInterpreter{
			FileIDs: openaiapi.F(fileIDs),
		})
	}

	_, err := client.Beta.Assistants.Update(ctx, d.Id(), openaiapi.BetaAssistantUpdateParams{
		Name:          openaiapi.F(d.Get("name").(string)),
		Model:         openaiapi.F(openaiapi.BetaAssistantUpdateParamsModel(d.Get("model").(string))),
		Description:   openaiapi.F(d.Get("description").(string)),
		Instructions:  openaiapi.F(d.Get("instructions").(string)),
		Tools:         openaiapi.F(tools),
		Metadata:      openaiapi.F(metadata),
		ToolResources: openaiapi.F(toolResources),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOpenAIAssistantRead(ctx, d, m)
}

func resourceOpenAIAssistantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	_, err := client.Beta.Assistants.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
