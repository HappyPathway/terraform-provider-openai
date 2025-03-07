terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {}

# Upload security scanning rules and compliance policies
resource "openai_file" "security_policies" {
  filename  = "security-policies.yaml"
  file_path = "${path.module}/policies/security-policies.yaml"
  purpose   = "assistants"
}

resource "openai_file" "compliance_standards" {
  filename  = "compliance-standards.json"
  file_path = "${path.module}/standards/compliance-standards.json"
  purpose   = "assistants"
}

resource "openai_file" "aws_best_practices" {
  filename  = "aws-security-best-practices.md"
  file_path = "${path.module}/guidelines/aws-security-best-practices.md"
  purpose   = "assistants"
}

# Create a compliance and security scanning assistant
resource "openai_assistant" "security_scanner" {
  name         = "Security and Compliance Scanner"
  model        = "gpt-4-turbo-preview"
  instructions = <<-EOT
    You are a specialized security and compliance assistant for infrastructure code.
    Your responsibilities include:
    1. Static analysis of infrastructure code for security issues
    2. Compliance verification against industry standards
    3. Best practices enforcement for cloud resources
    4. Security posture assessment and recommendations
    5. Generation of detailed compliance reports

    When analyzing configurations:
    - Check for security misconfigurations and vulnerabilities
    - Validate against organizational policies and compliance standards
    - Provide specific references to violated policies
    - Suggest concrete fixes with code examples
    - Consider both security and compliance implications
    - Reference relevant AWS security best practices
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
      # Files that the assistant can analyze with code interpreter
      file_ids = [
        openai_file.security_policies.id,
        openai_file.compliance_standards.id
      ]
    }
    file_search {
      # Files available for semantic search during analysis
      vector_store_ids = [
        openai_file.security_policies.id,
        openai_file.compliance_standards.id,
        openai_file.aws_best_practices.id
      ]
    }
  }

  metadata = {
    purpose      = "Security Scanning"
    department   = "Security"
    version      = "2.0"
    last_updated = "2024-04-17"
    standards    = "CIS,NIST,SOC2"
  }
}

# Create variables for infrastructure code to be scanned
variable "infrastructure_code" {
  description = "Infrastructure code to be scanned (as a string)"
  type        = string
}

# Create a thread for the security scanning session
resource "openai_thread" "security_scan" {
  metadata = {
    scan_type   = "security-and-compliance"
    scan_level  = "detailed"
    environment = "production"
  }

  # Initialize thread with tool resources
  tool_resources = {
    code_interpreter = {
      file_ids = [
        openai_file.security_policies.id,
        openai_file.compliance_standards.id
      ]
    }
    file_search = {
      vector_store_ids = [
        openai_file.security_policies.id,
        openai_file.compliance_standards.id,
        openai_file.aws_best_practices.id
      ]
    }
  }
}

# Initialize the scanning process with a message
resource "openai_message" "scan_request" {
  thread_id = openai_thread.security_scan.id
  role      = "user"
  content   = <<-EOT
    Please perform a comprehensive security and compliance analysis of the following infrastructure code.
    Provide a detailed report including:
    1. Security vulnerabilities and misconfigurations
    2. Compliance violations against standard frameworks (CIS, NIST, SOC2)
    3. Risk assessment for each finding
    4. Specific remediation steps with code examples
    5. References to relevant security policies and best practices

    Infrastructure Code:
    ${var.infrastructure_code}
  EOT

  # Attach relevant files for the analysis
  attachments = [
    {
      file_id = openai_file.security_policies.id
      tools   = ["file_search", "code_interpreter"]
    },
    {
      file_id = openai_file.compliance_standards.id
      tools   = ["file_search", "code_interpreter"]
    },
    {
      file_id = openai_file.aws_best_practices.id
      tools   = ["file_search"]
    }
  ]

  wait_for_response = true
  assistant_id      = openai_assistant.security_scanner.id
}

# Output the security analysis results
output "security_analysis" {
  description = "Detailed security and compliance analysis results"
  value       = openai_message.scan_request.response_content
}

# Example infrastructure code
locals {
  example_config = <<-EOT
    resource "aws_instance" "web" {
      ami           = "ami-12345678"
      instance_type = "t3.micro"

      root_block_device {
        encrypted = false
      }

      vpc_security_group_ids = ["sg-12345678"]

      tags = {
        Environment = "Production"
      }
    }

    resource "aws_s3_bucket" "data" {
      bucket = "my-important-data"
    }

    resource "aws_s3_bucket_public_access_block" "data" {
      bucket = aws_s3_bucket.data.id
      
      block_public_acls       = false
      block_public_policy     = false
      ignore_public_acls      = false
      restrict_public_buckets = false
    }
  EOT
}

# Example usage outputs
output "usage_example" {
  description = "Example of how to use this configuration"
  value       = <<-EOT
    # Apply this configuration with:
    terraform apply -var="infrastructure_code=${local.example_config}"

    # The assistant will analyze the code and provide:
    # 1. Security issues (e.g., unencrypted EBS, public S3 bucket)
    # 2. Compliance violations (e.g., CIS AWS Foundations)
    # 3. Remediation steps
  EOT
}
