terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {}

# Reference an existing assistant or create a new one
resource "openai_assistant" "support_assistant" {
  name         = "Support Assistant"
  description  = "An assistant that helps with support requests"
  model        = "gpt-4"
  instructions = "You are a helpful support assistant. Answer user questions clearly and concisely."

  tools = ["code_interpreter"]
}

# Create a thread
resource "openai_thread" "support_thread" {
  metadata = {
    user_id    = "user_12345"
    session_id = "session_67890"
    topic      = "billing_inquiry"
  }
}

# Add an initial user message to the thread
resource "openai_message" "initial_inquiry" {
  thread_id = openai_thread.support_thread.id
  role      = "user"
  content   = "I have a question about my recent invoice. The amount seems higher than usual."
  metadata = {
    source     = "web_chat"
    importance = "high"
  }
}

# Add another user message with more details
resource "openai_message" "follow_up_details" {
  thread_id = openai_thread.support_thread.id
  role      = "user"
  content   = "I was charged $59.99 but my usual subscription is $39.99. Can you help me understand why?"
  metadata = {
    source = "web_chat"
  }
}

# Example of getting a response by using a special parameter
resource "openai_message" "assistant_response" {
  thread_id = openai_thread.support_thread.id
  role      = "user"
  content   = "Please analyze the pricing difference and give me possible reasons."
  # This parameter tells the provider to send the message and wait for the assistant's response
  assistant_id = openai_assistant.support_assistant.id

  metadata = {
    source = "web_chat"
  }
}

# Create a run to process the messages with the assistant
resource "openai_run" "support_analysis" {
  thread_id    = openai_thread.support_thread.id
  assistant_id = openai_assistant.support_assistant.id

  # Optional overrides as shown in deep dive docs
  model        = "gpt-4-turbo-preview" # Using the correct model name
  instructions = "Please analyze the pricing discrepancy in detail and provide a clear explanation."

  # Control token usage and context window
  max_prompt_tokens     = 4000
  max_completion_tokens = 1000

  # Wait for completion with reasonable timeout
  wait_for_completion = true
  polling_interval    = "5s"
  timeout             = "10m"
}

output "thread_id" {
  value       = openai_thread.support_thread.id
  description = "ID of the created thread"
}

output "run_status" {
  value       = openai_run.support_analysis.status
  description = "Status of the assistant run"
}

output "assistant_response" {
  value       = openai_run.support_analysis.response_content
  description = "The assistant's response to the pricing inquiry"
}
