terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

# Configure the OpenAI Provider
provider "openai" {
  # api_key = "your-api-key"               # or use OPENAI_API_KEY env var
  # organization = "your-organization-id"  # or use OPENAI_ORGANIZATION env var
  # base_url = "https://api.openai.com/v1" # for custom endpoints
  # enable_debug_logging = true            # for verbose logging
}

# Get information about GPT-4 model
data "openai_model" "gpt4" {
  model_id = "gpt-4"
}

# Upload a file to be used with the assistant
resource "openai_file" "knowledge_base" {
  filename  = "company_info.json"
  file_path = "${path.module}/data/company_info.json"
  purpose   = "assistants"
}

# Create an assistant with retrieval and code interpreter capabilities
resource "openai_assistant" "support_assistant" {
  name         = "Customer Support Assistant"
  description  = "An assistant that helps with customer inquiries using our knowledge base"
  model        = data.openai_model.gpt4.model_id
  instructions = <<-EOT
    You are a helpful customer support assistant. 
    Use the company information provided to answer questions accurately.
    If you don't know the answer, just say so - don't make up information.
    Always be polite and professional in your responses.
  EOT

  tools {
    type = "retrieval" # Enables file search
  }

  tools {
    type = "code_interpreter" # Enables code and data analysis
  }

  file_ids = [openai_file.knowledge_base.object_id]

  metadata = {
    "team"    = "customer-support"
    "version" = "1.0"
  }
}

# Create a thread for a conversation
resource "openai_thread" "customer_inquiry" {
  metadata = {
    "customer_id" = "cust_12345"
    "topic"       = "product-question"
  }
}

# Send a user message to the thread
resource "openai_message" "initial_question" {
  thread_id = openai_thread.customer_inquiry.object_id
  role      = "user"
  content {
    type = "text"
    text = "What are your company's refund policies?"
  }
  assistant_id = openai_assistant.support_assistant.object_id
  metadata = {
    "source"  = "web"
    "browser" = "chrome"
  }
}

# Create a vector embedding for semantic search
resource "openai_embedding" "search_query" {
  model = "text-embedding-ada-002"
  input = "How do I request a refund for a subscription?"
}

# Generate a completion about the company's refund policies
resource "openai_chat_completion" "policy_summary" {
  model = data.openai_model.gpt4.model_id

  messages {
    role    = "system"
    content = "You are a helpful assistant that summarizes information clearly and concisely."
  }
  messages {
    role    = "user"
    content = <<-EOT
        Based on the following refund policy, provide a brief 2-3 sentence summary:
        
        Our refund policy allows customers to request a full refund within 30 days of purchase. 
        For subscription products, we offer prorated refunds for the unused portion of the billing period.
        Digital products that have been downloaded may not be eligible for refunds unless they contain 
        technical defects. Hardware products must be returned in original packaging to qualify for a refund.
      EOT
  }

  temperature = 0.3
  max_tokens  = 150
}

# Output the assistant's response to the customer inquiry
output "assistant_response" {
  value = openai_message.initial_question.response_content
}

# Output the embedding vector (first 5 dimensions only, for readability)
output "embedding_vector_sample" {
  value = slice(openai_embedding.search_query.embedding, 0, 5)
}

# Output the summary of the refund policy
output "refund_policy_summary" {
  value = openai_chat_completion.policy_summary.response_content[0]
}
