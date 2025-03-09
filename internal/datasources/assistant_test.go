package datasources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssistantDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create an assistant using the resource
			{
				Config: acctest.ProviderConfig() + testAccAssistantDataSourceConfig_resource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test", "name", "Test Assistant"),
				),
			},
			// Then test reading it with the data source
			{
				Config: acctest.ProviderConfig() + testAccAssistantDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_assistant.test", "name", "Test Assistant"),
					resource.TestCheckResourceAttr("data.openai_assistant.test", "model", "gpt-4"),
					resource.TestCheckResourceAttrSet("data.openai_assistant.test", "created_at"),
				),
			},
		},
	})
}

func TestAccAssistantDataSource_withTools(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccAssistantDataSourceConfig_withTools(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_assistant.test", "name", "Test Assistant with Tools"),
					resource.TestCheckResourceAttr("data.openai_assistant.test", "tools.#", "1"),
					resource.TestCheckResourceAttr("data.openai_assistant.test", "tools.0", "code_interpreter"),
				),
			},
		},
	})
}

func testAccAssistantDataSourceConfig_resource() string {
	return `
resource "openai_assistant" "test" {
  name         = "Test Assistant"
  description  = "A test assistant"
  model        = "gpt-4"
  instructions = "You are a helpful test assistant."
}
`
}

func testAccAssistantDataSourceConfig_basic() string {
	return `
resource "openai_assistant" "test" {
  name         = "Test Assistant"
  description  = "A test assistant"
  model        = "gpt-4"
  instructions = "You are a helpful test assistant."
}

data "openai_assistant" "test" {
  assistant_id = openai_assistant.test.id
}
`
}

func testAccAssistantDataSourceConfig_withTools() string {
	return `
resource "openai_assistant" "test" {
  name         = "Test Assistant with Tools"
  description  = "A test assistant with tools"
  model        = "gpt-4"
  instructions = "You are a helpful test assistant that can write code."
  tools        = ["code_interpreter"]
}

data "openai_assistant" "test" {
  assistant_id = openai_assistant.test.id
}
`
}
