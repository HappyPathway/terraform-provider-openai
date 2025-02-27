---
page_title: "openai_model Data Source - OpenAI Provider"
subcategory: ""
description: |-
  Get information about a specific OpenAI model.
---

# Data Source: openai_model

Use this data source to get information about a specific OpenAI model.

## Example Usage

```terraform
data "openai_model" "gpt4" {
  id = "gpt-4"
}
```

## Argument Reference

* `id` - (Required) The ID of the model to retrieve information about.

## Attributes Reference

* `owned_by` - The organization that owns the model.
* `permission` - List of permissions for the model. Each permission contains:
  * `allow_create_engine` - Whether the model can be used to create engines.
  * `allow_fine_tuning` - Whether the model can be fine-tuned.
  * `allow_sampling` - Whether the model allows sampling.
  * `allow_search_indices` - Whether the model can be used for search indices.
  * `allow_view` - Whether the model can be viewed.
  * `created` - The timestamp when the permission was created.
  * `group` - The group the permission belongs to.
  * `is_blocking` - Whether the permission is blocking.
  * `organization` - The organization the permission applies to.