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
  filename  = "data.json"
  file_path = "${path.module}/data.json"
  purpose   = "assistants"
}

resource "openai_file" "secondary_data_file" {
  filename  = "secondary_data.json"
  file_path = "${path.module}/secondary_data.json"
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
    days   = 90
    anchor = "last_active_at" # This will make it expire 90 days after last use
  }
}

# Add files to the vector store for semantic search
resource "openai_vector_store_file" "data_vectors" {
  vector_store_id = openai_vector_store.analysis_store.id
  file_id         = openai_file.data_file.id

  depends_on = [
    openai_file.data_file,
    openai_vector_store.analysis_store
  ]
}

resource "openai_vector_store_file" "secondary_data_vectors" {
  vector_store_id = openai_vector_store.analysis_store.id
  file_id         = openai_file.secondary_data_file.id

  depends_on = [
    openai_file.secondary_data_file,
    openai_vector_store.analysis_store
  ]
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
      file_ids = toset([
        openai_file.data_file.id,
        openai_file.secondary_data_file.id,
        openai_file.code_file.id
      ])
    }
  }

  depends_on = [
    openai_file.data_file,
    openai_file.secondary_data_file,
    openai_file.code_file
  ]

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
  tool_resources {
    file_search {
      vector_store_ids = [openai_vector_store.analysis_store.id]
    }
    code_interpreter {
      file_ids = toset([openai_file.code_file.id])
    }
  }
}

resource "openai_thread" "analysis_session_file_search" {
  metadata = {
    session_type = "data_analysis"
    project      = "example"
  }

  # Configure the thread to use the vector store for semantic search
  tool_resources {
    file_search {
      vector_store_ids = [openai_vector_store.analysis_store.id]
    }
  }
}

resource "openai_thread" "analysis_session_code_interpreter" {
  metadata = {
    session_type = "data_analysis"
    project      = "example"
  }

  # Configure the thread to use the vector store for semantic search
  tool_resources {
    code_interpreter {
      file_ids = toset([openai_file.code_file.id])
    }
  }
}

# Create a thread without tools
resource "openai_thread" "analysis_session_no_tools" {
  metadata = {
    session_type = "data_analysis"
    project      = "example"
  }
}

# Create a message in the thread
resource "openai_message" "initial_message" {
  thread_id = openai_thread.analysis_session.id
  role      = "user"
  content   = "Please analyze the data in data.json using the provided analysis.py script."

  # attachment {
  #   file_id = openai_file.data_file.id
  #   tools   = ["code_interpreter"]
  # }
}

# Output vector store information for reference
output "vector_store_info" {
  value = {
    id          = openai_vector_store.analysis_store.id
    name        = openai_vector_store.analysis_store.name
    status      = openai_vector_store.analysis_store.status
    usage_bytes = openai_vector_store.analysis_store.usage_bytes
  }
  description = "Information about the vector store and its contents"
}
