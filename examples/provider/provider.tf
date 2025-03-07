terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

# Configure the OpenAI Provider
provider "openai" {
  # Configure OpenAI API key - this can also be set via OPENAI_API_KEY environment variable
  # api_key = "your-api-key"

  # Optional: Configure organization ID if you're part of multiple organizations
  # organization_id = "your-org-id"
}

# Use the OpenAI model data source 
data "openai_model" "gpt4" {
  model_id = "gpt-4"
}

# Create a chat completion resource
resource "openai_chat_completion" "example" {
  model = data.openai_model.gpt4.model_id

  messages {
    role    = "system"
    content = "You are a helpful assistant."
  }
  messages {
    role    = "user"
    content = "Write a haiku about Terraform."
  }

  temperature = 0.7
  max_tokens  = 150
}

# Output the generated text
output "chat_response" {
  value = openai_chat_completion.example.response_content[0]
}
