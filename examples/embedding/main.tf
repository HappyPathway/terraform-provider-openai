terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

resource "openai_embedding" "text_embedding" {
  model = "text-embedding-ada-002"
  input = "The quick brown fox jumps over the lazy dog"
}

# Example of batch embedding
resource "openai_embedding" "batch_embedding" {
  model = "text-embedding-ada-002"
  input = [
    "The quick brown fox jumps over the lazy dog",
    "Pack my box with five dozen liquor jugs",
    "How vexingly quick daft zebras jump"
  ]
}

output "single_embedding" {
  value = openai_embedding.text_embedding.embedding[0]
}

output "batch_embeddings" {
  value = openai_embedding.batch_embedding.embedding
}