terraform {
  required_providers {
    openai = {
      source  = "darnold/openai"
      version = "0.1.0"
    }
  }
}

provider "openai" {}

# First, upload the training data file
resource "openai_file" "training_data" {
  file_path = "${path.module}/data/training_data.jsonl"
  purpose   = "fine-tune"
}

# Optional: Upload validation data
resource "openai_file" "validation_data" {
  file_path = "${path.module}/data/validation_data.jsonl"
  purpose   = "fine-tune"
}

# Create a fine-tuning job
resource "openai_fine_tune" "custom_model" {
  training_file   = openai_file.training_data.id
  validation_file = openai_file.validation_data.id
  model           = "gpt-3.5-turbo"

  # Optional: add a suffix to the model name
  suffix = "customer-service-assistant"

  # Optional: configure hyperparameters
  hyperparameters = {
    n_epochs                 = 4
    batch_size               = null # Let OpenAI determine optimal batch size
    learning_rate_multiplier = null # Let OpenAI determine optimal learning rate
  }
}

output "fine_tuned_model" {
  value       = openai_fine_tune.custom_model.fine_tuned_model
  description = "The name of the fine-tuned model (available once training completes)"
}

output "fine_tune_status" {
  value       = openai_fine_tune.custom_model.status
  description = "Current status of the fine-tuning job"
}

output "fine_tune_events" {
  value       = openai_fine_tune.custom_model.events
  description = "Events from the fine-tuning process"
}

# Note: You need to create the following files before running this example:
# - data/training_data.jsonl - JSONL file with training data in the chat format
# - data/validation_data.jsonl - JSONL file with validation data in the chat format
