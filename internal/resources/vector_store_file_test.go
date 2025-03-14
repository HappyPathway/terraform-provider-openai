package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVectorStoreFileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreFileResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_vector_store_file.test", "id"),
					resource.TestCheckResourceAttrSet("openai_vector_store_file.test", "created_at"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "vector_store_id", "vs_abc123"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "file_id", "file-abc123"),
				),
			},
		},
	})
}

func TestAccVectorStoreFileResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreFileResourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_vector_store_file.test", "id"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "vector_store_id", "vs_abc123"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "file_id", "file-abc123"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "metadata.department", "engineering"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "metadata.type", "documentation"),
				),
			},
		},
	})
}

func testAccVectorStoreFileResourceConfig_basic() string {
	return `
resource "openai_vector_store_file" "test" {
  vector_store_id = "vs_abc123"
  file_id        = "file-abc123"
}
`
}

func testAccVectorStoreFileResourceConfig_withMetadata() string {
	return `
resource "openai_vector_store_file" "test" {
  vector_store_id = "vs_abc123"
  file_id        = "file-abc123"
  metadata = {
    department = "engineering"
    type      = "documentation"
  }
}
`
}
