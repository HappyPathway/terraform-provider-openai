---
page_title: "openai_assistant Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Create and manage OpenAI Assistants for building integrated AI experiences.
---

# openai_assistant

Creates and manages OpenAI Assistants, which are specialized AI entities that use specific instructions, capabilities, and knowledge to help users with tasks.

## Example Usage

```terraform
resource "openai_file" "knowledge_base" {
  file_path = "./data/knowledge_base.pdf"
  purpose   = "assistants"
}

resource "openai_assistant" "customer_support" {
  name         = "Customer Support Agent"
  description  = "An assistant that helps with customer inquiries about our products"
  model        = "gpt-4-1106-preview"
  instructions = "You are a customer support assistant for a tech company. Answer questions helpfully and accurately based on the provided knowledge base."

  tools = [
    {
      type = "retrieval"
    },
    {
      type = "code_interpreter"
    }
  ]

  file_ids = [openai_file.knowledge_base.id]

  metadata = {
    department = "customer_support"
    team       = "technical"
    version    = "1.0"
  }
}
```

## Argument Reference

- `name` - (Optional) The name of the assistant.
- `description` - (Optional) The description of the assistant.
- `model` - (Required) The model to use for the assistant, e.g., "gpt-4-1106-preview" or "gpt-3.5-turbo".
- `instructions` - (Optional) Instructions for how the assistant should behave and respond.
- `tools` - (Optional) A list of tools enabled for the assistant. Each tool has a `type` which can be:
  - `"retrieval"` - Allows the assistant to access its knowledge base.
  - `"code_interpreter"` - Enables the assistant to write and execute code.
  - `"function"` - Lets the assistant call external functions.
- `file_ids` - (Optional) A list of file IDs that have been uploaded to OpenAI with the purpose "assistants".
- `metadata` - (Optional) A map of key-value pairs that can be used to organize and categorize the assistant.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The OpenAI-assigned ID for this assistant.
- `created_at` - The timestamp when the assistant was created.
- `object` - The object type, always "assistant".

## Import

Assistants can be imported using the OpenAI assistant ID:

```
$ terraform import openai_assistant.example asst_abc123
```
