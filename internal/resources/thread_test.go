package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccThreadResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccThreadResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_thread.test", "id"),
					resource.TestCheckResourceAttrSet("openai_thread.test", "created_at"),
				),
			},
		},
	})
}

func TestAccThreadResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccThreadResourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_thread.test", "metadata.user_id", "test-123"),
					resource.TestCheckResourceAttr("openai_thread.test", "metadata.type", "test"),
				),
			},
		},
	})
}

func testAccThreadResourceConfig_basic() string {
	return `
resource "openai_thread" "test" {}
`
}

func testAccThreadResourceConfig_withMetadata() string {
	return `
resource "openai_thread" "test" {
  metadata = {
    user_id = "test-123"
    type    = "test"
  }
}
`
}
