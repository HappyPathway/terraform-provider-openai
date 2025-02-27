package openai

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIModel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + testAccDataSourceOpenAIModelConfig(),
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

func TestAccDataSourceModel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceModelConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_model.gpt4", "id", "gpt-4"),
					resource.TestCheckResourceAttr("data.openai_model.gpt4", "owned_by", "openai"),
				),
			},
		},
	})
}

func TestAccDataSourceModel_nonexistent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceModelConfig_nonexistent(),
				ExpectError: regexp.MustCompile("Model not found"),
			},
		},
	})
}

const testAccProviderConfig = `
provider "openai" {
  api_key = "mock-api-key"
}
`

func testAccDataSourceOpenAIModelConfig() string {
	return `
data "openai_model" "test" {
  model_id = "gpt-3.5-turbo"
}
`
}

func testAccDataSourceModelConfig_basic() string {
	return fmt.Sprintf(`
data "openai_model" "gpt4" {
  model = "gpt-4"
}
`)
}

func testAccDataSourceModelConfig_nonexistent() string {
	return fmt.Sprintf(`
data "openai_model" "nonexistent" {
  model = "nonexistent-model"
}
`)
}
