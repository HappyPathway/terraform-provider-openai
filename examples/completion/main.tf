terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {}

resource "openai_completion" "story" {
  model       = "gpt-3.5-turbo-instruct"
  prompt      = "Write a short story about AI and Terraform"
  max_tokens  = 150
  temperature = 0.7
}

output "story" {
  value = openai_completion.story.choices[0].text
}