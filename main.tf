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

# Upload a file for use with the OpenAI API
resource "openai_file" "data_file" {
  filename = "data.json"
  content  = jsonencode({ "example" : "data" })
  purpose  = "assistants"
}

# Create a thread for analysis
resource "openai_thread" "analysis_session" {
  metadata = {
    session_type = "data_analysis"
    project      = "example"
  }

  tool_resources = {
    code_interpreter = {
      file_ids = []
    }
    file_search = {
      vector_store_ids = [] # Empty list of strings
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

  attachment = [
    {
      file_id = openai_file.data_file.id
      tools   = ["code_interpreter"]
    }
  ]
}
