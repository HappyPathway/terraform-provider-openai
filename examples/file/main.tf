terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {}

// Upload a file for fine-tuning
resource "openai_file" "fine_tune_data" {
  file_path = "${path.module}/data/fine_tune_data.jsonl"
  filename  = "fine_tune_data.jsonl"
  purpose   = "fine-tune"
}

// Upload a file for assistant retrieval
resource "openai_file" "assistant_knowledge" {
  file_path = "${path.module}/data/knowledge_base.pdf"
  filename  = "knowledge_base.pdf"
  purpose   = "assistants"
}

output "fine_tune_file_id" {
  value       = openai_file.fine_tune_data.id
  description = "The OpenAI-assigned file ID for the fine-tuning file"
}

output "assistant_file_id" {
  value       = openai_file.assistant_knowledge.id
  description = "The OpenAI-assigned file ID for the assistant knowledge file"
}

// Note: You need to create the following files before running this example:
// - data/fine_tune_data.jsonl - JSONL file with fine-tuning data
// - data/knowledge_base.pdf - PDF file with knowledge for assistants
