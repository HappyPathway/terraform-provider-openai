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
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + testAccResourceOpenAIFineTuningJobConfig(testDataPath),
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

func TestAccResourceFineTuningJob_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFineTuningJobConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_fine_tuning_job.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet("openai_fine_tuning_job.test", "training_file"),
					resource.TestCheckResourceAttrSet("openai_fine_tuning_job.test", "status"),
					resource.TestCheckResourceAttrSet("openai_fine_tuning_job.test", "fine_tuned_model"),
				),
			},
		},
	})
}

func TestAccResourceFineTuningJob_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFineTuningJobConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_fine_tuning_job.test_full", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttrSet("openai_fine_tuning_job.test_full", "training_file"),
					resource.TestCheckResourceAttrSet("openai_fine_tuning_job.test_full", "validation_file"),
					resource.TestCheckResourceAttr("openai_fine_tuning_job.test_full", "hyperparameters.n_epochs", "3"),
					resource.TestCheckResourceAttr("openai_fine_tuning_job.test_full", "suffix", "custom-model"),
					resource.TestCheckResourceAttrSet("openai_fine_tuning_job.test_full", "status"),
					resource.TestCheckResourceAttrSet("openai_fine_tuning_job.test_full", "fine_tuned_model"),
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

func testAccResourceFineTuningJobConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_file" "training" {
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

resource "openai_fine_tuning_job" "test" {
  model         = "gpt-3.5-turbo"
  training_file = openai_file.training.id
}
`)
}

func testAccResourceFineTuningJobConfig_full() string {
	return fmt.Sprintf(`
resource "openai_file" "training" {
  content = jsonencode([
    {
      "messages": [
        {"role": "system", "content": "You are a specialized assistant."},
        {"role": "user", "content": "Hello!"},
        {"role": "assistant", "content": "Greetings! I'm here to assist you with specialized tasks."}
      ]
    }
  ])
  filename = "training_data.jsonl"
  purpose  = "fine-tune"
}

resource "openai_file" "validation" {
  content = jsonencode([
    {
      "messages": [
        {"role": "system", "content": "You are a specialized assistant."},
        {"role": "user", "content": "Hi there!"},
        {"role": "assistant", "content": "Hello! I'm ready to help with specialized tasks."}
      ]
    }
  ])
  filename = "validation_data.jsonl"
  purpose  = "fine-tune"
}

resource "openai_fine_tuning_job" "test_full" {
  model           = "gpt-3.5-turbo"
  training_file   = openai_file.training.id
  validation_file = openai_file.validation.id
  
  hyperparameters {
    n_epochs = 3
  }
  
  suffix = "custom-model"
}
`)
}
