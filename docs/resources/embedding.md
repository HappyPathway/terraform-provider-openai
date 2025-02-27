---
page_title: "openai_embedding Resource - terraform-provider-openai"
description: |-
  Creates embeddings with OpenAI models.
---

# openai_embedding (Resource)

Creates embeddings with OpenAI models. Given a piece of text, the model will generate a vector representation (embedding) that captures its semantic meaning.

~> **Note:** This resource creates an embedding each time it is created and cannot be updated. Changes to any of the arguments will result in a new embedding being created.

## Example Usage

```terraform
resource "openai_embedding" "example" {
  model = "text-embedding-ada-002"
  input = "The quick brown fox jumps over the lazy dog"
}

output "embedding_vector" {
  value = openai_embedding.example.embeddings[0].embedding
}
```

## Argument Reference

The following arguments are supported:

* `model` - (Required, Forces new resource) ID of the model to use for generating embeddings (e.g., "text-embedding-ada-002").
* `input` - (Required, Forces new resource) The text to generate embeddings for.
* `dimensions` - (Optional, Forces new resource) The number of dimensions the resulting embeddings should have. Only supported in `text-embedding-3` and later models.
* `encoding_format` - (Optional, Forces new resource) The format to return the embeddings in. Can be either "float" (default) or "base64".
* `user` - (Optional, Forces new resource) A unique identifier representing your end-user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique identifier for this embedding operation.
* `object` - The object type, which is always "list".
* `embeddings` - A list of generated embeddings. Each contains:
  * `embedding` - The embedding vector, represented as an array of floats.
  * `index` - The index of this embedding in the list.
* `usage` - Information about token usage, contains:
  * `prompt_tokens` - Number of tokens in the input text.
  * `total_tokens` - Total number of tokens used.