---
page_title: "openai_fine_tuning_job Resource - terraform-provider-openai"
description: |-
  Manages OpenAI fine-tuning jobs.
---

# openai_fine_tuning_job (Resource)

Manages fine-tuning jobs in OpenAI. Fine-tuning lets you customize a model with your training data to get better results for your specific use case.

## Example Usage

```terraform
resource "openai_file" "training_data" {
  content  = filebase64("${path.module}/training_data.jsonl")
  filename = "training_data.jsonl"
  purpose  = "fine-tune"
}

resource "openai_fine_tuning_job" "example" {
  model          = "gpt-3.5-turbo"
  training_file  = openai_file.training_data.id
  hyperparameters = {
    n_epochs = 3
  }
  suffix = "custom-model-v1"
}
```

## Argument Reference

The following arguments are supported:

* `model` - (Required, Forces new resource) The name of the model to fine-tune. You can use GPT-3.5-turbo or other models enabled for fine-tuning in your organization.
* `training_file` - (Required, Forces new resource) The ID of the file to use for training.
* `validation_file` - (Optional, Forces new resource) The ID of the file to use for validation.
* `hyperparameters` - (Optional, Forces new resource) The hyperparameters used for the fine-tuning job. Contains:
  * `n_epochs` - (Optional) Number of epochs to train for. Defaults to auto-detect.
* `suffix` - (Optional, Forces new resource) A string that will be added to the fine-tuned model name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the fine-tuning job.
* `status` - The current status of the fine-tuning job. Can be one of "validating", "queued", "running", "succeeded", "failed", or "cancelled".
* `created_at` - The Unix timestamp (in seconds) for when the fine-tuning job was created.
* `finished_at` - The Unix timestamp (in seconds) for when the fine-tuning job was completed, if it has completed.
* `error` - Error information, if the job failed. Contains:
  * `code` - The error code.
  * `message` - The error message.

## Import

Fine-tuning jobs can be imported using the job ID, e.g.:

```shell
$ terraform import openai_fine_tuning_job.example ftjob-abc123
```