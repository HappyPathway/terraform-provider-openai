---
page_title: "openai_model Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  Retrieve information about OpenAI models.
---

# openai_model Data Source

Retrieves information about available OpenAI models and their capabilities. This data source can be used to find the appropriate model for your specific needs.

## Example Usage

```terraform
# Get information about a specific model by ID
data "openai_model" "gpt4" {
  model_id = "gpt-4"
}

# Filter models by type (e.g., get all chat models)
data "openai_model" "chat_models" {
  filter_by_type = "chat"
}

output "model_details" {
  value = data.openai_model.gpt4
}
```

## Argument Reference

- `model_id` - (Optional) The ID of a specific model to retrieve information about. If provided, the data source will return details for only this model.
- `filter_by_type` - (Optional) Filter models by type. Supported values include "chat", "embedding", "completion", etc. If provided, the data source will return all models of the specified type.

## Attribute Reference

When `model_id` is specified, the following attributes are exported:

- `id` - The ID of the model.
- `object` - The object type, always "model".
- `created` - The timestamp when the model was created or made available by OpenAI.
- `owned_by` - The organization that owns the model.
- `permission` - A list of permissions associated with the model, each containing:
  - `id` - The permission ID.
  - `object` - The permission object type.
  - `created` - When the permission was created.
  - `allow_create_engine` - Whether creating an engine is allowed.
  - `allow_sampling` - Whether sampling is allowed.
  - `allow_logprobs` - Whether log probabilities are allowed.
  - `allow_search_indices` - Whether search indices are allowed.
  - `allow_view` - Whether viewing is allowed.
  - `allow_fine_tuning` - Whether fine-tuning is allowed.
  - `organization` - The organization this permission applies to.
  - `group` - The group this permission applies to.
  - `is_blocking` - Whether this permission is blocking.

When `filter_by_type` is specified or neither argument is provided, the data source exports:

- `models` - A list of models matching the filter criteria, each containing the same attributes as above.
- `total` - The total number of models returned.
