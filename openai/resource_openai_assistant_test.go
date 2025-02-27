package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIAssistant_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIAssistantConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "name", "Test Assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "description", "Test assistant for acceptance tests"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "instructions", "You are a test assistant."),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.#", "1"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.0.type", "code_interpreter"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "metadata.test", "true"),
				),
			},
			{
				Config: testAccResourceOpenAIAssistantConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "name", "Updated Test Assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "description", "Updated test assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.#", "0"),
				),
			},
		},
	})
}

func testAccResourceOpenAIAssistantConfig() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test" {
  name         = "Test Assistant"
  description  = "Test assistant for acceptance tests"
  model        = "gpt-3.5-turbo"
  instructions = "You are a test assistant."

  tools {
    type = "code_interpreter"
  }

  metadata = {
    test = "true"
  }
}
`)
}

func testAccResourceOpenAIAssistantConfigUpdated() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test" {
  name         = "Updated Test Assistant"
  description  = "Updated test assistant"
  model        = "gpt-3.5-turbo"
  instructions = "You are a test assistant."

  metadata = {
    test = "true"
  }
}
`)
}
