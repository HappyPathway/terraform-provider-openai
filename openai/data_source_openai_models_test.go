package openai

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIModels_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + testAccDataSourceOpenAIModelsConfig(),
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

func TestAccDataSourceModels_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceModelsConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Check that we get a list of models
					resource.TestCheckResourceAttrSet("data.openai_models.all", "models.#"),
					// Check that the models belong to OpenAI
					resource.TestCheckResourceAttr("data.openai_models.all", "models.0.owned_by", "openai"),
					// Check that permissions are populated
					resource.TestCheckResourceAttrSet("data.openai_models.all", "models.0.permission.#"),
				),
			},
		},
	})
}

func TestAccDataSourceModels_filtered(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceModelsConfig_filtered(),
				Check: resource.ComposeTestCheckFunc(
					// Verify we get results and they're filtered
					resource.TestCheckResourceAttrSet("data.openai_models.gpt", "models.#"),
					resource.TestMatchResourceAttr("data.openai_models.gpt", "models.0.id",
						regexp.MustCompile(`^gpt-`)),
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

func testAccDataSourceModelsConfig_basic() string {
	return `
data "openai_models" "all" {}
`
}

func testAccDataSourceModelsConfig_filtered() string {
	return `
data "openai_models" "gpt" {
  filter_prefix = "gpt-"
}
`
}
