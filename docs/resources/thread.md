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
# Basic thread with metadata
resource "openai_thread" "support_conversation" {
  metadata = {
    user_id            = "user_123456"
    conversation_type  = "technical_support"
    priority          = "high"
  }
}

# Thread with initial messages
resource "openai_thread" "onboarding" {
  metadata = {
    user_id = "new_user_789"
  }

  messages = [
    {
      role    = "user"
      content = "Hi! I'm new here and would like to learn about the platform."
      metadata = {
        topic = "onboarding"
      }
    }
  ]
}

# Thread with tool resources
resource "openai_thread" "data_analysis" {
  tools = ["code_interpreter", "file_search"]

  tool_resources = {
    code_interpreter = {
      file_ids = [openai_file.dataset.id]
    }
    file_search = {
      vector_store_ids = [openai_vector_store.documentation.id]
    }
  }

  metadata = {
    project = "quarterly_analysis"
  }
}
```

## Argument Reference

- `metadata` - (Optional) A map of key-value pairs that can be used to store additional information about the thread.
- `tools` - (Optional) A list of tools enabled for this thread. Valid values are:
  - `code_interpreter` - For executing code and working with files.
  - `file_search` - For semantic search in uploaded files.
- `tool_resources` - (Optional) Resources made available to the thread's tools. Contains nested blocks:
  - `code_interpreter` - (Optional) Configuration for the code interpreter tool:
    - `file_ids` - (Optional) File IDs that the code interpreter can use.
  - `file_search` - (Optional) Configuration for the file search tool:
    - `vector_store_ids` - (Optional) Vector store IDs for the file search capability.
- `messages` - (Optional) A list of initial messages to add to the thread. Each message contains:
  - `role` - (Required) The role of the entity creating the message. Can be "user" or "assistant".
  - `content` - (Required) The content of the message.
  - `metadata` - (Optional) A map of key-value pairs with additional information about the message.
  - `file_ids` - (Optional, Deprecated) A list of file IDs to attach to the message. Deprecated in v2, use message attachments instead.

~> **Note** In v2 of the API, files are managed through tool_resources and message attachments rather than direct file_ids.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The OpenAI-assigned ID for this thread.
- `created_at` - The Unix timestamp (in seconds) for when the thread was created.
- `object` - The object type, always "thread".

## Import

Threads can be imported using the OpenAI thread ID:

```shell
$ terraform import openai_thread.example thread_abc123
```
