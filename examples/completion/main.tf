terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {
  # API key can be provided via OPENAI_API_KEY environment variable
}

# Basic text completion
resource "openai_content_generator" "basic" {
  model = "gpt-3.5-turbo"
  temperature = 0.7

  messages {
    role    = "system"
    content = "You are a helpful programming assistant."
  }

  messages {
    role    = "user"
    content = "Write a brief explanation of recursion in programming."
  }
}

# Structured JSON output
resource "openai_content_generator" "json" {
  model = "gpt-4"
  temperature = 0.7

  messages {
    role    = "user"
    content = "Generate a recipe for chocolate chip cookies"
  }

  response_format {
    type = "json_object"
    schema = jsonencode({
      type = "object"
      properties = {
        name = {
          type = "string"
        }
        prep_time_minutes = {
          type = "integer"
        }
        ingredients = {
          type = "array"
          items = {
            type = "object"
            properties = {
              item = {
                type = "string"
              }
              amount = {
                type = "string"
              }
            }
          }
        }
        instructions = {
          type = "array"
          items = {
            type = "string"
          }
        }
      }
      required = ["name", "prep_time_minutes", "ingredients", "instructions"]
    })
  }
}

output "explanation" {
  value = openai_content_generator.basic.content
}

output "recipe" {
  value = jsondecode(openai_content_generator.json.content)
}