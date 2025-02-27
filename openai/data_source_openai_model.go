package openai

import (
	"context"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenAIModel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenAIModelRead,
		Schema: map[string]*schema.Schema{
			"model_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the OpenAI model to retrieve",
			},
			"owned_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization that owns the model",
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
	}
}

func dataSourceOpenAIModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	modelID := d.Get("model_id").(string)
	model, err := client.GetModel(ctx, modelID)
	if err != nil {
		return diag.FromErr(err)
	}
	if model == nil {
		return diag.Errorf("Model %s not found", modelID)
	}

	d.SetId(model.ID)
	if err := d.Set("owned_by", model.OwnedBy); err != nil {
		return diag.FromErr(err)
	}

	permissions := make([]interface{}, len(model.Permission))
	for i, p := range model.Permission {
		permission := make(map[string]interface{})
		permission["id"] = p.ID
		permission["object"] = p.Object
		permission["created"] = p.Created
		permission["allow_create_engine"] = p.AllowCreateEngine
		permission["allow_sampling"] = p.AllowSampling
		permission["allow_fine_tuning"] = p.AllowFineTuning
		permissions[i] = permission
	}
	if err := d.Set("permission", permissions); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
