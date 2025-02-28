terraform {
  required_providers {
    openai = {
      source = "openai/openai"
    }
  }
  required_version = ">= 0.13"
}

provider "openai" {
  # Configuration loaded from environment variables:
  # OPENAI_API_KEY
  # OPENAI_ORGANIZATION_ID (optional)
}

# Upload a file for fine-tuning
resource "openai_file" "training_data" {
  file    = "${path.module}/training_data.jsonl"
  purpose = "fine-tune"
}

# Example of retrieving information about a fine-tuned model
# Note: The model must be created through the OpenAI API or CLI first
data "openai_fine_tuned_model" "custom_model" {
  model_id = "ft:gpt-3.5-turbo:my-org:custom_model:id123" # Replace with your model ID
}

output "model_info" {
  value = {
    id        = data.openai_fine_tuned_model.custom_model.id
    created   = data.openai_fine_tuned_model.custom_model.created
    owned_by  = data.openai_fine_tuned_model.custom_model.owned_by
  }
}