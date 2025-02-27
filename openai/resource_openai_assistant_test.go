package openai

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIAssistant_basic(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIAssistantConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "name", "Test Assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "model", "gpt-4"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "description", "Test assistant for acceptance tests"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "instructions", "You are a test assistant."),
					resource.TestCheckResourceAttrSet(
						"openai_assistant.test", "created_at"),
				),
			},
		},
	})
}

func testAccResourceOpenAIAssistantConfig(apiKey string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_assistant" "test" {
  name         = "Test Assistant"
  description  = "Test assistant for acceptance tests"
  model        = "gpt-4"
  instructions = "You are a test assistant."

  tools {
    type = "code_interpreter"
  }

  metadata = {
    test = "true"
  }
}
`, apiKey)
}
