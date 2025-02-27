terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {}

# Upload training data file
resource "openai_file" "training_data" {
  content  = filebase64("${path.module}/training_data.jsonl")
  filename = "training_data.jsonl"
  purpose  = "fine-tune"
}

# Create fine-tuning job
resource "openai_fine_tuning_job" "custom_model" {
  model          = "gpt-3.5-turbo"
  training_file  = openai_file.training_data.id
}

# Optional: Data source to check available models
data "openai_model" "fine_tuned" {
  model_id = openai_fine_tuning_job.custom_model
}

output "fine_tuned_model" {
  value = data.openai_model.fine_tuned.id
}

output "training_status" {
  value = openai_fine_tuning_job.custom_model.status
}