package openai

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIModels_basic(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenAIModelsConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.openai_models.test", "models.#"),
				),
			},
		},
	})
}

func testAccDataSourceOpenAIModelsConfig(apiKey string) string {
	return `
provider "openai" {
  api_key = "` + apiKey + `"
}

data "openai_models" "test" {
}
`
}
