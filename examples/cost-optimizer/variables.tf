variable "infrastructure_code" {
  description = "The infrastructure code to analyze for cost optimization"
  type        = string
}

variable "openai_api_key" {
  description = "OpenAI API key"
  type        = string
  sensitive   = true
}
