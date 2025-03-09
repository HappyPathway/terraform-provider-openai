package datasources_test

import (
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccModelDataSourceConfig_gpt4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_model.test", "model_id", "gpt-4"),
					resource.TestCheckResourceAttr("data.openai_model.test", "owned_by", "openai"),
					resource.TestCheckResourceAttrSet("data.openai_model.test", "created"),
				),
			},
		},
	})
}

func TestAccModelDataSource_ada002(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig() + testAccModelDataSourceConfig_ada002(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openai_model.test", "model_id", "text-embedding-ada-002"),
					resource.TestCheckResourceAttr("data.openai_model.test", "owned_by", "openai"),
					resource.TestCheckResourceAttrSet("data.openai_model.test", "created"),
				),
			},
		},
	})
}

func testAccModelDataSourceConfig_gpt4() string {
	return `
data "openai_model" "test" {
  model_id = "gpt-4"
}
`
}

func testAccModelDataSourceConfig_ada002() string {
	return `
data "openai_model" "test" {
  model_id = "text-embedding-ada-002"
}
`
}
