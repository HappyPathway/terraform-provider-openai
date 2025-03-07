terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {
  # Configure the OpenAI Provider
}

# Create a vector store for knowledge base management
resource "openai_vector_store" "knowledge_base" {
  name = "company-kb"

  metadata = {
    environment = "production"
    department  = "documentation"
    category    = "technical-docs"
  }

  # Configure expiration for compliance
  expires_after {
    days   = 365
    anchor = "now"
  }
}

# Upload documentation files
resource "openai_file" "product_docs" {
  filename  = "product-documentation.pdf"
  file_path = "${path.module}/files/product-documentation.pdf"
  purpose   = "assistants"
}

resource "openai_file" "api_docs" {
  filename  = "api-documentation.md"
  file_path = "${path.module}/files/api-documentation.md"
  purpose   = "assistants"
}

resource "openai_file" "user_guides" {
  filename  = "user-guides.pdf"
  file_path = "${path.module}/files/user-guides.pdf"
  purpose   = "assistants"
}

# Add files to the vector store for semantic search
resource "openai_vector_store_file" "product_docs_vectors" {
  vector_store_id = openai_vector_store.knowledge_base.id
  file_id         = openai_file.product_docs.id
}

resource "openai_vector_store_file" "api_docs_vectors" {
  vector_store_id = openai_vector_store.knowledge_base.id
  file_id         = openai_file.api_docs.id
}

resource "openai_vector_store_file" "user_guides_vectors" {
  vector_store_id = openai_vector_store.knowledge_base.id
  file_id         = openai_file.user_guides.id
}

# Use data source to get information about the vector store
data "openai_vector_store" "kb_info" {
  id = openai_vector_store.knowledge_base.id
}

# Output detailed information about the vector store
output "knowledge_base_stats" {
  value = {
    total_files     = data.openai_vector_store.kb_info.file_counts.total
    processed_files = data.openai_vector_store.kb_info.file_counts.completed
    total_size      = "${data.openai_vector_store.kb_info.usage_bytes / 1024 / 1024} MB"
    status          = data.openai_vector_store.kb_info.status
    expires_at      = data.openai_vector_store.kb_info.expires_at
  }
  description = "Statistics about the knowledge base vector store"
}

# Create an assistant that uses the vector store for document search
resource "openai_assistant" "documentation_assistant" {
  name         = "Documentation Helper"
  model        = "gpt-4-turbo-preview"
  description  = "An assistant that helps users find and understand documentation"
  instructions = "You are a documentation assistant. Use the company knowledge base to help answer questions about our products, APIs, and user guides."

  tool_resources {
    file_search {
      vector_store_ids = [openai_vector_store.knowledge_base.id]
    }
  }

  metadata = {
    type       = "documentation_assistant"
    department = "technical_support"
  }
}
