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

resource "openai_file" "code_file" {
  filename  = "analysis.py"
  file_path = "${path.module}/analysis.py"
  purpose   = "assistants"
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
      file_ids = [openai_file.data_file.id, openai_file.code_file.id]
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
