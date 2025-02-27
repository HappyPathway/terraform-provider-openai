package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceModeration_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceModerationConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_moderation.test", "input", "I want to create art"),
					resource.TestCheckResourceAttr("openai_moderation.test", "model", "text-moderation-latest"),
					resource.TestCheckResourceAttrSet("openai_moderation.test", "flagged"),
					resource.TestCheckResourceAttrSet("openai_moderation.test", "categories.0.harassment"),
					resource.TestCheckResourceAttrSet("openai_moderation.test", "category_scores.0.harassment"),
				),
			},
		},
	})
}

func TestAccResourceModeration_sensitive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceModerationConfig_sensitive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_moderation.test_sensitive", "input", "I hate everyone"),
					resource.TestCheckResourceAttr("openai_moderation.test_sensitive", "model", "text-moderation-stable"),
					resource.TestCheckResourceAttrSet("openai_moderation.test_sensitive", "flagged"),
					resource.TestCheckResourceAttrSet("openai_moderation.test_sensitive", "categories.0.hate"),
					resource.TestCheckResourceAttrSet("openai_moderation.test_sensitive", "category_scores.0.hate"),
				),
			},
		},
	})
}

func TestAccResourceModeration_multipleInputs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceModerationConfig_multipleInputs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_moderation.test_multi", "input.0", "Hello world"),
					resource.TestCheckResourceAttr("openai_moderation.test_multi", "input.1", "I love coding"),
					resource.TestCheckResourceAttr("openai_moderation.test_multi", "model", "text-moderation-latest"),
					resource.TestCheckResourceAttrSet("openai_moderation.test_multi", "results.#"),
					resource.TestCheckResourceAttrSet("openai_moderation.test_multi", "results.0.flagged"),
					resource.TestCheckResourceAttrSet("openai_moderation.test_multi", "results.1.flagged"),
				),
			},
		},
	})
}

func testAccResourceModerationConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_moderation" "test" {
  input = "I want to create art"
  model = "text-moderation-latest"
}
`)
}

func testAccResourceModerationConfig_sensitive() string {
	return fmt.Sprintf(`
resource "openai_moderation" "test_sensitive" {
  input = "I hate everyone"
  model = "text-moderation-stable"
}
`)
}

func testAccResourceModerationConfig_multipleInputs() string {
	return fmt.Sprintf(`
resource "openai_moderation" "test_multi" {
  input = ["Hello world", "I love coding"]
  model = "text-moderation-latest"
}
`)
}
