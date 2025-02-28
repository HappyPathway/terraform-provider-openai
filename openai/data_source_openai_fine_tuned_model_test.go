package openai

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"openai": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func TestAccDataSourceOpenAIFineTunedModel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenAIFineTunedModelConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.openai_fine_tuned_model.test", "id", regexp.MustCompile("^ft:")),
					resource.TestCheckResourceAttrSet(
						"data.openai_fine_tuned_model.test", "created"),
					resource.TestCheckResourceAttrSet(
						"data.openai_fine_tuned_model.test", "owned_by"),
					resource.TestCheckResourceAttrSet(
						"data.openai_fine_tuned_model.test", "object"),
				),
			},
		},
	})
}

const testAccDataSourceOpenAIFineTunedModelConfig = `
data "openai_fine_tuned_model" "test" {
  model_id = "ft:gpt-3.5-turbo-0613:personal::8Da2zh9P" # Replace with a valid fine-tuned model ID in your test environment
}
`
