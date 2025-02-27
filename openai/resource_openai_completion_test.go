package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCompletion_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCompletionConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_completion.test", "model", "gpt-3.5-turbo-instruct"),
					resource.TestCheckResourceAttr("openai_completion.test", "prompt", "Say hello"),
					resource.TestCheckResourceAttr("openai_completion.test", "max_tokens", "50"),
					resource.TestCheckResourceAttrSet("openai_completion.test", "choices.0.text"),
				),
			},
		},
	})
}

func TestAccResourceCompletion_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCompletionConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_completion.test_full", "model", "gpt-3.5-turbo-instruct"),
					resource.TestCheckResourceAttr("openai_completion.test_full", "prompt", "Write a story about AI"),
					resource.TestCheckResourceAttr("openai_completion.test_full", "max_tokens", "100"),
					resource.TestCheckResourceAttr("openai_completion.test_full", "temperature", "0.8"),
					resource.TestCheckResourceAttr("openai_completion.test_full", "top_p", "0.9"),
					resource.TestCheckResourceAttr("openai_completion.test_full", "n", "2"),
					resource.TestCheckResourceAttr("openai_completion.test_full", "stop.0", "THE END"),
					resource.TestCheckResourceAttrSet("openai_completion.test_full", "choices.0.text"),
					resource.TestCheckResourceAttrSet("openai_completion.test_full", "choices.1.text"),
				),
			},
		},
	})
}

func testAccResourceCompletionConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_completion" "test" {
  model      = "gpt-3.5-turbo-instruct"
  prompt     = "Say hello"
  max_tokens = 50
}
`)
}

func testAccResourceCompletionConfig_full() string {
	return fmt.Sprintf(`
resource "openai_completion" "test_full" {
  model      = "gpt-3.5-turbo-instruct"
  prompt     = "Write a story about AI"
  max_tokens = 100
  temperature = 0.8
  top_p      = 0.9
  n          = 2
  stop       = ["THE END"]
}
`)
}
