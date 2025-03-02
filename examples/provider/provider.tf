terraform {
  required_providers {
    openai = {
      source  = "darnold/openai"
      version = "0.1.0"
    }
  }
}

# Configure the OpenAI Provider
provider "openai" {
  # api_key = "your-api-key" # or use OPENAI_API_KEY env var
  # organization = "your-org-id" # or use OPENAI_ORGANIZATION env var
  # enable_debug_logging = true
}

# Use the OpenAI model data source 
data "openai_model" "gpt4" {
  model_id = "gpt-4"
}

# Create a chat completion resource
resource "openai_chat_completion" "example" {
  model = data.openai_model.gpt4.model_id

  messages = [
    {
      role    = "system"
      content = "You are a helpful assistant."
    },
    {
      role    = "user"
      content = "Write a haiku about Terraform."
    }
  ]

  temperature = 0.7
  max_tokens  = 150
}

# Output the generated text
output "chat_response" {
  value = openai_chat_completion.example.response_content[0]
}
