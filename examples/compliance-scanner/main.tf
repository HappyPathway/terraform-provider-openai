terraform {
  required_providers {
    openai = {
      source = "happypathway/openai"
    }
  }
}

provider "openai" {}

# Upload inline security policies
resource "openai_file" "security_policies" {
  filename = "security-policies.json"
  content = jsonencode({
    policies = {
      encryption = {
        required    = true
        algorithms  = ["AES-256", "TLS-1.2"]
        description = "All data must be encrypted at rest and in transit"
      }
      network_access = {
        public_access = false
        vpn_required  = true
        description   = "Resources must not be publicly accessible"
      }
      authentication = {
        mfa_required = true
        password_policy = {
          min_length            = 12
          require_special_chars = true
        }
        description = "Strong authentication controls must be enabled"
      }
    }
  })
  purpose = "assistants"
}

# Upload inline compliance standards
resource "openai_file" "compliance_standards" {
  filename = "compliance-standards.json"
  content = jsonencode({
    standards = {
      cis = {
        version = "1.4.0"
        controls = [
          {
            id          = "1.1"
            title       = "Maintain current contact details"
            description = "Ensure contact email and phone are current for security notifications"
          },
          {
            id          = "1.2"
            title       = "Multi factor authentication"
            description = "Enable MFA for all human users in the AWS Account"
          }
        ]
      }
      nist = {
        framework = "CSF"
        version   = "1.1"
        functions = ["Identify", "Protect", "Detect", "Respond", "Recover"]
      }
      soc2 = {
        principles = ["Security", "Availability", "Processing Integrity"]
        controls = {
          CC1 = "Control Environment"
          CC2 = "Communication and Information"
          CC3 = "Risk Assessment"
        }
      }
    }
  })
  purpose = "assistants"
}

# Upload AWS security guidelines
resource "openai_file" "aws_guidelines" {
  filename  = "aws-security-best-practices.md"
  file_path = "${path.module}/aws-security-best-practices.md"
  purpose   = "assistants"
}

# Create a compliance and security scanning assistant
resource "openai_assistant" "security_scanner" {
  name         = "Security and Compliance Scanner"
  model        = "gpt-4-turbo-preview"
  instructions = file("${path.module}/assistant_instructions.md")

  tools = ["code_interpreter", "file_search"]

  tool_resources {
    code_interpreter {
      file_ids = [
        openai_file.security_policies.id,
        openai_file.compliance_standards.id,
        openai_file.aws_guidelines.id
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

# Create a thread for the security scanning session
resource "openai_thread" "security_scan" {
  metadata = {
    scan_type   = "security-and-compliance"
    scan_level  = "detailed"
    environment = "production"
  }
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

  # Configure detailed instructions for the security scanner
  scanner_instructions = <<-EOT
    EOT
}

# Initialize the scanning process with a message
resource "openai_message" "scan_request" {
  thread_id = openai_thread.security_scan.id
  role      = "user"
  content   = <<-EOT
Please perform a comprehensive security and compliance analysis of the following infrastructure code:

```hcl
${local.example_config}
```
  EOT
}

# Run the assistant on the thread to analyze the code
resource "openai_run" "security_analysis" {
  thread_id    = openai_thread.security_scan.id
  assistant_id = openai_assistant.security_scanner.id

  # Use the detailed instructions from locals
  instructions = local.scanner_instructions

  # Wait for completion with reasonable timeout
  wait_for_completion = true
  polling_interval    = "5s"
  timeout             = "10m"
}

# Output the security analysis results
output "security_analysis" {
  description = "Detailed security and compliance analysis results"
  value = {
    status             = openai_run.security_analysis.status
    started_at         = openai_run.security_analysis.started_at
    completed_at       = openai_run.security_analysis.completed_at
    last_error         = openai_run.security_analysis.last_error
    response_content   = openai_run.security_analysis.response_content
    incomplete_details = openai_run.security_analysis.incomplete_details
  }
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
