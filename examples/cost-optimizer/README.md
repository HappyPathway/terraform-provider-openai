# Cost Optimization Assistant Example

This example demonstrates how to use OpenAI's assistant API to analyze and optimize cloud infrastructure costs. The assistant examines Terraform configurations and provides detailed cost optimization recommendations.

## Features

- Infrastructure cost analysis
- Resource rightsizing recommendations
- Reserved instance/savings plan suggestions
- Architecture optimization tips
- Cost-benefit analysis of changes

## Usage

1. Configure your OpenAI API key:
   ```bash
   export OPENAI_API_KEY="your-api-key"
   ```

2. (Optional) Edit `terraform.tfvars` to modify the infrastructure code you want to analyze.

3. Initialize and apply:
   ```bash
   terraform init
   terraform apply
   ```

## How It Works

1. Creates an AI assistant specialized in cloud cost optimization
2. Analyzes provided infrastructure code through a thread/message
3. Runs the analysis with appropriate token limits
4. Outputs detailed cost optimization recommendations

## Example Response

The assistant will analyze the infrastructure and provide:
- Identification of oversized resources
- Specific cost-saving recommendations
- Reserved instance/savings plan opportunities
- Architectural improvements for cost efficiency
- Estimated cost savings

## Notes

- The assistant uses GPT-4 Turbo for optimal analysis
- Token limits are set high to handle large infrastructure configurations
- Response times may vary based on infrastructure complexity