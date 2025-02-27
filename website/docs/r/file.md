---
page_title: "openai_file Resource - OpenAI Provider"
subcategory: ""
description: |-
  Manages an OpenAI File resource for use with assistants, fine-tuning, and other features.
---

# Resource: openai_file

This resource allows you to manage files in your OpenAI account. Files can be used for various purposes including fine-tuning models, providing context to assistants, and storing assistant outputs.

## Example Usage

```terraform
resource "openai_file" "example" {
  file    = "${path.module}/data.jsonl"
  purpose = "fine-tune"
}
```

## Argument Reference

* `file` - (Required) Path to the file to upload. The file must be accessible from where Terraform is run.
* `purpose` - (Required) The intended purpose of the file. Possible values are: 'fine-tune', 'fine-tune-results', 'assistants', or 'assistants_output'.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `bytes` - Size of the file in bytes.
* `created_at` - The Unix timestamp (in seconds) for when the file was created.
* `filename` - Name of the file.

## Import

Files can be imported using the file ID, e.g.,

```shell
terraform import openai_file.example file-abc123
```