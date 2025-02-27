---
page_title: "openai_assistant Resource - terraform-provider-openai"
description: |-
  Manages OpenAI assistants.
---

# openai_assistant (Resource)

Manages OpenAI assistants. Assistants can be customized with specific instructions, tools, and knowledge to help with a variety of tasks.

## Example Usage

```terraform
resource "openai_assistant" "example" {
  name         = "Customer Support Assistant"
  description  = "An assistant that helps with customer support queries"
  model        = "gpt-4-1106-preview"
  instructions = "You are a helpful customer support assistant. Be concise and friendly in your responses."
  tools = [
    {
      type = "code_interpreter"
    }
  ]
  metadata = {
    department = "customer_support"
    team       = "tier_1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the assistant.
* `model` - (Required) ID of the model to use. You can use GPT-4, GPT-3.5, or other models enabled for your organization.
* `description` - (Optional) A description of the assistant.
* `instructions` - (Optional) Instructions that the assistant should follow when interacting.
* `tools` - (Optional) A list of tools enabled for the assistant. Each tool has a `type` field which must be one of "code_interpreter", "retrieval", or "function".
* `file_ids` - (Optional) A list of file IDs attached to the assistant. These files are used by the retrieval tool.
* `metadata` - (Optional) A map of key-value pairs that can be used to store additional information about the assistant.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the assistant.
* `object` - The object type, which is always "assistant".
* `created_at` - The Unix timestamp (in seconds) for when the assistant was created.

## Import

Assistants can be imported using the assistant ID, e.g.:

```shell
$ terraform import openai_assistant.example asst_abc123
```