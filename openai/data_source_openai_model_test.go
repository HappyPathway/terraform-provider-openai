package openai

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIModel_basic(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenAIModelConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openai_model.test", "model_id", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet(
						"data.openai_model.test", "owned_by"),
					resource.TestCheckResourceAttrSet(
						"data.openai_model.test", "permission.#"),
				),
			},
		},
	})
}

func testAccDataSourceOpenAIModelConfig(apiKey string) string {
	return `
provider "openai" {
  api_key = "` + apiKey + `"
}

data "openai_model" "test" {
  model_id = "gpt-3.5-turbo"
}
`
}
