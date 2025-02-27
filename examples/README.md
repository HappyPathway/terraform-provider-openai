# OpenAI Provider Examples

This directory contains examples that demonstrate various use cases for the OpenAI Terraform Provider.

## Prerequisites

1. Set up your OpenAI API key as an environment variable:
```bash
export TF_VAR_openai_api_key="your-api-key"
```

2. (Optional) If you're using an organization, set the organization ID:
```bash
export TF_VAR_organization_id="your-org-id"
```

## Examples

- **provider-configuration/** - Basic provider configuration example
- **completion/** - Text completion examples using various models
- **embedding/** - Text embedding generation for vectorization
- **assistant/** - Creating and configuring OpenAI assistants with tools and knowledge bases
- **fine-tuning/** - Fine-tuning models with custom training data

## Running the Examples

Each directory contains a standalone Terraform configuration. To run an example:

1. Change to the example directory:
```bash
cd provider-configuration
```

2. Initialize Terraform:
```bash
terraform init
```

3. Apply the configuration:
```bash
terraform apply
```

## Notes

- The examples use sample data files for demonstration purposes. Replace them with your actual data files in production.
- Some operations (like fine-tuning) can take significant time to complete.
- Consider your API usage and costs when running these examples.