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

# Upload files for use with code interpreter and retrieval
resource "openai_file" "data_file" {
  filename  = "data.csv"
  file_path = "${path.module}/data.csv"
  purpose   = "assistants"
}

resource "openai_file" "secondary_data_file" {
  filename  = "secondary_data.csv"
  file_path = "${path.module}/secondary_data.csv"
  purpose   = "assistants"
}

resource "openai_file" "code_file" {
  filename  = "analysis.py"
  file_path = "${path.module}/analysis.py"
  purpose   = "assistants"
}

# Create a vector store for semantic search
resource "openai_vector_store" "analysis_store" {
  name = "data-analysis-store"

  metadata = {
    project     = "example"
    environment = "development"
    purpose     = "data-analysis"
  }

  # Optional: Configure expiration
  expires_after {
    days = 90
  }
}

# Add files to the vector store for semantic search
resource "openai_vector_store_file" "data_vectors" {
  vector_store_id = openai_vector_store.analysis_store.id
  file_id         = openai_file.data_file.id
}

resource "openai_vector_store_file" "secondary_data_vectors" {
  vector_store_id = openai_vector_store.analysis_store.id
  file_id         = openai_file.secondary_data_file.id
}

# Create an assistant that can use both code interpreter and file search
resource "openai_assistant" "data_analyst" {
  name         = "Data Analyst Assistant"
  model        = "gpt-4-turbo-preview"
  description  = "An assistant that helps analyze data using code interpreter and file search"
  instructions = "You are a data analysis assistant. Use the provided files and code interpreter to help analyze data and answer questions."

  tools = ["code_interpreter"]

  # Tool resources block is optional
  tool_resources {
    # Code interpreter block is optional
    code_interpreter {
      file_ids = [
        openai_file.data_file.id,
        openai_file.secondary_data_file.id,
        openai_file.code_file.id
      ]
    }
  }

  metadata = {
    session_type = "data_analysis"
    project      = "example"
  }
}

# Create a thread for analysis
resource "openai_thread" "analysis_session" {
  metadata = {
    session_type = "data_analysis"
    project      = "example"
  }

  # Configure the thread to use the vector store for semantic search
  tool_resources = {
    file_search = {
      vector_store_ids = [openai_vector_store.analysis_store.id]
    }
  }
}

# Create a message in the thread
resource "openai_message" "initial_message" {
  thread_id = openai_thread.analysis_session.id
  role      = "user"
  content   = "Please analyze the data in data.csv using the provided analysis.py script."

  attachments = [
    {
      file_id = openai_file.data_file.id
      tools   = ["code_interpreter"]
    }
  ]
}

# Output vector store information for reference
output "vector_store_info" {
  value = {
    id          = openai_vector_store.analysis_store.id
    name        = openai_vector_store.analysis_store.name
    status      = openai_vector_store.analysis_store.status
    file_counts = openai_vector_store.analysis_store.file_counts
    usage_bytes = openai_vector_store.analysis_store.usage_bytes
  }
  description = "Information about the vector store and its contents"
}
