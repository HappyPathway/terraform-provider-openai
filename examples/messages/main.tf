terraform {
  required_providers {
    openai = {
      source = "openai/openai"
    }
  }
}

# Create a thread to hold the messages
resource "openai_thread" "example" {
  metadata = {
    purpose = "example"
  }
}

# Create a text message in the thread
resource "openai_message" "text_example" {
  thread_id = openai_thread.example.id
  role      = "user"
  
  content {
    type = "text"
    text = "What is machine learning?"
  }

  metadata = {
    message_type = "question"
  }
}

# Create a message with an image file attachment
resource "openai_message" "image_example" {
  thread_id = openai_thread.example.id
  role      = "user"

  content {
    type = "image_file"
    image_file {
      file_id = "file-abc123" # Replace with actual file ID
      detail  = "auto"
    }
  }

  # Attach the file for code interpretation
  attachments {
    file_id = "file-abc123" # Replace with actual file ID
    tool    = "code_interpreter"
  }
}

# Create a message with an image URL
resource "openai_message" "image_url_example" {
  thread_id = openai_thread.example.id
  role      = "user"

  content {
    type = "image_url"
    image_url {
      url    = "https://example.com/image.jpg"
      detail = "low"
    }
  }
}

# Example of an assistant message
resource "openai_message" "assistant_example" {
  thread_id = openai_thread.example.id
  role      = "assistant"

  content {
    type = "text"
    text = "Machine learning is a subset of artificial intelligence..."
  }
}