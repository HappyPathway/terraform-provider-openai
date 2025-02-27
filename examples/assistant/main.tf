terraform {
  required_version = ">= 1.0.0"
  required_providers {
    openai = {
      source  = "HappyPathway/openai"
    }
  }
}

provider "openai" {
  # Configuration options will be populated from environment variables
  # OPENAI_API_KEY or from .terraformrc credentials block
}

resource "openai_file" "knowledge_base" {
  content  = filebase64("${path.module}/knowledge_base.txt")
  filename = "knowledge_base.txt"
  purpose  = "assistants"
}

resource "openai_assistant" "research_assistant" {
  name         = "Research Assistant"
  description  = "A research assistant that helps with data analysis and academic writing"
  model        = "gpt-4-1106-preview"
  instructions = "You are a helpful research assistant. Use the provided knowledge base to answer questions accurately."

  tools {
    type = "retrieval"
  }

  tools {
    type = "function"
    
    function {
      name = "search_papers"
      description = "Search for relevant research papers in the database"
      parameters = jsonencode({
        type = "object"
        properties = {
          query = {
            type = "string"
            description = "The search query"
          }
          year = {
            type = "integer"
            description = "Optional: Filter by publication year"
          }
        }
        required = ["query"]
      })
    }
  }

  file_ids = [openai_file.knowledge_base.id]
}

resource "openai_assistant" "infra_assistant" {
  name         = "Infrastructure Assistant"
  description  = "An assistant that helps with infrastructure and DevOps tasks"
  model        = "gpt-4-1106-preview"
  instructions = "You are a helpful infrastructure and DevOps assistant. Use the provided knowledge base to answer questions accurately and provide practical solutions."

  tools {
    type = "retrieval"
  }

  tools {
    type = "code_interpreter"
  }

  file_ids = [openai_file.knowledge_base.id]
}

output "assistant_id" {
  value       = openai_assistant.research_assistant.id
  description = "The ID of the created OpenAI Assistant"
}

output "infra_assistant_id" {
  value = openai_assistant.infra_assistant.id
}