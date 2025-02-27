package openai

import (
	"context"
	"strings"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenAIModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenAIModelsRead,
		Schema: map[string]*schema.Schema{
			"filter_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter models by ID prefix",
			},
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
	client := m.(testutil.ClientInterface)

	models, err := client.ListModels(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	prefix := d.Get("filter_prefix").(string)
	var filteredModels []testutil.Model
	if prefix != "" {
		for _, model := range models {
			if strings.HasPrefix(model.ID, prefix) {
				filteredModels = append(filteredModels, model)
			}
		}
		models = filteredModels
	}

	modelsList := make([]interface{}, len(models))
	for i, model := range models {
		m := make(map[string]interface{})
		m["id"] = model.ID
		m["owned_by"] = model.OwnedBy

		permissions := make([]interface{}, len(model.Permission))
		for j, p := range model.Permission {
			permission := make(map[string]interface{})
			permission["id"] = p.ID
			permission["object"] = p.Object
			permission["created"] = p.Created
			permission["allow_create_engine"] = p.AllowCreateEngine
			permission["allow_sampling"] = p.AllowSampling
			permission["allow_fine_tuning"] = p.AllowFineTuning
			permissions[j] = permission
		}
		m["permission"] = permissions
		modelsList[i] = m
	}

	if err := d.Set("models", modelsList); err != nil {
		return diag.FromErr(err)
	}

	// Generate a unique ID for this list
	d.SetId("openai-models")

	return nil
}
