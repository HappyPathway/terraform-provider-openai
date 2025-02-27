package openai

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIAssistant_basic(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenAIAssistantConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.openai_assistant.test", "name"),
					resource.TestCheckResourceAttrSet(
						"data.openai_assistant.test", "model"),
					resource.TestCheckResourceAttrSet(
						"data.openai_assistant.test", "created_at"),
					resource.TestCheckResourceAttr(
						"data.openai_assistant.test", "tools.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openai_assistant.test", "tools.0.type", "code_interpreter"),
				),
			},
		},
	})
}

func testAccDataSourceOpenAIAssistantConfig(apiKey string) string {
	return `
provider "openai" {
	api_key = "` + apiKey + `"
}

resource "openai_assistant" "test" {
	name = "Test Assistant"
	model = "gpt-4"
	tools {
		type = "code_interpreter"
	}
}

data "openai_assistant" "test" {
	assistant_id = openai_assistant.test.id
}`
}
