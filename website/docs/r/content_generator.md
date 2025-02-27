---
page_title: "openai_content_generator Resource - OpenAI Provider"
subcategory: ""
description: |-
  Manages content generation using OpenAI models.
---

# Resource: openai_content_generator

This resource allows you to generate content using OpenAI models and store the results in Terraform state. This is useful for generating text content that needs to be consistent across infrastructure deployments.

## Example Usage

```terraform
resource "openai_content_generator" "description" {
  model = "gpt-4"
  prompt = "Write a brief description of a cloud infrastructure monitoring system:"
  
  parameters {
    temperature = 0.7
    max_tokens  = 150
  }
}

output "generated_description" {
  value = openai_content_generator.description.content
}
```

## Argument Reference

* `model` - (Required) The ID of the model to use for content generation.
* `prompt` - (Required) The prompt to send to the model.
* `parameters` - (Optional) Configuration for the content generation. Supports:
  * `temperature` - (Optional) Sampling temperature between 0 and 2. Higher values make output more random, lower values make it more deterministic.
  * `max_tokens` - (Optional) Maximum number of tokens to generate.
  * `top_p` - (Optional) Alternative to temperature. Nucleus sampling parameter between 0 and 1.
  * `frequency_penalty` - (Optional) Number between -2.0 and 2.0. Positive values penalize new tokens based on their frequency.
  * `presence_penalty` - (Optional) Number between -2.0 and 2.0. Positive values penalize tokens based on whether they've appeared in the text so far.

## Attributes Reference

* `content` - The generated content from the model.
* `finish_reason` - The reason why the model stopped generating text.
* `usage` - Information about token usage:
  * `prompt_tokens` - Number of tokens in the prompt.
  * `completion_tokens` - Number of tokens in the generated completion.
  * `total_tokens` - Total number of tokens used.