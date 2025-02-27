package openai

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIFineTuningJob_basic(t *testing.T) {
	// Skip if required variables aren't set
	projectPrompt := os.Getenv("TF_VAR_project_prompt")
	repoOrg := os.Getenv("TF_VAR_repo_org")
	projectName := os.Getenv("TF_VAR_project_name")
	apiKey := os.Getenv("OPENAI_API_KEY")

	if projectPrompt == "" || repoOrg == "" || projectName == "" {
		t.Skip("Required variables TF_VAR_project_prompt, TF_VAR_repo_org, or TF_VAR_project_name not set")
	}

	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	testDataPath := "testdata/test.jsonl"

	// Create test training file
	err := os.WriteFile(testDataPath, []byte(
		"{\"messages\": [{\"role\": \"system\", \"content\": \"You analyze company revenue.\"}, {\"role\": \"user\", \"content\": \"Company: Acme Corp, Revenue: $10M\"}, {\"role\": \"assistant\", \"content\": \"Revenue analysis: Acme Corp reported $10M in revenue.\"}]}\n"+
			"{\"messages\": [{\"role\": \"system\", \"content\": \"You analyze company revenue.\"}, {\"role\": \"user\", \"content\": \"Company: TechCo, Revenue: $5M\"}, {\"role\": \"assistant\", \"content\": \"Revenue analysis: TechCo reported $5M in revenue.\"}]}\n",
	), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDataPath)

	modelPattern := regexp.MustCompile("^gpt-3\\.5-turbo")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenAIFineTuningJobConfig(apiKey, testDataPath, projectPrompt, repoOrg, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"openai_fine_tuning_job.test", "model",
						modelPattern,
					),
					resource.TestCheckResourceAttrSet(
						"openai_fine_tuning_job.test", "status"),
					resource.TestCheckResourceAttrSet(
						"openai_fine_tuning_job.test", "created_at"),
				),
			},
		},
	})
}

func testAccResourceOpenAIFineTuningJobConfig(apiKey, trainingFilePath, projectPrompt, repoOrg, projectName string) string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = "%s"
}

resource "openai_file" "training" {
  file    = "%s"
  purpose = "fine-tune"
}

resource "openai_fine_tuning_job" "test" {
  model          = "gpt-3.5-turbo-0613"
  training_file  = openai_file.training.id
  n_epochs       = 3
  suffix         = "test-model"

  depends_on = [openai_file.training]
}

# Required variables for the test environment
variable "project_prompt" {
  default = "%s"
}

variable "repo_org" {
  default = "%s"
}

variable "project_name" {
  default = "%s"
}
`, apiKey, trainingFilePath, projectPrompt, repoOrg, projectName)
}
