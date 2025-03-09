package datasources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVectorStoreDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccVectorStoreDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_vector_store.test", "name", "Test Store"),
					resource.TestCheckResourceAttrSet("data.openai_vector_store.test", "created_at"),
				),
			},
		},
	})
}

func TestAccVectorStoreDataSource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccVectorStoreDataSourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_vector_store.test", "name", "Test Store with Metadata"),
					resource.TestCheckResourceAttr("data.openai_vector_store.test", "metadata.project", "test"),
					resource.TestCheckResourceAttr("data.openai_vector_store.test", "metadata.environment", "acceptance"),
				),
			},
		},
	})
}

func testAccVectorStoreDataSourceConfig_basic() string {
	return `
resource "openai_vector_store" "test" {
  name = "Test Store"
}

data "openai_vector_store" "test" {
  id = openai_vector_store.test.id
}
`
}

func testAccVectorStoreDataSourceConfig_withMetadata() string {
	return `
resource "openai_vector_store" "test" {
  name = "Test Store with Metadata"
  metadata = {
    project     = "test"
    environment = "acceptance"
  }
}

data "openai_vector_store" "test" {
  id = openai_vector_store.test.id
}
`
}
