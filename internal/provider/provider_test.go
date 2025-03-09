package provider_test

import (
	"os"
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `provider "openai" {}`,
			},
		},
	})
}

func TestAccProviderConfigure(t *testing.T) {
	t.Run("can be configured with API token in config", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
						provider "openai" {
							api_token = "test-token"
						}
					`,
				},
			},
		})
	})

	t.Run("can be configured with API token from environment", func(t *testing.T) {
		os.Setenv("OPENAI_API_TOKEN", "test-token")
		defer os.Unsetenv("OPENAI_API_TOKEN")

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `provider "openai" {}`,
				},
			},
		})
	})

	t.Run("can be configured with API token and base URL", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
						provider "openai" {
							api_token = "test-token"
							base_url = "https://api.custom-openai.example"
						}
					`,
				},
			},
		})
	})
}

func testAccPreCheck(t *testing.T) {
	// You can add any necessary setup for acceptance tests here
	if os.Getenv("OPENAI_API_TOKEN") == "" {
		t.Skip("OPENAI_API_TOKEN environment variable must be set for acceptance tests")
	}
}
