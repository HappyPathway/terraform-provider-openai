package openai

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIContentGenerator_basic(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIContentGeneratorConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_content_generator.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet(
						"openai_content_generator.test", "content"),
					resource.TestCheckResourceAttrSet(
						"openai_content_generator.test", "raw_response"),
					resource.TestCheckResourceAttrSet(
						"openai_content_generator.test", "usage.prompt_tokens"),
					resource.TestCheckResourceAttrSet(
						"openai_content_generator.test", "usage.completion_tokens"),
					resource.TestCheckResourceAttrSet(
						"openai_content_generator.test", "usage.total_tokens"),
				),
			},
		},
	})
}

func TestAccResourceOpenAIContentGenerator_withJsonResponse(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIContentGeneratorConfigWithJsonResponse(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_content_generator.test_json", "model", "gpt-4"),
					resource.TestCheckResourceAttrSet(
						"openai_content_generator.test_json", "content"),
					// Verify the response is valid JSON
					resource.TestMatchResourceAttr(
						"openai_content_generator.test_json", "content",
						regexp.MustCompile(`^\{.*\}$`)),
				),
			},
		},
	})
}

func testAccResourceOpenAIContentGeneratorConfig(apiKey string) string {
	return `
provider "openai" {
  api_key = "` + apiKey + `"
}

resource "openai_content_generator" "test" {
  model = "gpt-3.5-turbo"
  temperature = 0.7

  messages {
    role    = "user"
    content = "Tell me a short joke about programming"
  }
}
`
}

func testAccResourceOpenAIContentGeneratorConfigWithJsonResponse(apiKey string) string {
	return `
provider "openai" {
  api_key = "` + apiKey + `"
}

resource "openai_content_generator" "test_json" {
  model = "gpt-4"
  temperature = 0.7

  messages {
    role    = "user"
    content = "Generate information about a random animal in JSON format"
  }

  response_format {
    type = "json_object"
    schema = jsonencode({
      type = "object"
      properties = {
        animal_name = {
          type = "string"
        }
        scientific_name = {
          type = "string"
        }
        habitat = {
          type = "string"
        }
        diet = {
          type = "string"
        }
      }
      required = ["animal_name", "scientific_name", "habitat", "diet"]
    })
  }
}
`
}
