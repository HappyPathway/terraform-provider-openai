package openai

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIFile_basic(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	testDataPath := "testdata/test.jsonl"

	// Create test data directory if it doesn't exist
	err := os.MkdirAll("testdata", 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create test file
	err = os.WriteFile(testDataPath, []byte("{\"prompt\": \"test\", \"completion\": \"test completion\"}\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDataPath)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIFileConfig(apiKey, testDataPath),
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

func testAccResourceOpenAIFileConfig(apiKey, filePath string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_file" "test" {
  file    = "%s"
  purpose = "fine-tune"
}
`, apiKey, filePath)
}
