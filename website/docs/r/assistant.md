---
page_title: "openai_assistant Resource - OpenAI Provider"
subcategory: ""
description: |-
  Manages an OpenAI Assistant that can use multiple tools and models to help with a variety of tasks.
---

# Resource: openai_assistant

This resource allows you to create and manage OpenAI Assistants. Assistants can be configured with specific models, instructions, and tools to help with various tasks.

## Example Usage

```terraform
resource "openai_assistant" "example" {
  name         = "Research Assistant"
  description  = "A research assistant that helps with data analysis"
  model        = "gpt-4-1106-preview"
  instructions = "You are a helpful research assistant. Use the provided knowledge base to answer questions accurately."

  tools {
    type = "retrieval"
  }

  tools {
    type = "code_interpreter"
  }

  tools {
    type        = "function"
    name        = "search_papers"
    description = "Search for relevant research papers"
    parameters  = jsonencode({
      type = "object"
      properties = {
        query = {
          type = "string"
          description = "The search query"
        }
        year = {
          type = "integer"
          description = "Filter by publication year"
        }
      }
      required = ["query"]
    })
  }

  file_ids = ["file-abc123"]
}
```

## Argument Reference

* `name` - (Required) The name of the assistant.
* `model` - (Required) ID of the model to use. For example, `gpt-4` or `gpt-4-1106-preview`.
* `description` - (Optional) A description of the assistant.
* `instructions` - (Optional) The system instructions that set the behavior and context for the assistant.
* `tools` - (Optional) A list of tools enabled on the assistant. Each tool block supports:
  * `type` - (Required) The type of tool. Can be one of: `code_interpreter`, `retrieval`, or `function`.
  * `name` - (Optional) The name of the function. Required when type is `function`.
  * `description` - (Optional) A description of what the function does. Required when type is `function`.
  * `parameters` - (Optional) The parameters the function accepts (in JSON schema format). Required when type is `function`.
* `file_ids` - (Optional) A list of file IDs attached to this assistant.
* `metadata` - (Optional) Key-value pairs that can be attached to the assistant.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the assistant.
* `created_at` - The Unix timestamp (in seconds) for when the assistant was created.

## Import

Assistants can be imported using the assistant ID, e.g.,

```shell
terraform import openai_assistant.example asst_abc123
```