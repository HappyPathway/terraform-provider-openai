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

# Upload documentation files for the assistant to use
resource "openai_file" "terraform_docs" {
  filename  = "terraform-infrastructure.md"
  file_path = "${path.module}/docs/terraform-infrastructure.md"
  purpose   = "assistants"
}

resource "openai_file" "architecture_docs" {
  filename  = "architecture-overview.md"
  file_path = "${path.module}/docs/architecture-overview.md"
  purpose   = "assistants"
}

# Create an assistant specialized in infrastructure documentation
resource "openai_assistant" "infra_docs" {
  name         = "Infrastructure Documentation Assistant"
  model        = "gpt-4-turbo-preview"
  instructions = <<-EOT
    You are a specialized infrastructure documentation assistant.
    Your tasks include:
    1. Analyzing and understanding infrastructure configurations
    2. Generating and updating documentation
    3. Answering questions about infrastructure design
    4. Providing best practices and recommendations
    5. Helping maintain documentation accuracy and completeness

    When responding:
    - Be precise and technical in your explanations
    - Reference specific sections of the documentation
    - Provide concrete examples when explaining concepts
    - Highlight important security and compliance considerations
    - Suggest improvements to existing documentation
  EOT

  tools = ["code_interpreter", "file_search"]

  metadata = {
    department    = "Infrastructure"
    purpose       = "Documentation Management"
    version       = "2.0"
    last_updated  = "2024-04-17"
    documentation = "https://platform.openai.com/docs/assistants-api"
  }
}

# Create a persistent thread for documentation management
resource "openai_thread" "docs_manager" {
  metadata = {
    project     = "Infrastructure Documentation"
    environment = "production"
    type        = "documentation-management"
  }
}

# Initialize the thread with a context-setting message
resource "openai_message" "init_docs" {
  thread_id = openai_thread.docs_manager.id
  role      = "user"
  content   = <<-EOT
    Please analyze our infrastructure documentation and provide:
    1. A high-level overview of our current infrastructure
    2. List any missing or outdated documentation sections
    3. Recommend improvements for better documentation coverage
    4. Identify any security or compliance gaps in documentation
  EOT

  attachment = [
    {
      file_id = openai_file.terraform_docs.id
      tools   = ["file_search"]
    },
    {
      file_id = openai_file.architecture_docs.id
      tools   = ["file_search"]
    }
  ]

  assistant_id = openai_assistant.infra_docs.id
  metadata = {
    type = "initial-analysis"
    task = "documentation-review"
  }
}

# Output the assistant's analysis
output "documentation_analysis" {
  value = openai_message.init_docs
}
