---
page_title: "openai_vector_store Resource - terraform-provider-openai"
subcategory: ""
description: |-
  A vector store resource for storing and managing embeddings.
---

# openai_vector_store

A vector store is a specialized database optimized for storing and searching vector embeddings. It enables efficient similarity search operations and supports semantic search functionality.

## Example Usage

```terraform
# Create a vector store
resource "openai_vector_store" "knowledge_base" {
  name = "company-knowledge-base"

  metadata = {
    environment = "production"
    department  = "engineering"
  }

  # Configure expiration (optional)
  expires_after {
    days   = 90
    anchor = "now"  # Optional
  }
}

# Add files to the vector store (optional)
resource "openai_file" "documentation" {
  filename = "api_docs.pdf"
  filepath = "./docs/api_docs.pdf"
  purpose  = "assistants"
}

resource "openai_vector_store_file" "doc_vectors" {
  vector_store_id = openai_vector_store.knowledge_base.id
  file_id        = openai_file.documentation.id
}
```

## Argument Reference

- `name` - (Required) The name of the vector store. This helps identify the store and its purpose.
- `metadata` - (Optional) A map of metadata key-value pairs to associate with the vector store. Use this to store custom attributes.
- `expires_after` - (Optional) Configuration block for setting up expiration of the vector store:
  - `days` - (Required) Number of days after which the vector store will expire.
  - `anchor` - (Optional) Reference time for expiration calculation. Defaults to "now".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the vector store.
- `created_at` - The Unix timestamp (in seconds) for when the vector store was created.
- `status` - The current status of the vector store.
- `usage_bytes` - The total size of the vector store in bytes.
- `expires_at` - The Unix timestamp (in seconds) when the vector store will expire (if expiration is configured).

## Import

Vector stores can be imported using their ID:

```shell
$ terraform import openai_vector_store.knowledge_base vs_abc123
```
