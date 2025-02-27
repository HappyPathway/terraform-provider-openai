package openai

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIThread_basic(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIThreadConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_thread.test", "id"),
					resource.TestCheckResourceAttrSet("openai_thread.test", "created_at"),
					resource.TestCheckResourceAttr("openai_thread.test", "metadata.test", "true"),
					resource.TestCheckResourceAttr("openai_thread.test", "metadata.env", "test"),
				),
			},
		},
	})
}

func TestAccResourceOpenAIThread_withAllTools(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIThreadConfigWithAllTools(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_thread.test_tools", "id"),
					resource.TestCheckResourceAttr("openai_thread.test_tools", "tool_resources.0.code_interpreter.0.file_ids.#", "1"),
					resource.TestCheckResourceAttr("openai_thread.test_tools", "tool_resources.0.file_search.0.vector_store_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceOpenAIThread_update(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIThreadConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_thread.test", "id"),
					resource.TestCheckResourceAttr("openai_thread.test", "metadata.test", "true"),
				),
			},
			{
				Config: testAccResourceOpenAIThreadConfigUpdated(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openai_thread.test", "id"),
					resource.TestCheckResourceAttr("openai_thread.test", "metadata.test", "true"),
					resource.TestCheckResourceAttr("openai_thread.test", "metadata.updated", "yes"),
				),
			},
		},
	})
}

func testAccResourceOpenAIThreadConfig(apiKey string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_thread" "test" {
  metadata = {
    test = "true"
    env  = "test"
  }

  messages {
    role    = "user"
    content = "Hello, this is a test message."
  }
}`, apiKey)
}

func testAccResourceOpenAIThreadConfigWithAllTools(apiKey string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_file" "test" {
  file    = "test.txt"
  purpose = "assistants"
}

resource "openai_thread" "test_tools" {
  metadata = {
    test = "true"
    type = "tools_test"
  }

  messages {
    role    = "user"
    content = "Analyze this data and search through the documentation."
    file_ids = [openai_file.test.id]
  }

  tool_resources {
    code_interpreter {
      file_ids = [openai_file.test.id]
    }

    file_search {
      vector_store_ids = ["test-store-1"]
    }
  }
}`, apiKey)
}

func testAccResourceOpenAIThreadConfigUpdated(apiKey string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_thread" "test" {
  metadata = {
    test    = "true"
    env     = "test"
    updated = "yes"
  }

  messages {
    role    = "user"
    content = "Hello, this is a test message."
  }
}`, apiKey)
}
