---
page_title: "openai_vector_store_file Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Manages a file within an OpenAI vector store.
---

# openai_vector_store_file

This resource manages files within an OpenAI vector store. When a file is added to a vector store, it is automatically processed and its contents are converted into embeddings for efficient similarity search.

## Example Usage

```terraform
# Create a vector store
resource "openai_vector_store" "knowledge_base" {
  name = "company-knowledge-base"
}

# Upload files to OpenAI
resource "openai_file" "documentation" {
  filename = "documentation.pdf"
  filepath = "./docs/documentation.pdf"
  purpose  = "assistants"
}

resource "openai_file" "procedures" {
  filename = "procedures.pdf"
  filepath = "./docs/procedures.pdf"
  purpose  = "assistants"
}

# Add files to the vector store
resource "openai_vector_store_file" "docs" {
  vector_store_id = openai_vector_store.knowledge_base.id
  file_id        = openai_file.documentation.id
}

resource "openai_vector_store_file" "procedures" {
  vector_store_id = openai_vector_store.knowledge_base.id
  file_id        = openai_file.procedures.id
}
```

## Argument Reference

- `vector_store_id` - (Required, Forces new resource) The ID of the vector store to add the file to.
- `file_id` - (Required, Forces new resource) The ID of the file to add to the vector store. The file must already exist in OpenAI.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for this file within the vector store.
- `created_at` - The Unix timestamp (in seconds) when the file was added to the vector store.
- `status` - The current status of the file in the vector store. Can be "processing", "succeeded", "failed", or "cancelled".
- `usage_bytes` - The size of the file in bytes.

## Import

Vector store files can be imported using the format `vector_store_id:file_id`:

```shell
$ terraform import openai_vector_store_file.docs vs_abc123:file-xyz789
```

## Timeouts

This resource supports the following timeouts:

- `create` - (Default `10 minutes`) Used when creating the vector store file and waiting for it to be processed.
- `read` - (Default `5 minutes`) Used when reading the vector store file.
- `delete` - (Default `5 minutes`) Used when deleting the vector store file.
