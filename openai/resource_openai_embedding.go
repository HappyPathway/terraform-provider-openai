package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenAIEmbedding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIEmbeddingCreate,
		ReadContext:   resourceOpenAIEmbeddingRead,
		DeleteContext: resourceOpenAIEmbeddingDelete,
		Schema: map[string]*schema.Schema{
			"model": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"input": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"encoding_format": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "float",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "float" && v != "base64" {
						errs = append(errs, fmt.Errorf("%q must be either 'float' or 'base64', got: %s", key, v))
					}
					return
				},
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed fields from API response
			"embeddings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"embedding": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"usage": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceOpenAIEmbeddingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	req := &CreateEmbeddingRequest{
		Model: d.Get("model").(string),
		Input: d.Get("input").(string),
	}

	if v, ok := d.GetOk("dimensions"); ok {
		req.Dimensions = v.(int)
	}

	if v, ok := d.GetOk("encoding_format"); ok {
		req.EncodingFormat = v.(string)
	}

	if v, ok := d.GetOk("user"); ok {
		req.User = v.(string)
	}

	embedding, err := client.CreateEmbedding(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating embedding: %v", err))
	}

	// Set ID to a combination of model and input
	d.SetId(fmt.Sprintf("%s-%s", d.Get("model").(string), d.Get("input").(string)))

	// Set computed values
	embeddings := make([]interface{}, len(embedding.Data))
	for i, data := range embedding.Data {
		embeddingMap := map[string]interface{}{
			"embedding": data.Embedding,
			"index":     data.Index,
		}
		embeddings[i] = embeddingMap
	}
	if err := d.Set("embeddings", embeddings); err != nil {
		return diag.FromErr(err)
	}

	usage := map[string]interface{}{
		"prompt_tokens": embedding.Usage.PromptTokens,
		"total_tokens":  embedding.Usage.TotalTokens,
	}
	if err := d.Set("usage", usage); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOpenAIEmbeddingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Embeddings are stateless and can't be retrieved after creation
	return nil
}

func resourceOpenAIEmbeddingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Embeddings are stateless and don't need explicit deletion
	return nil
}
