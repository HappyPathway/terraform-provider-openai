package openai

import (
	"context"
	"fmt"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenAIModeration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIModerationCreate,
		ReadContext:   resourceOpenAIModerationRead,
		DeleteContext: resourceOpenAIModerationDelete,
		Schema: map[string]*schema.Schema{
			"input": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The text to moderate.",
			},
			"model": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The moderation model to use.",
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flagged": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"categories": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
						"category_scores": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
						},
					},
				},
			},
		},
	}
}

func resourceOpenAIModerationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	req := &testutil.CreateModerationRequest{
		Input: d.Get("input").(string),
		Model: d.Get("model").(string),
	}

	moderation, err := client.CreateModeration(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating moderation: %v", err))
	}

	// Set ID to input string
	d.SetId(d.Get("input").(string))

	results := make([]interface{}, len(moderation.Results))
	for i, result := range moderation.Results {
		categories := map[string]interface{}{
			"hate":                   result.Categories.Hate,
			"hate/threatening":       result.Categories.HateThreatening,
			"harassment":             result.Categories.Harassment,
			"harassment/threatening": result.Categories.HarassmentThreatening,
			"self-harm":              result.Categories.SelfHarm,
			"self-harm/instructions": result.Categories.SelfHarmInstructions,
			"sexual":                 result.Categories.Sexual,
			"sexual/minors":          result.Categories.SexualMinors,
			"violence":               result.Categories.Violence,
			"violence/graphic":       result.Categories.ViolenceGraphic,
			"illicit":                result.Categories.Illicit,
			"illicit/violent":        result.Categories.IllicitViolent,
		}

		categoryScores := map[string]interface{}{
			"hate":             result.CategoryScores.Hate,
			"hate/threatening": result.CategoryScores.HateThreatening,
			"self-harm":        result.CategoryScores.SelfHarm,
			"sexual":           result.CategoryScores.Sexual,
			"sexual/minors":    result.CategoryScores.SexualMinors,
			"violence":         result.CategoryScores.Violence,
			"violence/graphic": result.CategoryScores.ViolenceGraphic,
		}

		resultMap := map[string]interface{}{
			"flagged":         result.Flagged,
			"categories":      categories,
			"category_scores": categoryScores,
		}

		results[i] = resultMap
	}

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOpenAIModerationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Moderations are stateless and can't be retrieved after creation
	return nil
}

func resourceOpenAIModerationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Moderations are stateless and don't need explicit deletion
	return nil
}
