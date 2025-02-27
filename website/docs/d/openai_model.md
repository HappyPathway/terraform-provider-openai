---
page_title: "OpenAI: openai_model Data Source"
subcategory: ""
description: |-
  Get information about a specific OpenAI model.
---

# openai_model Data Source

Use this data source to get information about a specific OpenAI model.

## Example Usage

```hcl
data "openai_model" "gpt4" {
  model_id = "gpt-4"
}
```

## Argument Reference

* `model_id` - (Required) The ID of the model to retrieve information about.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `owned_by` - The organization that owns the model.
* `permission` - A list of permissions that the model has. Each permission contains:
  * `allow_create_engine` - Whether the model can be used to create an engine.
  * `allow_fine_tuning` - Whether the model can be fine-tuned.
  * `allow_sampling` - Whether the model supports sampling.
  * `allow_search_indices` - Whether the model can be used for search indices.
  * `allow_view` - Whether the model can be viewed.
  * `created` - The timestamp when the permission was created.
  * `group` - The group that the permission belongs to.
  * `id` - The ID of the permission.
  * `is_blocking` - Whether the permission is blocking.
  * `organization` - The organization that the permission belongs to.
* `root` - The root model that this model is based on.
* `created` - The timestamp when the model was created.