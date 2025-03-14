package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFineTuneResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_fine_tune.test", "model", "davinci"),
					resource.TestCheckResourceAttrSet("openai_fine_tune.test", "id"),
					resource.TestCheckResourceAttrSet("openai_fine_tune.test", "created_at"),
					resource.TestCheckResourceAttr("openai_fine_tune.test", "training_file_id", "file-abc123"),
				),
			},
		},
	})
}

func TestAccFineTuneResource_withHyperparameters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFineTuneResourceConfig_withHyperparameters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_fine_tune.test", "model", "davinci"),
					resource.TestCheckResourceAttr("openai_fine_tune.test", "hyperparameters.n_epochs", "4"),
					resource.TestCheckResourceAttr("openai_fine_tune.test", "hyperparameters.batch_size", "32"),
				),
			},
		},
	})
}

func testAccFineTuneResourceConfig_basic() string {
	return `
resource "openai_fine_tune" "test" {
  training_file_id = "file-abc123"
  model           = "davinci"
}
`
}

func testAccFineTuneResourceConfig_withHyperparameters() string {
	return `
resource "openai_fine_tune" "test" {
  training_file_id = "file-abc123"
  model          = "davinci"
  hyperparameters = {
    n_epochs   = 4
    batch_size = 32
  }
}
`
}
