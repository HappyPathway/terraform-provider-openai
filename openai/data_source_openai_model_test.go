package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIModel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenAIModelConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openai_model.test", "model_id", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(
						"data.openai_model.test", "owned_by", "openai"),
					resource.TestCheckResourceAttr(
						"data.openai_model.test", "permission.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openai_model.test", "permission.0.allow_fine_tuning", "true"),
				),
			},
		},
	})
}

func testAccDataSourceOpenAIModelConfig() string {
	return `
data "openai_model" "test" {
  model_id = "gpt-3.5-turbo"
}
`
}
