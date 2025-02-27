package openai

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIFile_basic(t *testing.T) {
	testDataDir := "testdata"
	testDataPath := filepath.Join(testDataDir, "test.jsonl")

	// Ensure test data directory exists
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test file
	testData := `{"messages": [{"role": "system", "content": "You are a test assistant."}, {"role": "user", "content": "Hello"}, {"role": "assistant", "content": "Hi there!"}]}`
	if err := os.WriteFile(testDataPath, []byte(testData+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDataPath)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + testAccResourceOpenAIFileConfig(testDataPath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_file.test", "purpose", "fine-tune"),
					resource.TestCheckResourceAttrSet(
						"openai_file.test", "bytes"),
					resource.TestCheckResourceAttrSet(
						"openai_file.test", "filename"),
				),
			},
		},
	})
}

func TestAccResourceFile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFileConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_file.test", "purpose", "fine-tune"),
					resource.TestCheckResourceAttrSet("openai_file.test", "filename"),
					resource.TestCheckResourceAttrSet("openai_file.test", "id"),
					resource.TestCheckResourceAttrSet("openai_file.test", "bytes"),
					resource.TestCheckResourceAttrSet("openai_file.test", "created_at"),
					resource.TestCheckResourceAttr("openai_file.test", "status", "processed"),
				),
			},
			{
				ResourceName:      "openai_file.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceFile_assistantFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFileConfig_assistantFile(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_file.test_assistant", "purpose", "assistants"),
					resource.TestCheckResourceAttrSet("openai_file.test_assistant", "filename"),
					resource.TestCheckResourceAttrSet("openai_file.test_assistant", "id"),
					resource.TestCheckResourceAttrSet("openai_file.test_assistant", "bytes"),
					resource.TestCheckResourceAttrSet("openai_file.test_assistant", "created_at"),
					resource.TestCheckResourceAttr("openai_file.test_assistant", "status", "processed"),
				),
			},
		},
	})
}

func testAccResourceOpenAIFileConfig(filePath string) string {
	return fmt.Sprintf(`
resource "openai_file" "test" {
  file    = "%s"
  purpose = "fine-tune"
}
`, filePath)
}

func testAccResourceFileConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_file" "test" {
  content = jsonencode([
    {
      "messages": [
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello!"},
        {"role": "assistant", "content": "Hi there! How can I help you today?"}
      ]
    }
  ])
  filename = "training_data.jsonl"
  purpose  = "fine-tune"
}
`)
}

func testAccResourceFileConfig_assistantFile() string {
	return fmt.Sprintf(`
resource "openai_file" "test_assistant" {
  content  = "Here is some example content for testing purposes."
  filename = "assistant_knowledge.txt"
  purpose  = "assistants"
}
`)
}
