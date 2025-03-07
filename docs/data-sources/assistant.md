---
page_title: "openai_assistant Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  Retrieve information about OpenAI Assistants.
---

# openai_assistant Data Source

Retrieves information about existing OpenAI Assistants. This data source allows you to look up and use assistants in your Terraform configurations without managing them directly.

## Example Usage

```terraform
# Get a specific assistant by ID
data "openai_assistant" "customer_support" {
  assistant_id = "asst_abc123"
}

# Create a thread that uses an existing assistant
resource "openai_thread" "new_conversation" {
  # Thread configuration
}

output "assistant_details" {
  value = {
    name  = data.openai_assistant.customer_support.name
    model = data.openai_assistant.customer_support.model
  }
}
```

## Argument Reference

- `assistant_id` - (Required) The ID of the assistant to retrieve.

## Attribute Reference

- `id` - The ID of the assistant.
- `name` - The name of the assistant, if set.
- `description` - The description of the assistant, if set.
- `model` - The model used by the assistant.
- `instructions` - The instructions that set the behavior and capabilities of the assistant.
- `tools` - A list of tools enabled for the assistant. Each tool contains:
  - `type` - The tool type, such as "file_search", "code_interpreter", or "function".
  - `function` - For function tools, the function definition as a JSON string.
- `tool_resources` - Resources made available to the assistant's tools:
  - `code_interpreter` - Configuration for the code interpreter tool:
    - `file_ids` - List of file IDs that the code interpreter can use.
  - `file_search` - Configuration for the file search tool:
    - `vector_store_ids` - List of vector store IDs for the file search capability.
- `metadata` - Additional key-value pairs associated with the assistant.
- `created_at` - The timestamp when the assistant was created.
- `object` - The object type, always "assistant".
