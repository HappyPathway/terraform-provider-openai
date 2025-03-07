# OpenAI Provider Examples

This directory contains examples of using the OpenAI Terraform Provider for various Infrastructure as Code (IaC) use cases.

## Examples Structure

- `provider/` - Basic provider configuration examples
- `assistants/`
  - `infrastructure-docs-assistant/` - AI-powered infrastructure documentation management
  - `compliance-scanner/` - Security and compliance scanning for infrastructure code
  - `cost-optimizer/` - Infrastructure cost analysis and optimization

## Use Cases

### Infrastructure Documentation Assistant

Located in `assistants/infrastructure-docs-assistant/`, this example demonstrates how to:

- Create an AI assistant specialized in infrastructure documentation
- Maintain and analyze infrastructure documentation
- Use the assistant to generate and update documentation
- Query existing documentation for specific information

### Security and Compliance Scanner

Located in `assistants/compliance-scanner/`, this example shows how to:

- Create an assistant for automated security scanning
- Define and enforce compliance policies
- Scan infrastructure code for security vulnerabilities
- Generate detailed security and compliance reports
- Get actionable remediation steps for identified issues

### Cost Optimization Assistant

Located in `assistants/cost-optimizer/`, this example illustrates:

- Creating an assistant for infrastructure cost analysis
- Analyzing resource configurations for cost optimization
- Getting specific recommendations for cost savings
- Planning infrastructure changes within budget constraints
- Receiving long-term cost optimization strategies

## Prerequisites

1. An OpenAI API key with access to the Assistants API
2. Terraform 1.0+
3. The OpenAI provider configured with your API key

## Getting Started

1. Set your OpenAI API key:

```bash
export OPENAI_API_KEY="your-api-key"
```

2. Navigate to the example you want to try:

```bash
cd examples/assistants/cost-optimizer
```

3. Initialize Terraform:

```bash
terraform init
```

4. Create the necessary files mentioned in the examples (e.g., security rules, documentation files)

5. Apply the configuration:

```bash
terraform apply
```

## Best Practices

- Keep sensitive information in Terraform variables or environment variables
- Use workspaces to manage different environments
- Regularly update your assistants' instructions and knowledge base
- Use meaningful metadata tags for better resource organization
- Consider response times when using `wait_for_response = true`

## Notes

- The Assistants API v2 is being used in these examples, which includes the latest features for file handling and tool management
- Cost considerations apply when using vector stores and running assistant queries
- Some features might require specific OpenAI model versions or capabilities
