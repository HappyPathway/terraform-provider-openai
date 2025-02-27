package openai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEmbedding_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEmbeddingConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_embedding.test", "model", "text-embedding-ada-002"),
					resource.TestCheckResourceAttr("openai_embedding.test", "input", "Sample text for embedding"),
					resource.TestCheckResourceAttr("openai_embedding.test", "encoding_format", "float"),
					resource.TestCheckResourceAttrSet("openai_embedding.test", "embeddings.0.embedding.#"),
					resource.TestCheckResourceAttrSet("openai_embedding.test", "embeddings.0.index"),
					resource.TestCheckResourceAttrSet("openai_embedding.test", "usage.prompt_tokens"),
					resource.TestCheckResourceAttrSet("openai_embedding.test", "usage.total_tokens"),
				),
			},
		},
	})
}

func TestAccResourceEmbedding_multipleInputs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEmbeddingConfig_multipleInputs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_embedding.test_multi", "model", "text-embedding-ada-002"),
					resource.TestCheckResourceAttr("openai_embedding.test_multi", "input.0", "First sample text"),
					resource.TestCheckResourceAttr("openai_embedding.test_multi", "input.1", "Second sample text"),
					resource.TestCheckResourceAttr("openai_embedding.test_multi", "encoding_format", "float"),
					resource.TestCheckResourceAttrSet("openai_embedding.test_multi", "embeddings.0.embedding.#"),
					resource.TestCheckResourceAttrSet("openai_embedding.test_multi", "embeddings.1.embedding.#"),
					resource.TestCheckResourceAttr("openai_embedding.test_multi", "embeddings.0.index", "0"),
					resource.TestCheckResourceAttr("openai_embedding.test_multi", "embeddings.1.index", "1"),
					resource.TestCheckResourceAttrSet("openai_embedding.test_multi", "usage.prompt_tokens"),
					resource.TestCheckResourceAttrSet("openai_embedding.test_multi", "usage.total_tokens"),
				),
			},
		},
	})
}

func testAccResourceEmbeddingConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_embedding" "test" {
  model           = "text-embedding-ada-002"
  input           = "Sample text for embedding"
  encoding_format = "float"
}
`)
}

func testAccResourceEmbeddingConfig_multipleInputs() string {
	return fmt.Sprintf(`
resource "openai_embedding" "test_multi" {
  model           = "text-embedding-ada-002"
  input           = ["First sample text", "Second sample text"]
  encoding_format = "float"
}
`)
}
