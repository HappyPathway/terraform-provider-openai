package openai

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIFineTuningJob_basic(t *testing.T) {
	testDataDir := "testdata"
	testDataPath := filepath.Join(testDataDir, "test.jsonl")

	// Create test training file
	testData := `{"messages": [{"role": "system", "content": "You analyze data."}, {"role": "user", "content": "Data: 123"}, {"role": "assistant", "content": "Analysis: The value is 123."}]}`
	if err := os.WriteFile(testDataPath, []byte(testData+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDataPath)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIFineTuningJobConfig(testDataPath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_fine_tuning_job.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet(
						"openai_fine_tuning_job.test", "status"),
					resource.TestCheckResourceAttrSet(
						"openai_fine_tuning_job.test", "created_at"),
				),
			},
		},
	})
}

func testAccResourceOpenAIFineTuningJobConfig(trainingFilePath string) string {
	return fmt.Sprintf(`
resource "openai_file" "training" {
  file    = "%s"
  purpose = "fine-tune"
}

resource "openai_fine_tuning_job" "test" {
  model          = "gpt-3.5-turbo"
  training_file  = openai_file.training.id
  n_epochs       = 3
  suffix         = "test-model"
}
`, trainingFilePath)
}
