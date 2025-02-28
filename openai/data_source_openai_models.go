package openai

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openaiapi "github.com/openai/openai-go"
)

func dataSourceOpenAIModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenAIModelsRead,
		Schema: map[string]*schema.Schema{
			"models": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owned_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"object": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"created": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"allow_create_engine": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"allow_sampling": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"allow_fine_tuning": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceOpenAIModelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	modelPage, err := client.Models.List(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("models", flattenModels(modelPage.Data)); err != nil {
		return diag.FromErr(err)
	}

	// Generate a consistent ID for this data source
	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenModels(models []openaiapi.Model) []interface{} {
	var result []interface{}

	for _, model := range models {
		m := make(map[string]interface{})
		m["id"] = model.ID
		m["owned_by"] = model.OwnedBy

		perms := make([]interface{}, 1)
		p := make(map[string]interface{})
		p["id"] = model.ID // Use model.ID instead of JSON field
		p["object"] = string(model.Object)
		p["created"] = model.Created

		// These fields might be in ExtraFields as JSON strings
		if v, ok := model.JSON.ExtraFields["allow_create_engine"]; ok {
			var b bool
			if err := json.Unmarshal([]byte(v.Raw()), &b); err == nil {
				p["allow_create_engine"] = b
			}
		}
		if v, ok := model.JSON.ExtraFields["allow_sampling"]; ok {
			var b bool
			if err := json.Unmarshal([]byte(v.Raw()), &b); err == nil {
				p["allow_sampling"] = b
			}
		}
		if v, ok := model.JSON.ExtraFields["allow_fine_tuning"]; ok {
			var b bool
			if err := json.Unmarshal([]byte(v.Raw()), &b); err == nil {
				p["allow_fine_tuning"] = b
			}
		}

		perms[0] = p
		m["permission"] = perms

		result = append(result, m)
	}

	return result
}
