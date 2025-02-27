package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceContentGenerator_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceContentGeneratorConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_content_generator.test", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_content_generator.test", "messages.#", "1"),
					resource.TestCheckResourceAttr("openai_content_generator.test", "messages.0.role", "user"),
					resource.TestCheckResourceAttrSet("openai_content_generator.test", "content"),
					resource.TestCheckResourceAttrSet("openai_content_generator.test", "raw_response"),
					resource.TestCheckResourceAttrSet("openai_content_generator.test", "usage.completion_tokens"),
					resource.TestCheckResourceAttrSet("openai_content_generator.test", "usage.prompt_tokens"),
					resource.TestCheckResourceAttrSet("openai_content_generator.test", "usage.total_tokens"),
				),
			},
		},
	})
}

func TestAccResourceContentGenerator_withSchema(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceContentGeneratorConfig_withSchema(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_content_generator.test_schema", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_content_generator.test_schema", "messages.#", "2"),
					resource.TestCheckResourceAttr("openai_content_generator.test_schema", "response_format.0.type", "json_object"),
					resource.TestCheckResourceAttrSet("openai_content_generator.test_schema", "content"),
				),
			},
		},
	})
}

func testAccResourceContentGeneratorConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_content_generator" "test" {
  model = "gpt-4"
  messages {
    role    = "user"
    content = "Write a haiku about coding"
  }
  temperature = 0.7
}
`)
}

func testAccResourceContentGeneratorConfig_withSchema() string {
	return fmt.Sprintf(`
resource "openai_content_generator" "test_schema" {
  model = "gpt-4"
  
  messages {
    role    = "system"
    content = "You are a helpful assistant that provides structured data about books."
  }
  
  messages {
    role    = "user"
    content = "Provide information about The Hobbit by J.R.R. Tolkien"
  }

  response_format {
    type = "json_object"
    schema = jsonencode({
      type = "object"
      properties = {
        title = {
          type = "string"
          description = "The book title"
        }
        author = {
          type = "string"
          description = "The book author"
        }
        year_published = {
          type = "integer"
          description = "Year the book was first published"
        }
        genre = {
          type = "array"
          items = {
            type = "string"
          }
          description = "List of genres"
        }
      }
      required = ["title", "author", "year_published", "genre"]
    })
  }

  temperature = 0.5
}
`)
}
