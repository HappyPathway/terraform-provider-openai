---
page_title: "openai_content_generator Resource - terraform-provider-openai"
description: |-
  Creates dynamic content with customizable output formats using OpenAI models.
---

# openai_content_generator (Resource)

Creates dynamic content with OpenAI models using chat completion APIs. Supports structured outputs through JSON schema definitions.

~> **Note:** This resource creates content each time it is created and cannot be updated. Changes to any of the arguments will result in new content being generated.

## Example Usage

### Basic Usage

```terraform
resource "openai_content_generator" "example" {
  model = "gpt-4"
  messages {
    role    = "user"
    content = "Write a haiku about coding"
  }
  temperature = 0.7
}

output "generated_content" {
  value = openai_content_generator.example.content
}
```

### JSON Output with Schema

```terraform
resource "openai_content_generator" "structured" {
  model = "gpt-4"
  
  messages {
    role    = "system"
    content = "You are a helpful assistant that provides structured data about movies."
  }
  
  messages {
    role    = "user"
    content = "Provide information about The Matrix (1999)"
  }

  response_format {
    type = "json_object"
    schema = jsonencode({
      type = "object"
      properties = {
        title = {
          type = "string"
          description = "The movie title"
        }
        year = {
          type = "integer"
          description = "Release year"
        }
        directors = {
          type = "array"
          items = {
            type = "string"
          }
          description = "List of directors"
        }
        genres = {
          type = "array"
          items = {
            type = "string"
          }
          description = "List of genres"
        }
      }
      required = ["title", "year", "directors", "genres"]
    })
  }

  temperature = 0.5
}

output "movie_info" {
  value = jsondecode(openai_content_generator.structured.content)
}
```

## Argument Reference

* `model` - (Required) The ID of the model to use (e.g., "gpt-4", "gpt-3.5-turbo").
* `messages` - (Required) A list of messages that form the conversation history. Each message has the following parameters:
  * `role` - (Required) The role of the message author. Must be one of: "system", "user", "assistant", or "function".
  * `content` - (Required) The content of the message.
* `response_format` - (Optional) Configuration block for specifying the output format. Maximum of 1 block.
  * `type` - (Required) The format to return the response in. Must be one of: "text" or "json_object".
  * `schema` - (Optional) JSON schema that the response should conform to when type is "json_object". Must be a valid JSON schema.
* `temperature` - (Optional) Sampling temperature between 0 and 2. Higher values like 0.8 make the output more random, while lower values like 0.2 make it more focused and deterministic. Defaults to 1.0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content` - The generated content from the model.
* `raw_response` - The complete raw response from the API in JSON format.
* `usage` - A map containing token usage statistics with the following keys:
  * `completion_tokens` - The number of tokens in the generated completion.
  * `prompt_tokens` - The number of tokens in the prompt.
  * `total_tokens` - The total number of tokens used in the request.