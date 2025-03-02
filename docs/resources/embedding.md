---
page_title: "openai_embedding Resource - terraform-provider-openai"
subcategory: ""
description: |-
  Generate embeddings for text using OpenAI's embedding models.
---

# openai_embedding

Generates vector embeddings for input text using OpenAI's embedding models. These embeddings can be used for various natural language processing tasks, like semantic search, text similarity, and classification.

## Example Usage

```terraform
resource "openai_embedding" "example" {
  model = "text-embedding-ada-002"
  input = "The food was delicious and the service was excellent."
}

output "embedding_vector" {
  value     = openai_embedding.example.embedding
  sensitive = true
}
```

## Argument Reference

- `model` - (Required) ID of the model to use for generating embeddings (e.g., 'text-embedding-ada-002').
- `input` - (Required) The text to generate embeddings for. This can be a single string or an array of strings.
- `encoding_format` - (Optional) The format in which the embeddings are returned. Can be either "float" (default) or "base64".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - Unique identifier for this resource.
- `embedding` - The generated embedding vector(s). For a single input string, this will be a list of floating-point numbers.
- `usage` - Information about token usage, containing:
  - `prompt_tokens` - Number of tokens in the input text.
  - `total_tokens` - Total number of tokens processed.

## Import

This resource does not support import as it is stateless and generates new embeddings on each apply.
