package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

func resourceOpenAIVectorStore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIVectorStoreCreate,
		ReadContext:   resourceOpenAIVectorStoreRead,
		UpdateContext: resourceOpenAIVectorStoreUpdate,
		DeleteContext: resourceOpenAIVectorStoreDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the vector store.",
			},
			"file_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of File IDs to add to the vector store. Maximum of 10000 files.",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of key-value pairs that can be attached to the vector store.",
			},
			"expires_after": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"anchor": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Anchor timestamp after which the expiration policy applies. Currently only supports 'last_active_at'.",
						},
						"days": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of days after the anchor time that the vector store will expire.",
						},
					},
				},
				Description: "The expiration policy for the vector store.",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp (in seconds) for when the vector store was created.",
			},
			"last_active_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp (in seconds) for when the vector store was last active.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the vector store (expired, in_progress, or completed).",
			},
			"usage_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of bytes used by the files in the vector store.",
			},
			"expires_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp (in seconds) for when the vector store will expire.",
			},
		},
	}
}

func resourceOpenAIVectorStoreCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	params := openaiapi.BetaVectorStoreNewParams{
		Name: openaiapi.F(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("file_ids"); ok {
		fileIDsList := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			fileIDsList[i] = id.(string)
		}
		params.FileIDs = openaiapi.F(fileIDsList)
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := shared.MetadataParam{}
		for key, value := range v.(map[string]interface{}) {
			metadata[key] = value.(string)
		}
		params.Metadata = openaiapi.F(metadata)
	}

	if v, ok := d.GetOk("expires_after"); ok {
		expiresAfter := v.([]interface{})[0].(map[string]interface{})
		params.ExpiresAfter = openaiapi.F(openaiapi.BetaVectorStoreNewParamsExpiresAfter{
			Anchor: openaiapi.F(openaiapi.BetaVectorStoreNewParamsExpiresAfterAnchor(expiresAfter["anchor"].(string))),
			Days:   openaiapi.F(int64(expiresAfter["days"].(int))),
		})
	}

	vectorStore, err := client.Beta.VectorStores.New(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vectorStore.ID)
	return resourceOpenAIVectorStoreRead(ctx, d, m)
}

func resourceOpenAIVectorStoreRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	vectorStore, err := client.Beta.VectorStores.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", vectorStore.Name)
	d.Set("created_at", vectorStore.CreatedAt)
	d.Set("last_active_at", vectorStore.LastActiveAt)
	d.Set("status", string(vectorStore.Status))
	d.Set("usage_bytes", vectorStore.UsageBytes)
	d.Set("expires_at", vectorStore.ExpiresAt)

	if vectorStore.Metadata != nil {
		d.Set("metadata", vectorStore.Metadata)
	}

	return nil
}

func resourceOpenAIVectorStoreUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	params := openaiapi.BetaVectorStoreUpdateParams{}

	if d.HasChange("name") {
		params.Name = openaiapi.F(d.Get("name").(string))
	}

	if d.HasChange("metadata") {
		metadata := shared.MetadataParam{}
		if v, ok := d.GetOk("metadata"); ok {
			for key, value := range v.(map[string]interface{}) {
				metadata[key] = value.(string)
			}
		}
		params.Metadata = openaiapi.F(metadata)
	}

	if d.HasChange("expires_after") {
		if v, ok := d.GetOk("expires_after"); ok {
			expiresAfter := v.([]interface{})[0].(map[string]interface{})
			params.ExpiresAfter = openaiapi.F(openaiapi.BetaVectorStoreUpdateParamsExpiresAfter{
				Anchor: openaiapi.F(openaiapi.BetaVectorStoreUpdateParamsExpiresAfterAnchor(expiresAfter["anchor"].(string))),
				Days:   openaiapi.F(int64(expiresAfter["days"].(int))),
			})
		}
	}

	_, err := client.Beta.VectorStores.Update(ctx, d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOpenAIVectorStoreRead(ctx, d, m)
}

func resourceOpenAIVectorStoreDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	_, err := client.Beta.VectorStores.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
