package resources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEmbeddingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccEmbeddingResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_embedding.test", "model", "text-embedding-ada-002"),
					resource.TestCheckResourceAttrSet("openai_embedding.test", "embedding.#"),
					resource.TestCheckResourceAttrSet("openai_embedding.test", "usage.0.total_tokens"),
				),
			},
		},
	})
}

func testAccEmbeddingResourceConfig_basic() string {
	return `
resource "openai_embedding" "test" {
  model = "text-embedding-ada-002"
  input = "This is a test input for embedding generation."
}
`
}
