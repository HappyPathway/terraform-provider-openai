package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTranslation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTranslationConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_translation.test", "audio_content"),
					resource.TestCheckResourceAttr("openai_translation.test", "model", "whisper-1"),
					resource.TestCheckResourceAttrSet("openai_translation.test", "text"),
				),
			},
		},
	})
}

func TestAccResourceTranslation_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTranslationConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_translation.test_full", "audio_content"),
					resource.TestCheckResourceAttr("openai_translation.test_full", "model", "whisper-1"),
					resource.TestCheckResourceAttr("openai_translation.test_full", "prompt", "This is a business meeting discussing sales figures."),
					resource.TestCheckResourceAttr("openai_translation.test_full", "response_format", "verbose_json"),
					resource.TestCheckResourceAttr("openai_translation.test_full", "temperature", "0.3"),
					resource.TestCheckResourceAttrSet("openai_translation.test_full", "text"),
				),
			},
		},
	})
}

func testAccResourceTranslationConfig_basic() string {
	// In a real test, you would provide actual base64-encoded audio content
	return fmt.Sprintf(`
resource "openai_translation" "test" {
  audio_content = "U29tZSBiYXNlNjQgYXVkaW8gY29udGVudA=="  # "Some base64 audio content" in base64
}
`)
}

func testAccResourceTranslationConfig_full() string {
	// In a real test, you would provide actual base64-encoded audio content
	return fmt.Sprintf(`
resource "openai_translation" "test_full" {
  audio_content    = "U29tZSBiYXNlNjQgYXVkaW8gY29udGVudA=="  # "Some base64 audio content" in base64
  model           = "whisper-1"
  prompt          = "This is a business meeting discussing sales figures."
  response_format = "verbose_json"
  temperature     = 0.3
}
`)
}
