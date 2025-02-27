---
page_title: "OpenAI Provider"
subcategory: ""
description: |-
  The OpenAI provider provides resources to interact with the OpenAI API.
---

# OpenAI Provider

The OpenAI provider allows you to interact with the [OpenAI API](https://platform.openai.com/docs/api-reference) from your Terraform configurations. It provides resources to manage fine-tuning jobs, assistants, files, and generate content using OpenAI's models.

## Example Usage

```hcl
terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

# Configure the OpenAI Provider
provider "openai" {
  api_key = var.openai_api_key # Or use OPENAI_API_KEY environment variable
}
```

## Authentication

The OpenAI provider supports authentication via API key. You can provide the API key in several ways:

1. Via the provider configuration
2. Via the `OPENAI_API_KEY` environment variable
3. Via a credentials configuration file

## Provider Configuration

### Required

- `api_key` (String) - OpenAI API Key. Can also be provided via the `OPENAI_API_KEY` environment variable.

### Optional

- `organization_id` (String) - OpenAI Organization ID. Can also be provided via the `OPENAI_ORGANIZATION_ID` environment variable.
- `retry_max` (Number) - Maximum number of retries for API requests. Defaults to 3.
- `retry_delay` (Number) - Delay between retries in seconds. Defaults to 1.
- `timeout` (Number) - Timeout for API requests in seconds. Defaults to 30.