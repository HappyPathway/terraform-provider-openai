package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description: "The text to classify.",
			},
			"model": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "text-moderation-latest",
				ValidateFunc: validation.StringInSlice([]string{
					"text-moderation-latest",
					"text-moderation-stable",
					"omni-moderation-2024-09-26",
				}, false),
				Description: "The moderation model to use.",
			},
			"flagged": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the content is flagged as potentially harmful.",
			},
			"categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"harassment": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that expresses, incites, or promotes harassing language towards any target.",
						},
						"harassment_threatening": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Harassment content that also includes violence or serious harm towards any target.",
						},
						"hate": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that expresses, incites, or promotes hate based on protected attributes.",
						},
						"hate_threatening": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Hate content that also includes violence or serious harm towards the targeted group.",
						},
						"illicit": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that promotes or facilitates illicit activities.",
						},
						"illicit_violent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that promotes illicit activities involving violence.",
						},
						"self_harm": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that promotes, encourages, or depicts acts of self-harm.",
						},
						"self_harm_instructions": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that encourages or provides instructions for performing acts of self-harm.",
						},
						"self_harm_intent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that expresses intention to engage in acts of self-harm.",
						},
						"sexual": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content meant to arouse sexual excitement.",
						},
						"sexual_minors": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Sexual content involving minors.",
						},
						"violence": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that depicts death, violence, or physical injury.",
						},
						"violence_graphic": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Content that depicts death, violence, or physical injury in graphic detail.",
						},
					},
				},
			},
			"category_scores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"harassment": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for harassment category.",
						},
						"harassment_threatening": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for threatening harassment category.",
						},
						"hate": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for hate category.",
						},
						"hate_threatening": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for threatening hate category.",
						},
						"illicit": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for illicit activities category.",
						},
						"illicit_violent": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for violent illicit activities category.",
						},
						"self_harm": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for self-harm category.",
						},
						"self_harm_instructions": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for self-harm instructions category.",
						},
						"self_harm_intent": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for self-harm intent category.",
						},
						"sexual": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for sexual content category.",
						},
						"sexual_minors": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for sexual content involving minors category.",
						},
						"violence": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for violence category.",
						},
						"violence_graphic": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Score for graphic violence category.",
						},
					},
				},
			},
		},
	}
}

func resourceOpenAIModerationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	input := d.Get("input").(string)
	model := d.Get("model").(string)

	req := &CreateModerationRequest{
		Input: input,
		Model: model,
	}

	moderation, err := client.CreateModeration(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating moderation: %v", err))
	}

	// Set ID as a hash of the input and model
	d.SetId(fmt.Sprintf("%s-%s", input, model))

	if len(moderation.Results) > 0 {
		result := moderation.Results[0]

		d.Set("flagged", result.Flagged)

		// Set categories
		categories := []map[string]interface{}{
			{
				"harassment":             result.Categories.Harassment,
				"harassment_threatening": result.Categories.HarassmentThreatening,
				"hate":                   result.Categories.Hate,
				"hate_threatening":       result.Categories.HateThreatening,
				"illicit":                result.Categories.Illicit,
				"illicit_violent":        result.Categories.IllicitViolent,
				"self_harm":              result.Categories.SelfHarm,
				"self_harm_instructions": result.Categories.SelfHarmInstructions,
				"self_harm_intent":       result.Categories.SelfHarmIntent,
				"sexual":                 result.Categories.Sexual,
				"sexual_minors":          result.Categories.SexualMinors,
				"violence":               result.Categories.Violence,
				"violence_graphic":       result.Categories.ViolenceGraphic,
			},
		}
		if err := d.Set("categories", categories); err != nil {
			return diag.FromErr(err)
		}

		// Set category scores
		categoryScores := []map[string]interface{}{
			{
				"harassment":             result.CategoryScores.Harassment,
				"harassment_threatening": result.CategoryScores.HarassmentThreatening,
				"hate":                   result.CategoryScores.Hate,
				"hate_threatening":       result.CategoryScores.HateThreatening,
				"illicit":                result.CategoryScores.Illicit,
				"illicit_violent":        result.CategoryScores.IllicitViolent,
				"self_harm":              result.CategoryScores.SelfHarm,
				"self_harm_instructions": result.CategoryScores.SelfHarmInstructions,
				"self_harm_intent":       result.CategoryScores.SelfHarmIntent,
				"sexual":                 result.CategoryScores.Sexual,
				"sexual_minors":          result.CategoryScores.SexualMinors,
				"violence":               result.CategoryScores.Violence,
				"violence_graphic":       result.CategoryScores.ViolenceGraphic,
			},
		}
		if err := d.Set("category_scores", categoryScores); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceOpenAIModerationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Moderation is stateless and can't be retrieved after creation
	return nil
}

func resourceOpenAIModerationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Moderation results don't need to be deleted as they are stateless
	return nil
}
