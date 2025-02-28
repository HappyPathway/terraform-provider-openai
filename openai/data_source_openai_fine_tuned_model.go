package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenAIFineTunedModel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenAIFineTunedModelRead,
		Schema: map[string]*schema.Schema{
			"model_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the fine-tuned model to retrieve",
			},
			"created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp (in seconds) for when the model was created",
			},
			"owned_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization that owns the model",
			},
			"object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The object type, which is always 'model'",
			},
		},
	}
}

func dataSourceOpenAIFineTunedModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	modelID := d.Get("model_id").(string)

	model, err := client.Models.Get(ctx, modelID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(model.ID)
	d.Set("created", model.Created)
	d.Set("owned_by", model.OwnedBy)
	d.Set("object", model.Object)

	return nil
}
