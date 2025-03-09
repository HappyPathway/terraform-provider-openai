package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccFileResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_file.test", "filename", "test.jsonl"),
					resource.TestCheckResourceAttr("openai_file.test", "purpose", "fine-tune"),
					resource.TestCheckResourceAttrSet("openai_file.test", "bytes"),
					resource.TestCheckResourceAttrSet("openai_file.test", "created_at"),
				),
			},
		},
	})
}

func TestAccFileResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccFileResourceConfig_withMetadata(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_file.test", "filename", "test.jsonl"),
					resource.TestCheckResourceAttr("openai_file.test", "purpose", "fine-tune"),
					resource.TestCheckResourceAttr("openai_file.test", "metadata.project", "test"),
					resource.TestCheckResourceAttr("openai_file.test", "metadata.environment", "acceptance"),
				),
			},
		},
	})
}

func testAccFileResourceConfig_basic() string {
	return `
resource "openai_file" "test" {
  filename = "test.jsonl"
  purpose  = "fine-tune"
  content  = "[{\"prompt\": \"Test prompt\", \"completion\": \"Test completion\"}]"
}
`
}

func testAccFileResourceConfig_withMetadata() string {
	return `
resource "openai_file" "test" {
  filename = "test.jsonl"
  purpose  = "fine-tune"
  content  = "[{\"prompt\": \"Test prompt\", \"completion\": \"Test completion\"}]"
  metadata = {
    project     = "test"
    environment = "acceptance"
  }
}
`
}
