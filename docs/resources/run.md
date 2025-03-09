---
page_title: "openai_run Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Creates and manages an OpenAI Assistant run.
---

# openai_run (Resource)

Creates and manages an OpenAI Assistant run. A run represents the execution of an assistant on a thread.

## Example Usage

```hcl
# Create a basic assistant
resource "openai_assistant" "data_viz" {
  name = "Data Visualizer"
  instructions = "You are great at creating beautiful data visualizations."
  model = "gpt-4-turbo-preview"
  tools = ["code_interpreter"]
}

# Create a thread
resource "openai_thread" "example" {
  messages = [
    {
      role = "user"
      content = "Can you help me analyze this data?"
    }
  ]
}

# Create a run
resource "openai_run" "example" {
  assistant_id = openai_assistant.data_viz.id
  thread_id = openai_thread.example.id

  # Optional: Override assistant defaults for this run
  model = "gpt-4-turbo-preview"
  instructions = "Focus on time series analysis"

  # Control run completion behavior
  wait_for_completion = true
  polling_interval = "5s"
  timeout = "10m"
}
```

## Argument Reference

- `assistant_id` - (Required) The ID of the assistant to use for this run.
- `thread_id` - (Required) The ID of the thread to run the assistant on.
- `model` - (Optional) Override the default model used by the assistant.
- `instructions` - (Optional) Override the default instructions of the assistant for this run.
- `tools` - (Optional) Override the default tools of the assistant for this run.
- `wait_for_completion` - (Optional) Whether to wait for the run to complete before marking the resource as created. Defaults to true.
- `polling_interval` - (Optional) How often to poll for run status when wait_for_completion is true. Defaults to 5s.
- `timeout` - (Optional) Maximum time to wait for run completion when wait_for_completion is true. Defaults to 10m.
- `metadata` - (Optional) Set of key-value pairs that can be used to store additional information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the run.
- `status` - The status of the run (queued, in_progress, completed, requires_action, expired, cancelling, cancelled, failed).
- `created_at` - Unix timestamp for when the run was created.
- `expires_at` - Unix timestamp for when the run will expire.
- `started_at` - Unix timestamp for when the run was started.
- `cancelled_at` - Unix timestamp for when the run was cancelled.
- `failed_at` - Unix timestamp for when the run failed.
- `completed_at` - Unix timestamp for when the run completed.
- `last_error` - The last error message if the run failed.
- `steps` - The list of step IDs taken during the run.
- `required_action` - Details about any required actions needed to continue the run.

## Import

Runs can be imported using their ID:

```shell
terraform import openai_run.example run_abc123
```

## Notes

- Runs cannot be updated after creation. Any changes will force creation of a new run.
- When a run is deleted, it is cancelled if still in progress.
- Setting `wait_for_completion = true` (the default) means Terraform will wait for the run to complete before considering the resource created. This ensures any outputs or state changes from the run are captured.
