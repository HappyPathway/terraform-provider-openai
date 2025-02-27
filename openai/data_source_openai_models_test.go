package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIModels_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenAIModelsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openai_models.test", "models.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openai_models.test", "models.0.id", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(
						"data.openai_models.test", "models.0.owned_by", "openai"),
				),
			},
		},
	})
}

func testAccDataSourceOpenAIModelsConfig() string {
	return `
data "openai_models" "test" {
}
`
}
