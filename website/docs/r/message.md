---
layout: "openai"
page_title: "OpenAI: openai_message"
description: |-
  Create and manage messages in threads.
---

# openai_message

This resource allows you to create, read, update, and delete messages within OpenAI threads. Messages can contain text content, image files, or image URLs, and can be associated with various tools like code interpreter or file search.

## Example Usage

### Basic Text Message

```hcl
resource "openai_thread" "example" {
  metadata = {
    purpose = "example"
  }
}

resource "openai_message" "text_example" {
  thread_id = openai_thread.example.id
  role      = "user"
  
  content {
    type = "text"
    text = "What is machine learning?"
  }

  metadata = {
    message_type = "question"
  }
}
```

### Message with Image File

```hcl
resource "openai_message" "image_example" {
  thread_id = openai_thread.example.id
  role      = "user"

  content {
    type = "image_file"
    image_file {
      file_id = "your-file-id"
      detail  = "auto"
    }
  }

  attachments {
    file_id = "your-file-id"
    tool    = "code_interpreter"
  }
}
```

### Message with Image URL

```hcl
resource "openai_message" "image_url_example" {
  thread_id = openai_thread.example.id
  role      = "user"

  content {
    type = "image_url"
    image_url {
      url    = "https://example.com/image.jpg"
      detail = "low"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `thread_id` - (Required) The ID of the thread to create a message in.
* `role` - (Required) The role of the entity creating the message. Must be either 'user' or 'assistant'.
* `content` - (Required) The content of the message. Can contain text, image files, or image URLs.
* `metadata` - (Optional) Set of key-value pairs that can be attached to the message.
* `attachments` - (Optional) List of files to attach to the message.

### Content Arguments

The `content` block supports:

* `type` - (Required) The type of content. Must be one of: text, image_file, or image_url.
* `text` - (Optional) The text content when type is "text".
* `image_file` - (Optional) Configuration block for image file content.
* `image_url` - (Optional) Configuration block for image URL content.

### Image File Arguments

The `image_file` block supports:

* `file_id` - (Required) The ID of the file to use.
* `detail` - (Optional) Level of detail for the image. Must be one of: auto, low, or high. Defaults to "auto".

### Image URL Arguments

The `image_url` block supports:

* `url` - (Required) The URL of the image.
* `detail` - (Optional) Level of detail for the image. Must be one of: auto, low, or high. Defaults to "auto".

### Attachments Arguments

The `attachments` block supports:

* `file_id` - (Required) The ID of the file to attach.
* `tool` - (Required) The tool to add this file to. Must be one of: code_interpreter or file_search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the message.
* `assistant_id` - ID of the assistant that authored this message, if applicable.
* `created_at` - Unix timestamp for when the message was created.
* `completed_at` - Unix timestamp for when the message was completed.
* `incomplete_at` - Unix timestamp for when the message was marked as incomplete.
* `incomplete_details` - Details about why the message is incomplete, if applicable.
* `run_id` - ID of the run associated with message creation.
* `status` - Status of the message: in_progress, incomplete, or completed.

## Import

Messages can be imported using a combination of the thread ID and message ID, separated by a slash:

```sh
$ terraform import openai_message.example thread_abc123/msg_xyz789
```