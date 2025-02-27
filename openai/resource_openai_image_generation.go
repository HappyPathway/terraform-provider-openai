package openai

import (
	"context"
	"fmt"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpenAIImageGeneration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIImageGenerationCreate,
		ReadContext:   resourceOpenAIImageGenerationRead,
		DeleteContext: resourceOpenAIImageGenerationDelete,
		Schema: map[string]*schema.Schema{
			"prompt": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A text description of the desired image(s). Maximum length is 1000 characters for dall-e-2 and 4000 characters for dall-e-3.",
			},
			"model": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "dall-e-2",
				ValidateFunc: validation.StringInSlice([]string{"dall-e-2", "dall-e-3"}, false),
				Description:  "The model to use for image generation. Defaults to dall-e-2.",
			},
			"n": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "The number of images to generate. Must be between 1 and 10. For dall-e-3, only n=1 is supported.",
			},
			"quality": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "hd"}, false),
				Description:  "The quality of the image that will be generated. 'hd' creates images with finer details. Only supported for dall-e-3.",
			},
			"response_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "url",
				ValidateFunc: validation.StringInSlice([]string{"url", "b64_json"}, false),
				Description:  "The format in which the generated images are returned. Must be one of 'url' or 'b64_json'.",
			},
			"size": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "1024x1024",
				ValidateFunc: validation.StringInSlice([]string{"256x256", "512x512", "1024x1024", "1792x1024", "1024x1792"}, false),
				Description:  "The size of the generated images.",
			},
			"style": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "vivid",
				ValidateFunc: validation.StringInSlice([]string{"vivid", "natural"}, false),
				Description:  "The style of the generated images. Only supported for dall-e-3.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A unique identifier representing your end-user.",
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"b64_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The base64-encoded JSON of the generated image.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the generated image.",
						},
						"revised_prompt": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The prompt that was used to generate the image, if there was any revision.",
						},
					},
				},
			},
		},
	}
}

func resourceOpenAIImageGenerationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	req := &testutil.CreateImageRequest{
		Prompt:         d.Get("prompt").(string),
		Model:          d.Get("model").(string),
		N:              d.Get("n").(int),
		Quality:        d.Get("quality").(string),
		ResponseFormat: d.Get("response_format").(string),
		Size:           d.Get("size").(string),
		Style:          d.Get("style").(string),
		User:           d.Get("user").(string),
	}

	response, err := client.CreateImage(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating image: %v", err))
	}

	// Generate a unique ID based on prompt and timestamp
	d.SetId(fmt.Sprintf("%s-%d", d.Get("prompt").(string), response.Created))

	images := make([]interface{}, len(response.Data))
	for i, img := range response.Data {
		imageMap := map[string]interface{}{
			"b64_json":       img.B64JSON,
			"url":            img.URL,
			"revised_prompt": img.RevisedPrompt,
		}
		images[i] = imageMap
	}

	if err := d.Set("images", images); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOpenAIImageGenerationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Images are stateless and can't be retrieved after creation
	return nil
}

func resourceOpenAIImageGenerationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Images don't need to be deleted as they are stateless
	return nil
}
