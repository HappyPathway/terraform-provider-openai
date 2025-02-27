# OpenAI Provider Examples

This directory contains examples demonstrating the use of the OpenAI Terraform Provider.

## Prerequisites

1. Export your OpenAI API key:
```bash
export OPENAI_API_KEY="your-api-key"
```

2. (Optional) If using an organization, set the organization ID:
```bash
export OPENAI_ORGANIZATION_ID="your-org-id"
```

## Examples

- **provider-configuration/** - Basic provider configuration example
- **completion/** - Text completion examples using various models
- **assistant/** - Creating and configuring OpenAI assistants with tools and knowledge bases
- **fine-tuning/** - Fine-tuning models with custom training data

## Running the Examples

Each directory contains a standalone Terraform configuration. To run an example:

1. Navigate to the example directory
2. Initialize Terraform: `terraform init`
3. Apply the configuration: `terraform apply`

## Notes

- Some operations (like fine-tuning) can take significant time to complete
- Be aware of your API usage and costs when running these examples
- Make sure to handle your API keys securely and never commit them to version control