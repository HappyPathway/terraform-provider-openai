---
page_title: "openai_fine_tune Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Create and manage fine-tuning jobs for OpenAI models.
---

# openai_fine_tune

Creates and manages fine-tuning jobs for OpenAI models. Fine-tuning allows you to customize a model with your own training data to achieve better results for specific use cases.

## Example Usage

```terraform
resource "openai_file" "training_data" {
  file_path = "./data/training_data.jsonl"
  purpose   = "fine-tune"
}

resource "openai_fine_tune" "custom_model" {
  training_file_id = openai_file.training_data.id
  model            = "gpt-3.5-turbo"
  suffix           = "customer-support-specialist"
  epochs           = 4
}

output "fine_tuned_model" {
  value = openai_fine_tune.custom_model.fine_tuned_model
}
```

## Argument Reference

- `training_file_id` - (Required) The ID of the file to use for training. This file must be uploaded with the purpose "fine-tune".
- `model` - (Required) The name of the base model to fine-tune. You can use "gpt-3.5-turbo" or a previously fine-tuned model.
- `validation_file_id` - (Optional) The ID of the file to use for validation.
- `suffix` - (Optional) A string that will be added to the fine-tuned model name to help you identify it.
- `epochs` - (Optional) Number of training epochs. Default is determined by the model.
- `batch_size` - (Optional) **Deprecated**: This parameter is no longer directly supported by OpenAI's new fine-tuning API and will be ignored.
- `learning_rate_multiplier` - (Optional) **Deprecated**: This parameter is no longer directly supported by OpenAI's new fine-tuning API and will be ignored.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The OpenAI-assigned ID for this fine-tuning job.
- `object_id` - The ID of the fine-tuning job.
- `status` - The status of the fine-tuning job. Can be "pending", "running", "succeeded", "failed", or "cancelled".
- `fine_tuned_model` - The name of the fine-tuned model once training has completed successfully.
- `created_at` - The timestamp when the fine-tuning job was created.
- `organization_id` - The organization that owns the fine-tuning job.
- `result_files` - IDs of result files created by the fine-tuning job.

## Import

Fine-tuning jobs can be imported using the OpenAI fine-tune ID:

```
$ terraform import openai_fine_tune.example ft-abc123
```
