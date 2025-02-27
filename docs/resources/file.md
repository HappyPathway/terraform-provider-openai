---
page_title: "openai_file Resource - terraform-provider-openai"
description: |-
  Manages files uploaded to OpenAI.
---

# openai_file (Resource)

Manages files uploaded to OpenAI. Files are used for training, fine-tuning, and other API operations.

## Example Usage

```terraform
resource "openai_file" "training_data" {
  content  = filebase64("${path.module}/training_data.jsonl")
  filename = "training_data.jsonl"
  purpose  = "fine-tune"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, Forces new resource) The base64-encoded content of the file. Use the `filebase64()` function to read a file from disk.
* `filename` - (Required, Forces new resource) The name of the file to be uploaded.
* `purpose` - (Required, Forces new resource) The intended purpose of the uploaded file. Valid values are "fine-tune", "assistants", or "fine-tune-results".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the file.
* `bytes` - The size of the file in bytes.
* `created_at` - The Unix timestamp (in seconds) for when the file was created.
* `object` - The object type, which is always "file".

## Import

Files can be imported using the file ID, e.g.:

```shell
$ terraform import openai_file.training_data file-abc123
```