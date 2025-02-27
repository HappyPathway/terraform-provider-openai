terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

# Configure the OpenAI Provider
provider "openai" {
  api_key      = var.openai_api_key # Set this using environment variable TF_VAR_openai_api_key
  organization = var.organization_id # Optional: Set this using TF_VAR_organization_id if you're using an org
  retry_max    = 3                  # Optional: Number of retries for failed API calls
  retry_delay  = 5                  # Optional: Delay between retries in seconds
  timeout      = 30                 # Optional: Timeout for API calls in seconds
}