package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChatCompletionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChatCompletionResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_chat_completion.test", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.0.role", "user"),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.0.content", "Say hello!"),
				),
			},
		},
	})
}

func TestAccChatCompletionResource_withTemperature(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChatCompletionResourceConfig_withTemperature(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_chat_completion.test", "temperature", "0.7"),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.0.role", "user"),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.0.content", "Say hello with creativity!"),
				),
			},
		},
	})
}

func TestAccChatCompletionResource_withSystemRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChatCompletionResourceConfig_withSystemRole(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.0.role", "system"),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.0.content", "You are a helpful assistant."),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.1.role", "user"),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "messages.1.content", "Say hello!"),
				),
			},
		},
	})
}

func testAccChatCompletionResourceConfig_basic() string {
	return `
resource "openai_chat_completion" "test" {
  model = "gpt-4"
  messages = [
    {
      role    = "user"
      content = "Say hello!"
    }
  ]
}
`
}

func testAccChatCompletionResourceConfig_withTemperature() string {
	return `
resource "openai_chat_completion" "test" {
  model = "gpt-4"
  temperature = 0.7
  messages = [
    {
      role    = "user"
      content = "Say hello with creativity!"
    }
  ]
}
`
}

func testAccChatCompletionResourceConfig_withSystemRole() string {
	return `
resource "openai_chat_completion" "test" {
  model = "gpt-4"
  messages = [
    {
      role    = "system"
      content = "You are a helpful assistant."
    },
    {
      role    = "user"
      content = "Say hello!"
    }
  ]
}
`
}
