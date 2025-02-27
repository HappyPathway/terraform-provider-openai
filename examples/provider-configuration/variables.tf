variable "openai_api_key" {
  description = "OpenAI API Key"
  type        = string
  sensitive   = true
}

variable "organization_id" {
  description = "OpenAI Organization ID (optional)"
  type        = string
  default     = null
}