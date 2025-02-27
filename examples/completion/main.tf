terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {}

resource "openai_content_generator" "story" {
  model       = "gpt-3.5-turbo"
  messages {
    role    = "user"
    content = "Write a short story about AI and Terraform"
  }
  temperature = 0.7
}

output "story" {
  value = openai_content_generator.story.content
}