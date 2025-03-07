---
page_title: "openai_vector_store Data Source - terraform-provider-openai"
subcategory: ""
description: |-
  Use this data source to get information about an existing OpenAI vector store.
---

# openai_vector_store (Data Source)

This data source provides information about an existing OpenAI vector store. Use this to retrieve information about a vector store for use in other resources or outputs.

## Example Usage

```terraform
# Get information about a specific vector store
data "openai_vector_store" "existing" {
  id = "vs_abc123"
}

# Use the vector store information in outputs
output "store_status" {
  value = data.openai_vector_store.existing.status
}

output "total_usage" {
  value = data.openai_vector_store.existing.usage_bytes
}

output "file_stats" {
  value = {
    total_files = data.openai_vector_store.existing.file_counts.total
    completed   = data.openai_vector_store.existing.file_counts.completed
    failed      = data.openai_vector_store.existing.file_counts.failed
  }
}
```

## Argument Reference

- `id` - (Required) The ID of the vector store to retrieve information about.

## Attribute Reference

- `name` - The name of the vector store.
- `created_at` - The Unix timestamp (in seconds) when the vector store was created.
- `status` - The current status of the vector store.
- `file_counts` - Statistics about files in the vector store:
  - `in_progress` - Number of files currently being processed.
  - `completed` - Number of successfully processed files.
  - `failed` - Number of files that failed processing.
  - `cancelled` - Number of cancelled file operations.
  - `total` - Total number of files.
- `usage_bytes` - The total size of the vector store in bytes.
- `expires_at` - The Unix timestamp (in seconds) when the vector store will expire (if expiration is configured).
- `expires_after` - Configuration block showing vector store expiration settings:
  - `days` - Number of days after which the vector store expires.
  - `anchor` - Reference time for expiration calculation.
- `metadata` - A map of metadata key-value pairs associated with the vector store.
