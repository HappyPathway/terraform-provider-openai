---
page_title: "OpenAI Provider"
description: |-
  The OpenAI provider provides resources to interact with the OpenAI API.
---

# OpenAI Provider

The OpenAI provider provides resources to interact with the [OpenAI API](https://platform.openai.com/docs/api-reference). This provider allows you to manage various OpenAI resources as part of your infrastructure as code workflow.

## Example Usage

```terraform
terraform {
  required_providers {
    openai = {
      source = "openai/openai"
    }
  }
}

# Configure the OpenAI Provider
provider "openai" {
  api_key = var.openai_api_key # Or use OPENAI_API_KEY environment variable
}
```

## Authentication

The OpenAI provider offers a flexible means of providing credentials for authentication. The following methods are supported:

* Static credentials
* Environment variables

### Static Credentials

!> **Warning:** Hard-coded credentials are not recommended in any Terraform configuration and risks secret exposure. Consider using environment variables instead.

```terraform
provider "openai" {
  api_key = "your_api_key"
  organization_id = "your_organization_id" # Optional
}
```

### Environment Variables

You can provide your credentials via the `OPENAI_API_KEY` and `OPENAI_ORGANIZATION_ID` environment variables:

```terraform
provider "openai" {}
```

```sh
$ export OPENAI_API_KEY="your_api_key"
$ export OPENAI_ORGANIZATION_ID="your_organization_id" # Optional
$ terraform plan
```

## Provider Configuration

The following arguments are supported:

* `api_key` - (Required) The API key for OpenAI. It can also be sourced from the `OPENAI_API_KEY` environment variable.
* `organization_id` - (Optional) The Organization ID for OpenAI. It can also be sourced from the `OPENAI_ORGANIZATION_ID` environment variable.
* `retry_max` - (Optional) Maximum number of retries for API requests. Defaults to 3.
* `retry_delay` - (Optional) Delay between retries in seconds. Defaults to 1.
* `timeout` - (Optional) Timeout for API requests in seconds. Defaults to 30.