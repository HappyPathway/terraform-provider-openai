terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {}

# Upload a file for the assistant to use
resource "openai_file" "knowledge_base" {
  file_path = "${path.module}/data/knowledge_base.json"
  filename  = "knowledge_base.pdf"
  purpose   = "assistants"
}

# Create an assistant for customer support
resource "openai_assistant" "customer_support" {
  name         = "Customer Support Assistant"
  description  = "An assistant that helps with customer inquiries about our products"
  model        = "gpt-4-1106-preview"
  instructions = <<-EOT
    You are a customer support assistant for a technology company.
    
    Follow these guidelines:
    1. Be friendly and professional
    2. Answer questions based on the provided knowledge base
    3. If you don't know the answer, say so and offer to escalate to a human agent
    4. Don't make up information not in the knowledge base
    5. Format responses with markdown when helpful
  EOT

  tools {
    type = "file_search"
  }
  
  tools {
    type = "code_interpreter"
  }

  # Attach the knowledge base file
  file_path = openai_file.knowledge_base.file_path

  # Add metadata for organization
  metadata = {
    department = "customer_support"
    team       = "technical"
    version    = "1.0"
  }
}

output "assistant_id" {
  value       = openai_assistant.customer_support.id
  description = "The ID of the created assistant"
}

# Note: You need to create a data/knowledge_base.pdf file before running this example
