terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

resource "openai_file" "training_data" {
  content  = filebase64("${path.module}/training_data.jsonl")
  purpose  = "fine-tune"
}

resource "openai_fine_tuning_job" "custom_model" {
  model           = "gpt-3.5-turbo"
  training_file   = openai_file.training_data.id
  validation_file = null # Optional: Add a validation file to evaluate the model

  # Optional configuration
  hyperparameters {
    n_epochs = 3
  }

  suffix = "custom-support-bot" # This will be added to your model name
}

# Optional: Data source to check available models
data "openai_model" "fine_tuned" {
  id = openai_fine_tuning_job.custom_model.fine_tuned_model
}

output "fine_tuned_model" {
  value = data.openai_model.fine_tuned.id
}

output "training_status" {
  value = openai_fine_tuning_job.custom_model.status
}