package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMessageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccMessageResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_message.test", "id"),
					resource.TestCheckResourceAttr("openai_message.test", "role", "user"),
					resource.TestCheckResourceAttr("openai_message.test", "content", "Hello!"),
					resource.TestCheckResourceAttrSet("openai_message.test", "created_at"),
				),
			},
		},
	})
}

func TestAccMessageResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccMessageResourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_message.test", "role", "user"),
					resource.TestCheckResourceAttr("openai_message.test", "metadata.user_id", "test-123"),
					resource.TestCheckResourceAttr("openai_message.test", "metadata.type", "greeting"),
				),
			},
		},
	})
}

func testAccMessageResourceConfig_basic() string {
	return `
resource "openai_thread" "test" {}

resource "openai_message" "test" {
  thread_id = openai_thread.test.id
  role      = "user"
  content   = "Hello!"
}
`
}

func testAccMessageResourceConfig_withMetadata() string {
	return `
resource "openai_thread" "test" {}

resource "openai_message" "test" {
  thread_id = openai_thread.test.id
  role      = "user"
  content   = "Hello with metadata!"
  metadata = {
    user_id = "test-123"
    type    = "greeting"
  }
}
`
}
