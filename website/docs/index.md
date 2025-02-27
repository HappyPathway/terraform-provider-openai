---
page_title: "Provider: OpenAI"
subcategory: ""
description: |-
  The OpenAI provider provides resources to interact with the OpenAI API.
---

# OpenAI Provider

The OpenAI provider provides resources to interact with the [OpenAI API](https://platform.openai.com/). This provider can be used to create and manage OpenAI resources like files, assistants, and fine-tuning jobs.

## Example Usage

```terraform
terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {
  # Configuration options will be populated from environment variables
  # OPENAI_API_KEY or from provider configuration
}
```

## Authentication

The OpenAI provider offers two ways to provide credentials:

1. Static credentials
2. Environment variables

### Static Credentials

!> **Warning:** Hard-coded credentials are not recommended in any Terraform configuration and risks secret leakage should this file ever be committed to a public version control system.

```terraform
provider "openai" {
  api_key         = "your_api_key"
  organization_id = "your_organization_id" # Optional
}
```

### Environment Variables

You can provide your credentials via the `OPENAI_API_KEY` environment variable:

```bash
export OPENAI_API_KEY="your_api_key"
```

## Argument Reference

* `api_key` - (Optional) This is the OpenAI API key. It can also be sourced from the `OPENAI_API_KEY` environment variable.
* `organization_id` - (Optional) The OpenAI organization ID to use for API requests.
* `retry_max` - (Optional) Maximum number of retries for API requests. Defaults to 3.
* `retry_delay` - (Optional) Delay between retries in seconds. Defaults to 1.
* `timeout` - (Optional) Timeout for API requests in seconds. Defaults to 30.