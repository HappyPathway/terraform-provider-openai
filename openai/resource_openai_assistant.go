package openai

import (
	"context"
	"encoding/json"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
							Description: "The type of tool",
						},
						"function": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the function",
									},
									"description": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The description of what the function does",
									},
									"parameters": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameters the functions accepts (JSON encoded string)",
									},
								},
							},
							Description: "The function definition when type is 'function'",
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

func parseParameters(parametersStr string) (map[string]interface{}, error) {
	var parameters map[string]interface{}
	if err := json.Unmarshal([]byte(parametersStr), &parameters); err != nil {
		return nil, err
	}
	return parameters, nil
}

func resourceOpenAIAssistantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	req := &testutil.CreateAssistantRequest{
		Name:         d.Get("name").(string),
		Model:        d.Get("model").(string),
		Description:  d.Get("description").(string),
		Instructions: d.Get("instructions").(string),
	}

	if v, ok := d.GetOk("tools"); ok {
		toolsList := v.([]interface{})
		tools := make([]testutil.AssistantTool, len(toolsList))
		for i, tool := range toolsList {
			toolMap := tool.(map[string]interface{})
			tools[i] = testutil.AssistantTool{
				Type: toolMap["type"].(string),
			}

			// Handle function tool type
			if tools[i].Type == "function" && toolMap["function"] != nil {
				functionList := toolMap["function"].([]interface{})
				if len(functionList) > 0 {
					function := functionList[0].(map[string]interface{})
					parameters, err := parseParameters(function["parameters"].(string))
					if err != nil {
						return diag.FromErr(err)
					}
					tools[i].Function = &testutil.AssistantFunction{
						Name:        function["name"].(string),
						Description: function["description"].(string),
						Parameters:  parameters,
					}
				}
			}
		}
		req.Tools = tools
	}

	if v, ok := d.GetOk("file_ids"); ok {
		fileIDs := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			fileIDs[i] = id.(string)
		}
		req.FileIDs = fileIDs
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			metadata[key] = value.(string)
		}
		req.Metadata = metadata
	}

	assistant, err := client.CreateAssistant(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(assistant.ID)

	return resourceOpenAIAssistantRead(ctx, d, m)
}

func resourceOpenAIAssistantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	assistant, err := client.GetAssistant(ctx, d.Id())
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
			"type": tool.Type,
		}

		if tool.Function != nil {
			toolMap["function"] = []interface{}{
				map[string]interface{}{
					"name":        tool.Function.Name,
					"description": tool.Function.Description,
					"parameters":  tool.Function.Parameters,
				},
			}
		}
		tools[i] = toolMap
	}
	d.Set("tools", tools)

	d.Set("file_ids", assistant.FileIDs)
	d.Set("metadata", assistant.Metadata)

	return nil
}

func resourceOpenAIAssistantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	req := &testutil.CreateAssistantRequest{
		Name:         d.Get("name").(string),
		Model:        d.Get("model").(string),
		Description:  d.Get("description").(string),
		Instructions: d.Get("instructions").(string),
	}

	if v, ok := d.GetOk("tools"); ok {
		toolsList := v.([]interface{})
		tools := make([]testutil.AssistantTool, len(toolsList))
		for i, tool := range toolsList {
			toolMap := tool.(map[string]interface{})
			tools[i] = testutil.AssistantTool{
				Type: toolMap["type"].(string),
			}

			// Handle function tool type
			if tools[i].Type == "function" && toolMap["function"] != nil {
				functionList := toolMap["function"].([]interface{})
				if len(functionList) > 0 {
					function := functionList[0].(map[string]interface{})
					parameters, err := parseParameters(function["parameters"].(string))
					if err != nil {
						return diag.FromErr(err)
					}
					tools[i].Function = &testutil.AssistantFunction{
						Name:        function["name"].(string),
						Description: function["description"].(string),
						Parameters:  parameters,
					}
				}
			}
		}
		req.Tools = tools
	}

	if v, ok := d.GetOk("file_ids"); ok {
		fileIDs := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			fileIDs[i] = id.(string)
		}
		req.FileIDs = fileIDs
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			metadata[key] = value.(string)
		}
		req.Metadata = metadata
	}

	_, err := client.UpdateAssistant(ctx, d.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOpenAIAssistantRead(ctx, d, m)
}

func resourceOpenAIAssistantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	err := client.DeleteAssistant(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
