package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVectorStoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { provider.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_vector_store.test", "name", "test-store"),
					resource.TestCheckResourceAttr("openai_vector_store.test", "description", "Test vector store"),
					resource.TestCheckResourceAttrSet("openai_vector_store.test", "id"),
					resource.TestCheckResourceAttrSet("openai_vector_store.test", "created_at"),
				),
			},
		},
	})
}

func TestAccVectorStoreResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { provider.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreResourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_vector_store.test", "name", "test-store-metadata"),
					resource.TestCheckResourceAttr("openai_vector_store.test", "metadata.project", "test"),
					resource.TestCheckResourceAttr("openai_vector_store.test", "metadata.environment", "acceptance"),
				),
			},
		},
	})
}

func testAccVectorStoreResourceConfig_basic() string {
	return `
resource "openai_vector_store" "test" {
  name        = "test-store"
  description = "Test vector store"
}
`
}

func testAccVectorStoreResourceConfig_withMetadata() string {
	return `
resource "openai_vector_store" "test" {
  name        = "test-store-metadata"
  description = "Test vector store with metadata"
  metadata = {
    project     = "test"
    environment = "acceptance"
  }
}
`
}
