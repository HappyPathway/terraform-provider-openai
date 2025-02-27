---
page_title: "openai_model Data Source - terraform-provider-openai"
description: |-
  Get information about a specific OpenAI model.
---

# openai_model (Data Source)

Use this data source to get information about a specific OpenAI model, such as its permissions and capabilities.

## Example Usage

```terraform
data "openai_model" "gpt4" {
  id = "gpt-4"
}

output "model_owned_by" {
  value = data.openai_model.gpt4.owned_by
}
```

## Argument Reference

* `id` - (Required) The ID of the model to retrieve information for.

## Attributes Reference

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