package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssistantResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccAssistantResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test", "name", "Test Assistant"),
					resource.TestCheckResourceAttr("openai_assistant.test", "model", "gpt-4"),
					resource.TestCheckResourceAttrSet("openai_assistant.test", "created_at"),
				),
			},
		},
	})
}

func TestAccAssistantResource_withTools(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccAssistantResourceConfig_withTools(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test", "name", "Test Assistant with Tools"),
					resource.TestCheckResourceAttr("openai_assistant.test", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_assistant.test", "tools.#", "2"),
					resource.TestCheckResourceAttr("openai_assistant.test", "tools.0.type", "code_interpreter"),
					resource.TestCheckResourceAttr("openai_assistant.test", "tools.1.type", "retrieval"),
				),
			},
		},
	})
}

func TestAccAssistantResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccAssistantResourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test", "name", "Test Assistant with Metadata"),
					resource.TestCheckResourceAttr("openai_assistant.test", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_assistant.test", "metadata.project", "test"),
					resource.TestCheckResourceAttr("openai_assistant.test", "metadata.environment", "acceptance"),
				),
			},
		},
	})
}

func testAccAssistantResourceConfig_basic() string {
	return `
resource "openai_assistant" "test" {
  name  = "Test Assistant"
  model = "gpt-4"
}
`
}

func testAccAssistantResourceConfig_withTools() string {
	return `
resource "openai_assistant" "test" {
  name  = "Test Assistant with Tools"
  model = "gpt-4"
  tools = [
    {
      type = "code_interpreter"
    },
    {
      type = "retrieval"
    }
  ]
}
`
}

func testAccAssistantResourceConfig_withMetadata() string {
	return `
resource "openai_assistant" "test" {
  name  = "Test Assistant with Metadata"
  model = "gpt-4"
  metadata = {
    project     = "test"
    environment = "acceptance"
  }
}
`
}
