---
page_title: "openai_models Data Source - terraform-provider-openai"
description: |-
  Get a list of all available OpenAI models.
---

# openai_models (Data Source)

Use this data source to get a list of all models available to your organization through OpenAI's API.

## Example Usage

```terraform
data "openai_models" "all" {}

output "available_models" {
  value = data.openai_models.all.models[*].id
}

# Filter GPT-4 models
output "gpt4_models" {
  value = [for model in data.openai_models.all.models : model.id if startswith(model.id, "gpt-4")]
}
```

## Argument Reference

This data source has no arguments.

## Attributes Reference

* `object` - The object type, which is always "list".
* `data` - A list of model objects. Each model contains:
  * `id` - The ID of the model.
  * `created` - The Unix timestamp (in seconds) when the model was created.
  * `owned_by` - The organization that owns the model.
  * `object` - The object type, which is always "model".
  * `permissions` - A list of permissions that the model has. Each permission contains:
    * `id` - The ID of the permission.
    * `object` - The object type, which is always "model_permission".
    * `created` - The Unix timestamp (in seconds) when the permission was created.
    * `allow_create_engine` - Whether the model can be used to create engines.
    * `allow_fine_tuning` - Whether the model can be fine-tuned.
    * `allow_logprobs` - Whether the model can return log probabilities.
    * `allow_sampling` - Whether the model supports sampling.
    * `allow_search_indices` - Whether the model can be used to search indices.
    * `allow_view` - Whether the model can be viewed.
    * `is_blocking` - Whether the model is blocking.
    * `organization` - The organization ID that this permission applies to.