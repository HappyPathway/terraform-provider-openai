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
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.#", "1"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.0.type", "code_interpreter"),
				),
			},
		},
	})
}

func TestAccResourceOpenAIAssistant_withAllTools(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIAssistantConfigWithAllTools(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_assistant.test_tools", "name", "Tool Test Assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test_tools", "tools.#", "3"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test_tools", "tools.0.type", "code_interpreter"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test_tools", "tools.1.type", "retrieval"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test_tools", "tools.2.type", "function"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test_tools", "tools.2.function.name", "get_weather"),
				),
			},
		},
	})
}

func TestAccResourceOpenAIAssistant_update(t *testing.T) {
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
				),
			},
			{
				Config: testAccResourceOpenAIAssistantConfigUpdated(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "name", "Updated Test Assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "model", "gpt-4-1106-preview"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "instructions", "You are an updated test assistant."),
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

func testAccResourceOpenAIAssistantConfigWithAllTools(apiKey string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_assistant" "test_tools" {
  name         = "Tool Test Assistant"
  description  = "Test assistant with all tool types"
  model        = "gpt-4"
  instructions = "You are a test assistant with all available tools."

  tools {
    type = "code_interpreter"
  }

  tools {
    type = "retrieval"
  }

  tools {
    type = "function"
    name = "get_weather"
    description = "Get the current weather in a location"
    parameters = jsonencode({
      type = "object"
      properties = {
        location = {
          type = "string",
          description = "The city and state, e.g. San Francisco, CA"
        }
        unit = {
          type = "string",
          enum = ["celsius", "fahrenheit"]
        }
      },
      required = ["location"]
    })
  }

  metadata = {
    test = "true"
    type = "full_tools"
  }
}
`, apiKey)
}

func testAccResourceOpenAIAssistantConfigUpdated(apiKey string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_assistant" "test" {
  name         = "Updated Test Assistant"
  description  = "Updated test assistant for acceptance tests"
  model        = "gpt-4-1106-preview"
  instructions = "You are an updated test assistant."

  tools {
    type = "code_interpreter"
  }

  tools {
    type = "retrieval"
  }

  metadata = {
    test = "true"
    updated = "yes"
  }
}
`, apiKey)
}
