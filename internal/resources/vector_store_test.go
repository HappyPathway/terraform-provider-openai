package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVectorStoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_vector_store.test", "name", "test-store"),
					resource.TestCheckResourceAttr("openai_vector_store.test", "expires_after.days", "90"),
				),
			},
		},
	})
}

func TestAccVectorStoreResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreResourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_vector_store.test", "name", "test-store-metadata"),
					resource.TestCheckResourceAttr("openai_vector_store.test", "metadata.project", "test"),
					resource.TestCheckResourceAttr("openai_vector_store.test", "metadata.environment", "acceptance"),
					resource.TestCheckResourceAttr("openai_vector_store.test", "expires_after.days", "90"),
				),
			},
		},
	})
}

func testAccVectorStoreResourceConfig_basic() string {
	return `
resource "openai_vector_store" "test" {
  name = "test-store"
  expires_after = {
    days = 90
  }
}
`
}

func testAccVectorStoreResourceConfig_withMetadata() string {
	return `
resource "openai_vector_store" "test" {
  name = "test-store-metadata"
  expires_after = {
    days = 90
  }
  metadata = {
    project     = "test"
    environment = "acceptance"
  }
}
`
}
