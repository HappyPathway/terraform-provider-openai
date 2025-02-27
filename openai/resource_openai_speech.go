package openai

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpenAISpeech() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAISpeechCreate,
		ReadContext:   resourceOpenAISpeechRead,
		DeleteContext: resourceOpenAISpeechDelete,
		Schema: map[string]*schema.Schema{
			"input": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The text to generate audio for. The maximum length is 4096 characters.",
			},
			"model": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "tts-1",
				ValidateFunc: validation.StringInSlice([]string{"tts-1", "tts-1-hd"}, false),
				Description:  "The TTS model to use. One of 'tts-1' or 'tts-1-hd'.",
			},
			"voice": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"alloy", "ash", "coral", "echo", "fable",
					"onyx", "nova", "sage", "shimmer",
				}, false),
				Description: "The voice to use. Available voices: alloy, ash, coral, echo, fable, onyx, nova, sage, shimmer.",
			},
			"response_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "mp3",
				ValidateFunc: validation.StringInSlice([]string{"mp3", "opus", "aac", "flac", "wav", "pcm"}, false),
				Description:  "The format to generate audio in. One of mp3, opus, aac, flac, wav, or pcm.",
			},
			"speed": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ForceNew:     true,
				Default:      1.0,
				ValidateFunc: validation.FloatBetween(0.25, 4.0),
				Description:  "The speed of the generated audio. Select a value from 0.25 to 4.0. Default is 1.0.",
			},
			"audio_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated audio content in base64 encoded format.",
			},
		},
	}
}

func resourceOpenAISpeechCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	input := d.Get("input").(string)
	model := d.Get("model").(string)
	voice := d.Get("voice").(string)
	responseFormat := d.Get("response_format").(string)
	speed := d.Get("speed").(float64)

	resp, err := client.CreateSpeech(ctx, &CreateSpeechRequest{
		Input:          input,
		Model:          model,
		Voice:          voice,
		ResponseFormat: responseFormat,
		Speed:          speed,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating speech: %v", err))
	}

	// Set ID to a combination of input and model
	d.SetId(fmt.Sprintf("%s-%s", input, model))

	// Set the audio content
	d.Set("audio_content", resp)

	return nil
}

func resourceOpenAISpeechRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Speech is stateless and can't be retrieved after creation
	return nil
}

func resourceOpenAISpeechDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Speech results don't need to be deleted as they are stateless
	return nil
}
