package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTranscription_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTranscriptionConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_transcription.test", "audio_content"),
					resource.TestCheckResourceAttr("openai_transcription.test", "model", "whisper-1"),
					resource.TestCheckResourceAttrSet("openai_transcription.test", "text"),
				),
			},
		},
	})
}

func TestAccResourceTranscription_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTranscriptionConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_transcription.test_full", "audio_content"),
					resource.TestCheckResourceAttr("openai_transcription.test_full", "model", "whisper-1"),
					resource.TestCheckResourceAttr("openai_transcription.test_full", "language", "en"),
					resource.TestCheckResourceAttr("openai_transcription.test_full", "prompt", "This is a business meeting discussion."),
					resource.TestCheckResourceAttr("openai_transcription.test_full", "response_format", "verbose_json"),
					resource.TestCheckResourceAttr("openai_transcription.test_full", "temperature", "0.3"),
					resource.TestCheckResourceAttrSet("openai_transcription.test_full", "text"),
				),
			},
		},
	})
}

func testAccResourceTranscriptionConfig_basic() string {
	// In a real test, you would provide actual base64-encoded audio content
	return fmt.Sprintf(`
resource "openai_transcription" "test" {
  audio_content = "U29tZSBiYXNlNjQgYXVkaW8gY29udGVudA=="  # "Some base64 audio content" in base64
}
`)
}

func testAccResourceTranscriptionConfig_full() string {
	// In a real test, you would provide actual base64-encoded audio content
	return fmt.Sprintf(`
resource "openai_transcription" "test_full" {
  audio_content    = "U29tZSBiYXNlNjQgYXVkaW8gY29udGVudA=="  # "Some base64 audio content" in base64
  model           = "whisper-1"
  language        = "en"
  prompt          = "This is a business meeting discussion."
  response_format = "verbose_json"
  temperature     = 0.3
}
`)
}
