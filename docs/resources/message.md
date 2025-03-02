---
page_title: "openai_message Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Create and manage OpenAI Messages within Threads.
---

# openai_message

Creates and manages individual messages within OpenAI Threads. Messages are the building blocks of conversations between users and assistants.

## Example Usage

```terraform
resource "openai_thread" "conversation" {
  # Thread configuration
}

resource "openai_message" "user_query" {
  thread_id = openai_thread.conversation.id
  role      = "user"
  content   = "Can you help me understand how to implement authentication in my Node.js application?"

  metadata = {
    source = "customer_portal"
    user_timezone = "America/New_York"
  }
}

# You can add assistant messages as well, though typically they are generated through runs
resource "openai_message" "assistant_response" {
  thread_id = openai_thread.conversation.id
  role      = "assistant"
  content   = "I'd be happy to help you implement authentication in your Node.js application. There are several approaches you can take..."

  # Optional: attach files to the message
  file_ids = []

  metadata = {
    response_type = "detailed"
  }
}
```

## Argument Reference

- `thread_id` - (Required) The ID of the thread to add the message to.
- `role` - (Required) The role of the message author. Can be "user" or "assistant".
- `content` - (Required) The content of the message.
- `file_ids` - (Optional) A list of file IDs to attach to the message. These files must already be uploaded to OpenAI with purpose "assistants".
- `metadata` - (Optional) A map of key-value pairs that can be used to store additional information about the message.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The OpenAI-assigned ID for this message.
- `created_at` - The timestamp when the message was created.
- `object` - The object type, always "thread.message".
- `assistant_id` - If applicable, the ID of the assistant that created the message.

## Import

Messages can be imported using the format `thread_id:message_id`:

```
$ terraform import openai_message.example thread_abc123:msg_xyz789
```
