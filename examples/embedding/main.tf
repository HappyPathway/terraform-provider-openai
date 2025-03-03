terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {}

# Generate embeddings for a piece of text
resource "openai_embedding" "text_embedding" {
  model = "text-embedding-ada-002"
  input = "Terraform is an infrastructure as code tool that enables you to safely and predictably manage infrastructure."
}

# Output the first 5 values of the embedding vector (truncated for readability)
output "embedding_sample" {
  value       = slice(openai_embedding.text_embedding.embedding, 0, 5)
  description = "First 5 dimensions of the embedding vector"
}

# Note: the multi-text embedding is not supported in the current version
# To embed multiple texts, create multiple resources:
