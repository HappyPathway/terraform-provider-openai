infrastructure_code = <<EOT
resource "aws_instance" "web" {
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
EOT

# openai_api_key = "your-api-key" # Set via environment variable OPENAI_API_KEY instead