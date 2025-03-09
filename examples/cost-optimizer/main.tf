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

# Create an assistant specialized in cloud cost optimization
resource "openai_assistant" "cost_optimizer" {
  name         = "Cloud Cost Optimization Assistant"
  model        = "gpt-4-turbo-preview"
  instructions = <<-EOT
    You are a specialized cloud cost optimization assistant.
    When analyzing infrastructure:
    1. Identify underutilized or oversized resources
    2. Suggest cost-effective alternatives
    3. Calculate potential cost savings
    4. Consider reserved instances and savings plans
    5. Recommend architectural improvements for cost efficiency
    6. Use specific pricing data when available
    7. Consider performance impact of recommendations
  EOT

  tools = ["code_interpreter"]
}

# Create a thread for cost analysis
resource "openai_thread" "cost_analysis" {
  metadata = {
    environment = "production"
    project     = "cloud-cost-optimization"
  }
}

# Add initial message with infrastructure details
resource "openai_message" "analyze_costs" {
  thread_id = openai_thread.cost_analysis.id
  role      = "user"
  content   = var.infrastructure_code

  metadata = {
    resource_count = "25"
    cloud_provider = "aws"
    region        = "us-west-2"
  }
}

# Run the cost analysis
resource "openai_run" "cost_analysis" {
  assistant_id = openai_assistant.cost_optimizer.id
  thread_id    = openai_thread.cost_analysis.id

  # Set high token limits for large infrastructure analysis
  max_prompt_tokens     = 4000
  max_completion_tokens = 4000
}

# Output the analysis results
output "cost_analysis" {
  value = {
    status           = openai_run.cost_analysis.status
    started_at       = openai_run.cost_analysis.started_at
    completed_at     = openai_run.cost_analysis.completed_at
    response_content = openai_run.cost_analysis.response_content
  }
}

output "usage_example" {
  value = <<EOT
# Apply this configuration with:
terraform apply -var='infrastructure_code=resource "aws_instance" "web" {
  instance_type = "t3.2xlarge"
  
  ebs_block_device {
    volume_size = 1000
    volume_type = "gp3"
  }
  
  tags = {
    Environment = "Production"
  }
}

resource "aws_rds_instance" "db" {
  instance_class = "db.r5.2xlarge"
  multi_az      = true
  allocated_storage = 500
}
'

# The assistant will analyze the infrastructure and provide:
# 1. Cost breakdown of current resources
# 2. Potential cost savings with specific recommendations
# 3. Alternative configurations 
# 4. Reserved instance recommendations
EOT
}