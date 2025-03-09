package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVectorStoreFileResource(t *testing.T) {
	// Create a temporary test file
	content := []byte("This is test content for vector embedding")
	tmpfile, err := os.CreateTemp("", "vector-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { provider.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreFileResourceConfig_basic(tmpfile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_vector_store_file.test", "id"),
					resource.TestCheckResourceAttrSet("openai_vector_store_file.test", "created_at"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "filename", "test.txt"),
					resource.TestCheckResourceAttrSet("openai_vector_store_file.test", "bytes"),
				),
			},
		},
	})
}

func TestAccVectorStoreFileResource_withMetadata(t *testing.T) {
	content := []byte("This is test content for vector embedding with metadata")
	tmpfile, err := os.CreateTemp("", "vector-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { provider.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccVectorStoreFileResourceConfig_withMetadata(tmpfile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_vector_store_file.test", "id"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "metadata.source", "acceptance-test"),
					resource.TestCheckResourceAttr("openai_vector_store_file.test", "metadata.type", "document"),
				),
			},
		},
	})
}

func testAccVectorStoreFileResourceConfig_basic(filepath string) string {
	return fmt.Sprintf(`
resource "openai_vector_store" "test" {
  name        = "test-store"
  description = "Test vector store"
}

resource "openai_vector_store_file" "test" {
  vector_store_id = openai_vector_store.test.id
  filename       = "test.txt"
  filepath       = "%s"
}
`, filepath)
}

func testAccVectorStoreFileResourceConfig_withMetadata(filepath string) string {
	return fmt.Sprintf(`
resource "openai_vector_store" "test" {
  name        = "test-store"
  description = "Test vector store"
}

resource "openai_vector_store_file" "test" {
  vector_store_id = openai_vector_store.test.id
  filename       = "test.txt"
  filepath       = "%s"
  metadata = {
    source = "acceptance-test"
    type   = "document"
  }
}
`, filepath)
}
