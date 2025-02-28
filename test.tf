terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {
  # Configuration will be loaded from environment variables:
  # OPENAI_API_KEY for api_key
  # OPENAI_ORGANIZATION_ID for organization_id (optional)
  retry_max   = 3
  retry_delay = 5
  timeout     = 30
}

# Data Sources
data "openai_models" "available" {
  # Lists all available models
}

data "openai_model" "gpt4" {
  model_id = "gpt-4"
}

# File Resource - Used by other resources
resource "openai_file" "knowledge_base" {
  file    = "${path.module}/knowledge_base.txt"
  purpose = "assistants"
}

# Vector Store Resource - For file search capability
resource "openai_beta_vector_store" "example" {
  name = "test-store-1"
  metadata = {
    purpose = "example"
    env     = "development"
  }
  file_ids = [openai_file.knowledge_base.id]
}

# Thread Resource - For managing conversations
resource "openai_thread" "example" {
  metadata = {
    purpose = "example"
    env     = "development"
  }

  tool_resources {
    code_interpreter {
      file_ids = [openai_file.knowledge_base.id]
    }
    file_search {
      vector_store_ids = [openai_beta_vector_store.example.id]
    }
  }

  # Move message to separate resource to avoid recreation cycles
  lifecycle {
    create_before_destroy = true
  }
}

# Message Resource - For adding messages to threads
resource "openai_message" "initial" {
  thread_id = openai_thread.example.id
  role      = "user"
  
  content {
    type = "text"
    text = "What insights can you provide from the knowledge base?"
  }

  file_ids = [openai_file.knowledge_base.id]
  metadata = {
    type = "initial_query"
    tags = "knowledge_base,analysis"
  }
}

# Assistant Resource - AI assistant with various capabilities
resource "openai_assistant" "data_analyst" {
  name         = "Data Analysis Assistant"
  description  = "An assistant that helps with data analysis and visualization"
  model        = "gpt-4-1106-preview"
  instructions = "You are a data analysis expert. Use the provided tools to analyze data and create visualizations."

  tools {
    type = "code_interpreter"
  }

  tools {
    type = "retrieval"
  }

  tools {
    type        = "function"
    name        = "analyze_dataset"
    description = "Analyze a dataset with specified parameters"
    parameters  = jsonencode({
      type = "object"
      properties = {
        dataset_name = {
          type        = "string"
          description = "Name of the dataset to analyze"
        }
        analysis_type = {
          type        = "string"
          enum        = ["statistical", "temporal", "categorical"]
          description = "Type of analysis to perform"
        }
        output_format = {
          type        = "string"
          enum        = ["json", "csv", "text"]
          description = "Desired output format"
        }
      }
      required = ["dataset_name", "analysis_type"]
    })
  }

  file_ids = [openai_file.knowledge_base.id]
  metadata = {
    capability     = "data_analysis"
    specialization = "statistical_analysis"
    version       = "1.0"
  }
}

# Outputs
output "test_thread" {
  value = {
    id = openai_thread.example.id
    metadata = openai_thread.example.metadata
    created_at = openai_thread.example.created_at
  }
}

output "test_message" {
  value = {
    id = openai_message.initial.id
    thread_id = openai_message.initial.thread_id
    status = openai_message.initial.status
    created_at = openai_message.initial.created_at
  }
}