---
page_title: "openai_thread Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Create and manage OpenAI Threads for conversations with Assistants.
---

# openai_thread

Creates and manages OpenAI Threads, which are conversation contexts that collect and organize messages between users and assistants.

## Example Usage

```terraform
resource "openai_thread" "support_conversation" {
  metadata = {
    user_id      = "user_123456"
    conversation_type = "technical_support"
    priority     = "high"
  }

  # Optional: Initialize the thread with messages
  messages = [
    {
      role     = "user"
      content  = "I'm having trouble with my account login."
      file_ids = []
      metadata = {
        importance = "high"
      }
    }
  ]
}

output "thread_id" {
  value = openai_thread.support_conversation.id
}
```

## Argument Reference

- `metadata` - (Optional) A map of key-value pairs that can be used to store additional information about the thread.
- `messages` - (Optional) A list of initial messages to add to the thread. Each message contains:
  - `role` - (Required) The role of the message author. Can be "user" or "assistant".
  - `content` - (Required) The content of the message.
  - `file_ids` - (Optional) A list of file IDs to attach to the message.
  - `metadata` - (Optional) A map of key-value pairs with additional information about the message.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The OpenAI-assigned ID for this thread.
- `created_at` - The timestamp when the thread was created.
- `object` - The object type, always "thread".

## Import

Threads can be imported using the OpenAI thread ID:

```
$ terraform import openai_thread.example thread_abc123
```
