---
page_title: "openai_fine_tuning_job Resource - OpenAI Provider"
subcategory: ""
description: |-
  Manages an OpenAI Fine-tuning Job to create custom models.
---

# Resource: openai_fine_tuning_job

This resource allows you to create and manage fine-tuning jobs to customize OpenAI models with your own training data.

## Example Usage

```terraform
resource "openai_file" "training_data" {
  file    = "${path.module}/training_data.jsonl"
  purpose = "fine-tune"
}

resource "openai_fine_tuning_job" "custom_model" {
  model = "gpt-3.5-turbo"
  training_file = openai_file.training_data.id
  
  hyperparameters {
    n_epochs = 3
  }

  validation_file = openai_file.validation_data.id
}
```

## Argument Reference

* `model` - (Required) The name of the model to fine-tune. Must be one of the models that supports fine-tuning.
* `training_file` - (Required) The ID of the file to use for training.
* `hyperparameters` - (Optional) The hyperparameters used for the fine-tuning job. Supports:
  * `n_epochs` - (Optional) Number of epochs to train the model for.
* `validation_file` - (Optional) The ID of the file to use for validation.
* `suffix` - (Optional) A string that will be added to the fine-tuned model name.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the fine-tuning job.
* `status` - The current status of the fine-tuning job.
* `fine_tuned_model` - The name of the resulting fine-tuned model, once the job is completed.
* `created_at` - The Unix timestamp (in seconds) for when the fine-tuning job was created.
* `finished_at` - The Unix timestamp (in seconds) for when the fine-tuning job was completed.

## Import

Fine-tuning jobs can be imported using the job ID, e.g.,

```shell
terraform import openai_fine_tuning_job.example ftjob-abc123
```