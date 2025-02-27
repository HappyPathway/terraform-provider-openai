---
page_title: "openai_models Data Source - OpenAI Provider"
subcategory: ""
description: |-
  Get information about all available OpenAI models.
---

# Data Source: openai_models

Use this data source to get information about all available OpenAI models. This is useful when you need to list all models available to your organization.

## Example Usage

```terraform
data "openai_models" "all" {}

output "available_models" {
  value = data.openai_models.all.models
}
```

## Argument Reference

This data source has no arguments.

## Attributes Reference

* `models` - A list of all available models. Each model contains:
  * `id` - The ID of the model.
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