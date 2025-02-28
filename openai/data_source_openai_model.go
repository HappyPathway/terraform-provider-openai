package openai

import (
	"context"

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
				Description: "The ID of the model to retrieve",
			},
			"owned_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization that owns the model",
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
	config := m.(*Config)
	client := config.Client

	modelID := d.Get("model_id").(string)
	model, err := client.Models.Get(ctx, modelID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(model.ID)
	if err := d.Set("owned_by", model.OwnedBy); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
