---
page_title: "openai_file Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Upload and manage files for use with OpenAI APIs.
---

# openai_file

Uploads and manages files for use with OpenAI APIs, such as fine-tuning data or files for retrieval in assistants.

## Example Usage

```terraform
resource "openai_file" "fine_tune_data" {
  file_path = "./data/fine_tune_data.jsonl"
  purpose   = "fine-tune"
}

resource "openai_file" "assistant_file" {
  file_path = "./data/knowledge_base.pdf"
  purpose   = "assistants"
}
```

## Argument Reference

- `file_path` - (Required) Path to the file to be uploaded.
- `purpose` - (Required) The intended purpose of the file. Allowed values are "fine-tune" or "assistants".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The OpenAI-assigned ID for this file.
- `bytes` - The size of the file in bytes.
- `created_at` - The timestamp when the file was created.
- `filename` - The name of the file.
- `status` - The status of the file. Can be "uploaded", "processed", or "error".
- `status_details` - Additional details about the file status, if available.

## Import

Files can be imported using the OpenAI file ID:

```
$ terraform import openai_file.example file-abc123
```
