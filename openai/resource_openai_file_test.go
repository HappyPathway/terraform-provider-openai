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
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIFileConfig(testDataPath),
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

func testAccResourceOpenAIFileConfig(filePath string) string {
	return fmt.Sprintf(`
resource "openai_file" "test" {
  file    = "%s"
  purpose = "fine-tune"
}
`, filePath)
}
