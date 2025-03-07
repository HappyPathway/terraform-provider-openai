provider "openai" {}

# Upload cost analysis guidelines and reference data
resource "openai_file" "cost_guidelines" {
  filename  = "cost-optimization-guidelines.yaml"
  file_path = "${path.module}/guidelines/cost-optimization-guidelines.yaml"
  purpose   = "assistants"
}

resource "openai_file" "pricing_data" {
  filename  = "cloud-pricing.csv"
  file_path = "${path.module}/data/cloud-pricing.csv"
  purpose   = "assistants"
}

resource "openai_file" "savings_strategies" {
  filename  = "savings-strategies.md"
  file_path = "${path.module}/strategies/savings-strategies.md"
  purpose   = "assistants"
}

# Create a cost optimization assistant
resource "openai_assistant" "cost_optimizer" {
  name         = "Infrastructure Cost Optimizer"
  model        = "gpt-4-turbo-preview"
  instructions = <<-EOT
    You are a specialized infrastructure cost optimization assistant.
    Your responsibilities include:
    1. Analyzing infrastructure configurations for cost optimization opportunities
    2. Providing detailed cost breakdowns and estimates
    3. Suggesting resource rightsizing recommendations
    4. Identifying unused or underutilized resources
    5. Recommending architectural improvements for cost efficiency

    When analyzing infrastructure:
    - Compare configurations against cost optimization best practices
    - Calculate potential savings with specific numbers
    - Consider reserved instances and savings plans
    - Evaluate resource utilization patterns
    - Suggest alternative instance types or services
    - Reference pricing data for cost comparisons
  EOT

  tools = [
    {
      type = "code_interpreter"
    },
    {
      type = "file_search"
    }
  ]

  tool_resources {
    code_interpreter {
      # Files that can be used for calculations and analysis
      file_ids = [
        openai_file.pricing_data.id,
        openai_file.cost_guidelines.id
      ]
    }
    file_search {
      # Files available for semantic search
      vector_store_ids = [
        openai_file.cost_guidelines.id,
        openai_file.pricing_data.id,
        openai_file.savings_strategies.id
      ]
    }
  }

  metadata = {
    purpose      = "Cost Optimization"
    department   = "FinOps"
    version      = "2.0"
    last_updated = "2024-04-17"
    data_source  = "AWS Pricing API"
  }
}

# Variables for infrastructure configurations
variable "infrastructure_config" {
  description = "Current infrastructure configuration to analyze"
  type        = string
}

variable "budget_target" {
  description = "Target monthly budget for infrastructure"
  type        = number
}

variable "environment" {
  description = "Environment being analyzed (e.g., production, staging)"
  type        = string
  default     = "production"
}

# Create a thread for cost analysis
resource "openai_thread" "cost_analysis" {
  metadata = {
    analysis_type = "cost-optimization"
    priority      = "high"
    environment   = var.environment
  }

  # Initialize thread with access to all cost-related resources
  tool_resources {
    code_interpreter {
      file_ids = [
        openai_file.pricing_data.id,
        openai_file.cost_guidelines.id
      ]
    }
    file_search {
      vector_store_ids = [
        openai_file.cost_guidelines.id,
        openai_file.pricing_data.id,
        openai_file.savings_strategies.id
      ]
    }
  }
}

# Request cost analysis
resource "openai_message" "cost_analysis_request" {
  thread_id    = openai_thread.cost_analysis.id
  role         = "user"
  content      = <<-EOT
    Please analyze the following infrastructure configuration for cost optimization opportunities.
    Our target monthly budget is $${var.budget_target} for the ${var.environment} environment.
    
    Current Infrastructure:
    ${var.infrastructure_config}
    
    Please provide:
    1. Current estimated monthly costs
    2. Potential cost optimization opportunities:
       - Resource rightsizing recommendations
       - Reserved Instance/Savings Plan opportunities
       - Architectural improvements for cost efficiency
    3. Projected savings for each recommendation
    4. Implementation plan prioritized by impact
    5. Long-term cost optimization strategies
  EOT
  assistant_id = openai_assistant.cost_optimizer.id
}

# Outputs
output "cost_analysis" {
  description = "Detailed cost optimization analysis and recommendations"
  value       = openai_message.cost_analysis_request.response_content
}

# Example configuration
locals {
  example_infrastructure = <<-EOT
    # Production Web Infrastructure
    resource "aws_instance" "web" {
      count         = 10
      instance_type = "m5.2xlarge"
      
      ebs_block_device {
        volume_size = 1000
        volume_type = "gp2"
        encrypted   = true
      }

      tags = {
        Environment = "Production"
        Service     = "Web"
      }
    }

    resource "aws_rds_cluster" "database" {
      cluster_identifier  = "prod-db"
      engine             = "aurora-postgresql"
      engine_version     = "13.7"
      instance_class     = "db.r6g.2xlarge"
      availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
      
      backup_retention_period = 7
      storage_encrypted      = true

      serverlessv2_scaling_configuration {
        max_capacity = 64.0
        min_capacity = 16.0
      }
    }

    resource "aws_elasticache_cluster" "cache" {
      cluster_id           = "prod-cache"
      engine              = "redis"
      node_type          = "cache.r6g.xlarge"
      num_cache_nodes    = 3
      parameter_group_name = "default.redis6.x"
    }
  EOT
}

# Usage example output
output "usage_example" {
  value = <<-EOT
    # Apply this configuration with:
    terraform apply \\
      -var="infrastructure_config=${local.example_infrastructure}" \\
      -var="budget_target=25000" \\
      -var="environment=production"

    # The assistant will analyze the infrastructure and provide:
    # 1. Current cost analysis
    # 2. Optimization recommendations
    # 3. Potential savings calculations
    # 4. Implementation priorities
  EOT
}
