package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSpeech_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpeechConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_speech.test", "input", "Hello, this is a test of the text to speech system."),
					resource.TestCheckResourceAttr("openai_speech.test", "voice", "alloy"),
					resource.TestCheckResourceAttr("openai_speech.test", "model", "tts-1"),
					resource.TestCheckResourceAttrSet("openai_speech.test", "audio_content"),
				),
			},
		},
	})
}

func TestAccResourceSpeech_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpeechConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_speech.test_full", "input", "Welcome to our application! We hope you enjoy your stay."),
					resource.TestCheckResourceAttr("openai_speech.test_full", "model", "tts-1-hd"),
					resource.TestCheckResourceAttr("openai_speech.test_full", "voice", "nova"),
					resource.TestCheckResourceAttr("openai_speech.test_full", "response_format", "mp3"),
					resource.TestCheckResourceAttr("openai_speech.test_full", "speed", "1.2"),
					resource.TestCheckResourceAttrSet("openai_speech.test_full", "audio_content"),
				),
			},
		},
	})
}

func testAccResourceSpeechConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_speech" "test" {
  input = "Hello, this is a test of the text to speech system."
  voice = "alloy"
}
`)
}

func testAccResourceSpeechConfig_full() string {
	return fmt.Sprintf(`
resource "openai_speech" "test_full" {
  input           = "Welcome to our application! We hope you enjoy your stay."
  model           = "tts-1-hd"
  voice           = "nova"
  response_format = "mp3"
  speed           = 1.2
}
`)
}
