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
				Config: acctest.ProviderConfig() + testAccChatCompletionResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_chat_completion.test", "model", "gpt-4"),
					resource.TestCheckResourceAttrSet("openai_chat_completion.test", "created_at"),
					resource.TestCheckResourceAttrSet("openai_chat_completion.test", "choices.0.message.content"),
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
				Config: acctest.ProviderConfig() + testAccChatCompletionResourceConfig_withTemperature(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_chat_completion.test", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_chat_completion.test", "temperature", "0.7"),
					resource.TestCheckResourceAttrSet("openai_chat_completion.test", "choices.0.message.content"),
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
				Config: acctest.ProviderConfig() + testAccChatCompletionResourceConfig_withSystemRole(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_chat_completion.test", "model", "gpt-4"),
					resource.TestCheckResourceAttrSet("openai_chat_completion.test", "choices.0.message.content"),
				),
			},
		},
	})
}

func testAccChatCompletionResourceConfig_basic() string {
	return `
resource "openai_chat_completion" "test" {
  model = "gpt-4"
  message {
    role    = "user"
    content = "Hello!"
  }
}
`
}

func testAccChatCompletionResourceConfig_withTemperature() string {
	return `
resource "openai_chat_completion" "test" {
  model       = "gpt-4"
  temperature = 0.7
  message {
    role    = "user"
    content = "Tell me a creative story."
  }
}
`
}

func testAccChatCompletionResourceConfig_withSystemRole() string {
	return `
resource "openai_chat_completion" "test" {
  model = "gpt-4"
  message {
    role    = "system"
    content = "You are a helpful assistant."
  }
  message {
    role    = "user"
    content = "Hello!"
  }
}
`
}
