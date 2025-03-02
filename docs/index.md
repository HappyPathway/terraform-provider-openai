---
page_title: "OpenAI Provider"
subcategory: ""
description: |-
  The OpenAI provider allows you to interact with the OpenAI API.
---

# OpenAI Provider

The OpenAI provider allows you to interact with [OpenAI's API](https://platform.openai.com/) services. With this provider, you can create and manage resources like chat completions, file uploads, fine-tuning jobs, and use the assistants API.

## Example Usage

```terraform
terraform {
  required_providers {
    openai = {
      source = "darnold/openai"
      version = "~> 0.1"
    }
  }
}

provider "openai" {
  api_key = var.openai_api_key # or use OPENAI_API_KEY env var
}

# Use the provider to create resources
resource "openai_chat_completion" "hello_world" {
  model = "gpt-4"

  messages = [
    {
      role    = "user"
      content = "Say hello to Terraform!"
    }
  ]
}

output "response" {
  value = openai_chat_completion.hello_world.response_content[0]
}
```

## Authentication

The provider supports the following authentication methods:

1. Static credentials in the provider block
2. Environment variables

### Static Credentials

```terraform
provider "openai" {
  api_key      = "your-api-key"
  organization = "your-organization-id" # Optional
  base_url     = "https://api.openai.com/v1" # Optional
}
```

### Environment Variables

```bash
export OPENAI_API_KEY="your-api-key"
export OPENAI_ORGANIZATION="your-organization-id" # Optional
export OPENAI_BASE_URL="https://api.openai.com/v1" # Optional
```

## Schema

### Provider Configuration

- **api_key** (String, Optional) - OpenAI API key. Can also be specified with the `OPENAI_API_KEY` environment variable.
- **organization** (String, Optional) - OpenAI Organization ID. Can also be specified with the `OPENAI_ORGANIZATION` environment variable.
- **base_url** (String, Optional) - OpenAI Base URL. Can also be specified with the `OPENAI_BASE_URL` environment variable.
- **enable_debug_logging** (Boolean, Optional) - Enable debug logging. Defaults to false.
