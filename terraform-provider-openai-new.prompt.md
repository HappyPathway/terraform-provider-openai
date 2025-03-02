# Terraform OpenAI Provider

The Terraform OpenAI Provider allows you to manage OpenAI resources and interact with the OpenAI API in your Terraform configurations.

## Provider Configuration

To configure the OpenAI provider, include the following in your Terraform configuration:

```hcl
terraform {
  required_providers {
    openai = {
      source = "registry.terraform.io/darnold/openai"
      version = "~> 1.0"
    }
  }
}

provider "openai" {
  api_key = var.openai_api_key # or set OPENAI_API_KEY environment variable
  # optional:
  # organization = var.openai_org_id # or set OPENAI_ORGANIZATION environment variable
  # base_url = "https://api.openai.com/v1" # or set OPENAI_BASE_URL environment variable
}
```

## Authentication

The provider supports API key authentication. You can provide the API key in several ways:

1. Set the `api_key` attribute in the provider configuration.
2. Set the `OPENAI_API_KEY` environment variable.

## Data Sources

### openai_model

Retrieve information about an available OpenAI model.

```hcl
data "openai_model" "gpt4" {
  id = "gpt-4"
}

output "model_info" {
  value = data.openai_model.gpt4
}
```

### openai_assistant

Retrieve information about an existing assistant.

```hcl
data "openai_assistant" "example" {
  assistant_id = "asst_abc123"
}

output "assistant_info" {
  value = data.openai_assistant.example
}
```

### openai_chat_completion

Generate a completion using the Chat API as a data source.

```hcl
data "openai_chat_completion" "example" {
  model = "gpt-4"

  messages = [
    {
      role    = "system"
      content = "You are a helpful assistant."
    },
    {
      role    = "user"
      content = "What is the capital of France?"
    }
  ]

  temperature = 0.7
}

output "response" {
  value = data.openai_chat_completion.example.response_content[0]
}
```

## Resources

### openai_chat_completion

Generate a completion using the Chat API.

```hcl
resource "openai_chat_completion" "example" {
  model = "gpt-4"

  messages = [
    {
      role    = "system"
      content = "You are a helpful assistant."
    },
    {
      role    = "user"
      content = "What is the capital of France?"
    }
  ]

  temperature = 0.7
}

output "response" {
  value = openai_chat_completion.example.response_content[0]
}
```

### openai_embedding

Generate embeddings for input text.

```hcl
resource "openai_embedding" "example" {
  model = "text-embedding-ada-002"
  input = "The quick brown fox jumps over the lazy dog"
}

output "embedding" {
  value = openai_embedding.example.embedding
}
```

### openai_file

Upload a file to OpenAI for use with other API endpoints.

```hcl
resource "openai_file" "example" {
  filename = "data.jsonl"
  filepath = "${path.module}/data.jsonl"
  purpose  = "fine-tune"
}

output "file_id" {
  value = openai_file.example.id
}
```

### openai_fine_tune

Create a fine-tuned model.

```hcl
resource "openai_file" "training_file" {
  filename = "training_data.jsonl"
  filepath = "${path.module}/training_data.jsonl"
  purpose  = "fine-tune"
}

resource "openai_fine_tune" "example" {
  training_file     = openai_file.training_file.id
  model             = "gpt-3.5-turbo"
  suffix            = "my-fine-tuned-model"
  hyperparameters = {
    n_epochs = 3
  }
}

output "fine_tune_id" {
  value = openai_fine_tune.example.id
}
```

### openai_assistant

Create an AI assistant.

```hcl
resource "openai_assistant" "example" {
  name         = "Customer Support Assistant"
  description  = "An assistant that helps with customer inquiries"
  model        = "gpt-4-turbo"
  instructions = "You are a helpful customer support assistant. Always be polite and concise."

  tools = [
    {
      type = "code_interpreter"
    }
  ]
}

output "assistant_id" {
  value = openai_assistant.example.id
}
```

### openai_thread

Create a thread for conversations.

```hcl
resource "openai_thread" "example" {
  metadata = {
    user_id = "user_123"
  }
}

output "thread_id" {
  value = openai_thread.example.id
}
```

### openai_message

Add a message to a thread.

```hcl
resource "openai_thread" "example" {
  metadata = {
    user_id = "user_123"
  }
}

resource "openai_message" "example" {
  thread_id = openai_thread.example.id
  role      = "user"
  content   = "I need help with my order #12345"
}

output "message_id" {
  value = openai_message.example.id
}
```

## Complete Assistant Workflow Example

```hcl
# Create an assistant
resource "openai_assistant" "support" {
  name         = "Customer Support Assistant"
  description  = "An assistant that helps with customer inquiries"
  model        = "gpt-4-turbo"
  instructions = "You are a helpful customer support assistant. Always be polite and concise."
}

# Create a thread
resource "openai_thread" "customer_inquiry" {
  metadata = {
    customer_id = "cust_123"
    topic       = "Order Status"
  }
}

# Add a user message
resource "openai_message" "inquiry" {
  thread_id = openai_thread.customer_inquiry.id
  role      = "user"
  content   = "I placed order #12345 yesterday, but I haven't received a shipping confirmation yet. Can you help?"
}

# Run the assistant on the thread
resource "openai_run" "support_response" {
  assistant_id = openai_assistant.support.id
  thread_id    = openai_thread.customer_inquiry.id

  # Wait for the run to complete
  wait_for_completion = true

  # Set a timeout to prevent indefinite waiting
  completion_timeout_seconds = 60
}

# Get the assistant's response
data "openai_messages" "thread_messages" {
  thread_id = openai_thread.customer_inquiry.id

  # Ensure this runs after the assistant has responded
  depends_on = [openai_run.support_response]
}

output "assistant_response" {
  # Output the latest message content from the assistant
  value = [
    for msg in data.openai_messages.thread_messages.messages :
    msg.content
    if msg.role == "assistant"
  ][0]
}
```

## Best Practices

1. **API Key Management**

   - Store your API key securely using Terraform variables or environment variables
   - Consider using a secrets management solution for production environments

2. **Resource Management**

   - Use the `openai_file` resource for uploading training data and fine-tuning
   - Use meaningful suffixes for fine-tuned models to identify them easily

3. **Error Handling**

   - The provider includes retry logic with exponential backoff for API rate limits and server errors
   - Set appropriate timeouts for long-running operations like fine-tuning

4. **Cost Management**

   - Monitor API usage to control costs
   - Use the appropriate model for your needs (smaller models cost less)

5. **Versioning**
   - Pin the provider version in your Terraform configuration
   - Test provider upgrades in a non-production environment first

## Troubleshooting

1. **API Rate Limits**

   - The provider automatically handles rate limiting with retries
   - If you're consistently hitting rate limits, consider implementing client-side throttling

2. **Authentication Issues**

   - Ensure your API key is valid and has the necessary permissions
   - Check that the organization ID is correctly specified if applicable

3. **Resource State Issues**

   - For stateless resources like completions, re-creation is normal on subsequent applies
   - For stateful resources, check the OpenAI API dashboard if resources aren't being found

4. **Long-Running Operations**
   - Fine-tuning and assistant runs can take time to complete
   - Use `wait_for_completion` and appropriate timeouts for operations that should block

## Additional Information

For detailed API documentation, visit the [OpenAI API documentation](https://platform.openai.com/docs/api-reference/introduction).

For provider-specific issues, please open an issue on the [GitHub repository](https://github.com/darnold/terraform-provider-openai).
