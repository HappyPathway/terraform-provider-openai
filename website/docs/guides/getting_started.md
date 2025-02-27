---
page_title: "Getting Started with OpenAI Provider"
subcategory: ""
description: |-
  Getting started with the OpenAI Terraform provider.
---

# Getting Started with OpenAI Provider

This guide will help you get started with the OpenAI provider. We'll cover basic setup and walk through some common use cases.

## Before you begin

1. Sign up for an OpenAI account at https://platform.openai.com/
2. Create an API key from the OpenAI dashboard
3. Install Terraform

## Configuration

First, create a new directory for your Terraform configuration and create a `main.tf` file:

```terraform
terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {
  # API key can be set via OPENAI_API_KEY environment variable
}
```

Set your API key as an environment variable:

```bash
export OPENAI_API_KEY="your-api-key"
```

## Example: Creating an Assistant

Here's a complete example of creating an AI assistant with file-based knowledge:

```terraform
# Upload a knowledge base file
resource "openai_file" "knowledge_base" {
  file    = "${path.module}/knowledge_base.txt"
  purpose = "assistants"
}

# Create an assistant that uses the knowledge base
resource "openai_assistant" "example" {
  name         = "Documentation Helper"
  description  = "An assistant that helps with documentation"
  model        = "gpt-4-1106-preview"
  instructions = "You are a helpful assistant that specializes in technical documentation."

  tools {
    type = "retrieval"
  }

  tools {
    type = "code_interpreter"
  }

  file_ids = [openai_file.knowledge_base.id]
}
```

## Example: Fine-tuning a Model

Here's how to fine-tune a model with custom training data:

```terraform
# Upload training data
resource "openai_file" "training_data" {
  file    = "${path.module}/training.jsonl"
  purpose = "fine-tune"
}

# Create a fine-tuning job
resource "openai_fine_tuning_job" "custom_model" {
  model         = "gpt-3.5-turbo"
  training_file = openai_file.training_data.id
}
```

## Example: Dynamic Content Generation

Use the content generator for creating dynamic infrastructure descriptions:

```terraform
resource "openai_content_generator" "service_description" {
  model  = "gpt-4"
  prompt = "Write a technical description for a highly available web service:"
  
  parameters {
    temperature = 0.7
    max_tokens  = 200
  }
}

resource "aws_ssm_parameter" "service_description" {
  name  = "/service/description"
  type  = "String"
  value = openai_content_generator.service_description.content
}
```

## Next Steps

- Explore the [provider documentation](../index.html) for detailed information about all available resources
- Look at the [examples directory](https://github.com/HappyPathway/terraform-provider-openai/tree/main/examples) for more complex configurations
- Join our [community](https://github.com/HappyPathway/terraform-provider-openai/discussions) for support and discussions