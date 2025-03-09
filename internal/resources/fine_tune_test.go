package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { provider.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFineTuneResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_fine_tune.test", "id"),
					resource.TestCheckResourceAttr("openai_fine_tune.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet("openai_fine_tune.test", "status"),
					resource.TestCheckResourceAttrSet("openai_fine_tune.test", "created_at"),
				),
			},
		},
	})
}

func TestAccFineTuneResource_withHyperparameters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { provider.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFineTuneResourceConfig_withHyperparameters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_fine_tune.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr("openai_fine_tune.test", "hyperparameters.n_epochs", "3"),
					resource.TestCheckResourceAttrSet("openai_fine_tune.test", "status"),
				),
			},
		},
	})
}

func testAccFineTuneResourceConfig_basic() string {
	return `
resource "openai_file" "test" {
  filename = "test.jsonl"
  content  = jsonencode([
    {"prompt": "Company: ACME Corp, Revenue: ", "completion": " $1M"},
    {"prompt": "Company: XYZ Inc, Revenue: ", "completion": " $2M"}
  ])
  purpose = "fine-tune"
}

resource "openai_fine_tune" "test" {
  training_file = openai_file.test.id
  model        = "gpt-3.5-turbo"
  suffix       = "test-model"
}
`
}

func testAccFineTuneResourceConfig_withHyperparameters() string {
	return `
resource "openai_file" "test" {
  filename = "test.jsonl"
  content  = jsonencode([
    {"prompt": "Company: ACME Corp, Revenue: ", "completion": " $1M"},
    {"prompt": "Company: XYZ Inc, Revenue: ", "completion": " $2M"}
  ])
  purpose = "fine-tune"
}

resource "openai_fine_tune" "test" {
  training_file = openai_file.test.id
  model        = "gpt-3.5-turbo"
  suffix       = "test-model-params"
  hyperparameters = {
    n_epochs = 3
  }
}
`
}
