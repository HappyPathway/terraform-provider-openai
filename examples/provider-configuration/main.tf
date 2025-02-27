terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

# Configure the OpenAI Provider
provider "openai" {
  api_key = var.openai_api_key # Or use OPENAI_API_KEY environment variable
}

# Optional: configure organization ID if needed
variable "openai_api_key" {
  description = "OpenAI API Key"
  type        = string
  sensitive   = true
}