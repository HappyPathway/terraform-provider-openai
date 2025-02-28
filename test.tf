terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
  required_version = ">= 0.13"
}

# Example of using the file resource
resource "openai_file" "training_data" {
  file = "training_data.jsonl"
  purpose = "fine-tune"
}

# Example of using the data source to get information about a model
data "openai_model" "gpt4" {
  model_id = "gpt-4"
}

output "model_info" {
  value = data.openai_model.gpt4
}

provider "openai" {
  # Configuration will be loaded from environment variables:
  # OPENAI_API_KEY for api_key
  # OPENAI_ORGANIZATION_ID for organization_id (optional)
  retry_max   = 3
  retry_delay = 5
  timeout     = 30
}

# Data Sources for Models
data "openai_models" "available" {
  # Lists all available models
}

# File Resources
resource "openai_file" "knowledge_base" {
  file    = "${path.module}/knowledge_base.txt"
  purpose = "assistants"
}

resource "openai_file" "fine_tuning_data" {
  file    = "${path.module}/training_data.jsonl"
  purpose = "fine-tune"
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

# Thread Resource with All Capabilities
resource "openai_thread" "example" {
  metadata = {
    purpose = "example"
    env     = "development"
  }

  # Initial messages for the thread
  messages {
    role    = "user"
    content = "Initialize analysis with the provided knowledge base."
  }

  tool_resources {
    code_interpreter {
      file_ids = [openai_file.knowledge_base.id]
    }
    file_search {
      vector_store_ids = [openai_beta_vector_store.example.id]
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Message Resources - Different types of messages
resource "openai_message" "initial_query" {
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

resource "openai_message" "with_image" {
  thread_id = openai_thread.example.id
  role      = "user"
  
  content {
    type = "text"
    text = "Analyzing the following image"
  }

  file_ids = [openai_file.knowledge_base.id]
  metadata = {
    type = "image_analysis"
  }
}

# Assistant Resource with All Tool Types
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

# Outputs for important resource information
output "available_models" {
  value = data.openai_models.available.models
}

output "gpt4_model" {
  value = data.openai_model.gpt4
}

output "thread_info" {
  value = {
    id = openai_thread.example.id
    metadata = openai_thread.example.metadata
    created_at = openai_thread.example.created_at
  }
}

output "assistant_info" {
  value = {
    id = openai_assistant.data_analyst.id
    name = openai_assistant.data_analyst.name
    created_at = openai_assistant.data_analyst.created_at
  }
}

output "vector_store_info" {
  value = {
    id = openai_beta_vector_store.example.id
    name = openai_beta_vector_store.example.name
  }
}