---
layout: "openai"
page_title: "OpenAI: openai_message"
description: |-
  Get information about a specific message in a thread.
---

# openai_message

Use this data source to get information about a specific message in a thread. This is useful when you want to reference an existing message's details or check its status.

## Example Usage

```hcl
data "openai_message" "example" {
  thread_id  = "thread_abc123"
  message_id = "msg_xyz789"
}

output "message_content" {
  value = data.openai_message.example.content
}
```

## Argument Reference

* `thread_id` - (Required) The ID of the thread containing the message.
* `message_id` - (Required) The ID of the message to retrieve.

## Attributes Reference

* `role` - The role of the entity that created the message (user or assistant).
* `content` - The content of the message, which can include:
  * `type` - The type of content (text, image_file, or image_url).
  * `text` - The text content when type is "text".
  * `image_file` - Configuration for image file content:
    * `file_id` - The ID of the file used.
    * `detail` - The detail level of the image (auto, low, or high).
  * `image_url` - Configuration for image URL content:
    * `url` - The URL of the image.
    * `detail` - The detail level of the image (auto, low, or high).
* `assistant_id` - ID of the assistant that authored this message.
* `run_id` - ID of the run associated with message creation.
* `created_at` - Unix timestamp for when the message was created.
* `completed_at` - Unix timestamp for when the message was completed.
* `incomplete_at` - Unix timestamp for when the message was marked as incomplete.
* `incomplete_details` - Details about why the message is incomplete:
  * `reason` - The reason for the message being incomplete.
* `metadata` - Set of key-value pairs attached to the message.
* `status` - Status of the message (in_progress, incomplete, or completed).
* `attachments` - List of files attached to the message:
  * `file_id` - The ID of the attached file.
  * `tool` - The tool the file is attached to (code_interpreter or file_search).