package openai

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpenAITranscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAITranscriptionCreate,
		ReadContext:   resourceOpenAITranscriptionRead,
		DeleteContext: resourceOpenAITranscriptionDelete,
		Schema: map[string]*schema.Schema{
			"audio_content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The base64-encoded audio file content to transcribe. Supported formats: flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.",
			},
			"model": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "whisper-1",
				ValidateFunc: validation.StringInSlice([]string{"whisper-1"}, false),
				Description:  "The model to use for transcription. Currently only whisper-1 is supported.",
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The language of the input audio in ISO-639-1 format. If not specified, the model will auto-detect the language.",
			},
			"prompt": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional text to guide the model's style or continue a previous audio segment.",
			},
			"response_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "text",
				ValidateFunc: validation.StringInSlice([]string{"json", "text", "srt", "verbose_json", "vtt"}, false),
				Description:  "The format of the transcript output. Must be one of: json, text, srt, verbose_json, or vtt.",
			},
			"temperature": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ForceNew:     true,
				Default:      0,
				ValidateFunc: validation.FloatBetween(0, 1),
				Description:  "The sampling temperature, between 0 and 1. Higher values make the output more random, lower values more deterministic.",
			},
			"text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The transcribed text.",
			},
		},
	}
}

func resourceOpenAITranscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	fileContent := d.Get("file_content").(string)
	fileBytes, err := base64.StdEncoding.DecodeString(fileContent)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error decoding file content: %v", err))
	}

	req := &testutil.TranscriptionRequest{
		File:           fileBytes,
		Model:          d.Get("model").(string),
		Language:       d.Get("language").(string),
		Prompt:         d.Get("prompt").(string),
		ResponseFormat: d.Get("response_format").(string),
		Temperature:    float32(d.Get("temperature").(float64)),
	}

	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating transcription: %v", err))
	}

	// Set ID to a hash of the audio content
	d.SetId(fmt.Sprintf("%x", fileBytes[:16]))

	// Set the transcribed text
	d.Set("text", resp.Text)

	return nil
}

func resourceOpenAITranscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Transcription is stateless and can't be retrieved after creation
	return nil
}

func resourceOpenAITranscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Transcription results don't need to be deleted as they are stateless
	return nil
}
