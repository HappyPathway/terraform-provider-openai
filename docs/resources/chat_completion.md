---
page_title: "openai_chat_completion Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Generate a chat completion using OpenAI's GPT models.
---

# openai_chat_completion

Generates a chat completion using OpenAI's GPT models. This resource allows you to interact with OpenAI's language models by providing conversation messages and receiving AI-generated responses.

## Example Usage

```terraform
resource "openai_chat_completion" "example" {
  model = "gpt-4"

  messages = [
    {
      role    = "system"
      content = "You are a helpful assistant that provides concise responses."
    },
    {
      role    = "user"
      content = "What is the capital of France?"
    }
  ]

  temperature = 0.7
  max_tokens  = 100
}

output "assistant_response" {
  value = openai_chat_completion.example.response_content[0]
}
```

## Argument Reference

- `model` - (Required) ID of the model to use for completion (e.g., 'gpt-4', 'gpt-3.5-turbo').
- `messages` - (Required) A list of messages comprising the conversation so far. Each message object contains:
  - `role` - (Required) The role of the message author. Can be 'system', 'user', or 'assistant'.
  - `content` - (Required) The content of the message.
- `temperature` - (Optional) Sampling temperature between 0 and 2. Higher values like 0.8 make output more random, while lower values like 0.2 make it more focused and deterministic.
- `top_p` - (Optional) An alternative to sampling with temperature, called nucleus sampling. Set this between 0 and 1.
- `n` - (Optional) How many completion choices to generate. Default is 1.
- `max_tokens` - (Optional) Maximum number of tokens to generate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - Unique identifier for this resource.
- `response_content` - The completion(s) generated by the model. For n=1, this will be a list with a single element.
- `response_role` - The role of the returned message. Typically 'assistant'.

## Import

This resource does not support import as it is stateless and generates new completions on each apply.
