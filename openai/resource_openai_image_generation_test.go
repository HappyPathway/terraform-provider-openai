package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceImageGeneration_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceImageGenerationConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_image_generation.test", "prompt", "A test image of a smiling robot"),
					resource.TestCheckResourceAttr("openai_image_generation.test", "model", "dall-e-2"),
					resource.TestCheckResourceAttr("openai_image_generation.test", "n", "1"),
					resource.TestCheckResourceAttr("openai_image_generation.test", "size", "1024x1024"),
					resource.TestCheckResourceAttrSet("openai_image_generation.test", "images.0.url"),
					resource.TestCheckResourceAttr("openai_image_generation.test", "images.0.b64_json", ""),
				),
			},
		},
	})
}

func TestAccResourceImageGeneration_dalle3(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceImageGenerationConfig_dalle3(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_image_generation.test_dalle3", "prompt", "A serene mountain landscape with a misty lake at sunrise"),
					resource.TestCheckResourceAttr("openai_image_generation.test_dalle3", "model", "dall-e-3"),
					resource.TestCheckResourceAttr("openai_image_generation.test_dalle3", "quality", "hd"),
					resource.TestCheckResourceAttr("openai_image_generation.test_dalle3", "style", "natural"),
					resource.TestCheckResourceAttr("openai_image_generation.test_dalle3", "size", "1024x1024"),
					resource.TestCheckResourceAttrSet("openai_image_generation.test_dalle3", "images.0.url"),
					resource.TestCheckResourceAttrSet("openai_image_generation.test_dalle3", "images.0.revised_prompt"),
				),
			},
		},
	})
}

func TestAccResourceImageGeneration_advanced(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceImageGenerationConfig_advanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_image_generation.test_advanced", "prompt", "A detailed digital artwork of a futuristic city"),
					resource.TestCheckResourceAttr("openai_image_generation.test_advanced", "model", "dall-e-2"),
					resource.TestCheckResourceAttr("openai_image_generation.test_advanced", "n", "2"),
					resource.TestCheckResourceAttr("openai_image_generation.test_advanced", "size", "1024x1024"),
					resource.TestCheckResourceAttr("openai_image_generation.test_advanced", "response_format", "b64_json"),
					resource.TestCheckResourceAttrSet("openai_image_generation.test_advanced", "images.0.b64_json"),
					resource.TestCheckResourceAttrSet("openai_image_generation.test_advanced", "images.1.b64_json"),
				),
			},
		},
	})
}

func testAccResourceImageGenerationConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_image_generation" "test" {
  prompt = "A test image of a smiling robot"
  model  = "dall-e-2"
  n      = 1
  size   = "1024x1024"
}
`)
}

func testAccResourceImageGenerationConfig_dalle3() string {
	return fmt.Sprintf(`
resource "openai_image_generation" "test_dalle3" {
  prompt  = "A serene mountain landscape with a misty lake at sunrise"
  model   = "dall-e-3"
  quality = "hd"
  style   = "natural"
  size    = "1024x1024"
}
`)
}

func testAccResourceImageGenerationConfig_advanced() string {
	return fmt.Sprintf(`
resource "openai_image_generation" "test_advanced" {
  prompt          = "A detailed digital artwork of a futuristic city"
  model          = "dall-e-2"
  n              = 2
  size           = "1024x1024"
  response_format = "b64_json"
}
`)
}
