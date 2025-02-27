terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {
  # Configure your API key via OPENAI_API_KEY environment variable
}

# Test the models data source
data "openai_models" "available" {}

# Test specific model data source
data "openai_model" "gpt4" {
  model_id = "gpt-4"
}

# Test file resource
resource "openai_file" "test_file" {
  file    = "knowledge_base.txt"
  purpose = "assistants"
}

# Test assistant resource with comprehensive configuration
resource "openai_assistant" "test" {
  name         = "Data Analysis Assistant"
  description  = "A sophisticated assistant for data analysis and visualization"
  model        = "gpt-4-turbo-preview"
  instructions = <<-EOT
    You are a specialized data analysis assistant with the following capabilities:
    1. Analyze and interpret data using Python and other tools
    2. Create visualizations and charts
    3. Search through provided documentation and data files
    4. Execute custom functions for data processing
    Please provide clear explanations with your analysis and always show your work.
  EOT

  tools {
    type = "code_interpreter"
  }

  tools {
    type = "retrieval"
  }

  tools {
    type = "function"
    function {
      name        = "process_data"
      description = "Process and analyze dataset with custom parameters"
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
  }

  file_ids = [openai_file.test_file.id]

  metadata = {
    capability     = "data_analysis"
    version       = "1.0"
    specialization = "statistical_analysis"
  }
}

# Test content generator resource
resource "openai_content_generator" "test" {
  model       = "gpt-3.5-turbo"
  temperature = 0.7

  messages {
    role    = "user"
    content = "What is 2 + 2? Respond with just the number."
  }
}

# Test content generator resource with JSON response
resource "openai_content_generator" "json_test" {
  model       = "gpt-4-turbo-preview"
  temperature = 0.7

  messages {
    role    = "system"
    content = <<-EOT
      You are a helpful assistant that provides structured data about movies.
      Always respond with valid JSON that has the following structure:
      {
        "title": "string",
        "year": number,
        "directors": ["string"],
        "main_actors": [{"name": "string", "character": "string"}],
        "rating": {"value": number, "source": "IMDb"}
      }
    EOT
  }

  messages {
    role    = "user"
    content = "Give me information about The Matrix movie from 1999."
  }

  response_format {
    type = "json_object"
  }
}

output "test" {
  value = openai_assistant.test
}

output "test_file" {
  value = openai_file.test_file
}

output "test_content_generator" {
  value = openai_content_generator.test
}

output test_content_generator_anser {
  value = jsondecode(openai_content_generator.test.raw_response)
}

# Output showing both raw and parsed JSON response
output "matrix_movie_raw" {
  value = openai_content_generator.json_test.raw_response
}

output "matrix_movie_content" {
  value = jsondecode(openai_content_generator.json_test.content)
}

output "matrix_movie_actors" {
  description = "Just the main actors from the movie"
  value = jsondecode(openai_content_generator.json_test.content).main_actors
}

output "matrix_movie_rating" {
  description = "Just the rating information"
  value = jsondecode(openai_content_generator.json_test.content).rating
}

output "models" {
  value = data.openai_models.available
}

output "model" {
  value = data.openai_model.gpt4
}