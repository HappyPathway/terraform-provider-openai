package datasources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChatCompletionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccChatCompletionDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_chat_completion.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet("data.openai_chat_completion.test", "response_content.0"),
					resource.TestCheckResourceAttrSet("data.openai_chat_completion.test", "response_role.0"),
					resource.TestCheckResourceAttrSet("data.openai_chat_completion.test", "usage.0.total_tokens"),
				),
			},
		},
	})
}

func TestAccChatCompletionDataSource_withTemperature(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccChatCompletionDataSourceConfig_withTemperature(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_chat_completion.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr("data.openai_chat_completion.test", "temperature", "0.7"),
					resource.TestCheckResourceAttrSet("data.openai_chat_completion.test", "response_content.0"),
				),
			},
		},
	})
}

func TestAccChatCompletionDataSource_withSystemRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccChatCompletionDataSourceConfig_withSystemRole(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_chat_completion.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet("data.openai_chat_completion.test", "response_content.0"),
				),
			},
		},
	})
}

func testAccChatCompletionDataSourceConfig_basic() string {
	return `
data "openai_chat_completion" "test" {
  model = "gpt-3.5-turbo"
  messages = [
    {
      role    = "user"
      content = "Say hello!"
    }
  ]
}
`
}

func testAccChatCompletionDataSourceConfig_withTemperature() string {
	return `
data "openai_chat_completion" "test" {
  model = "gpt-3.5-turbo"
  temperature = 0.7
  messages = [
    {
      role    = "user"
      content = "Say hello!"
    }
  ]
}
`
}

func testAccChatCompletionDataSourceConfig_withSystemRole() string {
	return `
data "openai_chat_completion" "test" {
  model = "gpt-3.5-turbo"
  messages = [
    {
      role    = "system"
      content = "You are a helpful assistant that always responds with 'Hello!'"
    },
    {
      role    = "user" 
      content = "Say hello!"
    }
  ]
}
`
}
