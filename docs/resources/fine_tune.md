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
  training_file     = openai_file.training_data.id
  model             = "gpt-3.5-turbo"
  suffix            = "customer-support-specialist"
  hyperparameters = {
    n_epochs = 4
  }
}

output "fine_tuned_model" {
  value = openai_fine_tune.custom_model.fine_tuned_model
}
```

## Argument Reference

- `training_file` - (Required) The ID of the file to use for training. This file must be uploaded with the purpose "fine-tune".
- `model` - (Required) The name of the base model to fine-tune. You can use "gpt-3.5-turbo" or a previously fine-tuned model.
- `validation_file` - (Optional) The ID of the file to use for validation.
- `suffix` - (Optional) A string that will be added to the fine-tuned model name to help you identify it.
- `hyperparameters` - (Optional) Hyperparameters used for fine-tuning:
  - `n_epochs` - (Optional) Number of training epochs. Default is determined by the model.
  - `batch_size` - (Optional) Batch size to use for training.
  - `learning_rate_multiplier` - (Optional) Learning rate multiplier to use for training.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The OpenAI-assigned ID for this fine-tuning job.
- `status` - The status of the fine-tuning job. Can be "pending", "running", "succeeded", "failed", or "cancelled".
- `fine_tuned_model` - The name of the fine-tuned model once training has completed successfully.
- `created_at` - The timestamp when the fine-tuning job was created.
- `updated_at` - The timestamp when the fine-tuning job was last updated.
- `finished_at` - The timestamp when the fine-tuning job finished, if applicable.
- `organization_id` - The organization that owns the fine-tuning job.
- `result_files` - IDs of result files created by the fine-tuning job.
- `validation_files` - IDs of validation files used in the fine-tuning job.
- `training_files` - IDs of training files used in the fine-tuning job.
- `events` - Array of events related to the fine-tuning job, each containing:
  - `object` - The object type, always "fine-tuning.job.event".
  - `created_at` - The timestamp when the event was created.
  - `level` - The event level, e.g., "info", "warning", or "error".
  - `message` - A human-readable description of the event.
  - `data` - Additional data related to the event, if available.

## Import

Fine-tuning jobs can be imported using the OpenAI fine-tune ID:

```
$ terraform import openai_fine_tune.example ft-abc123
```
