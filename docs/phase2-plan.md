# Phase 2 Implementation Plan: Assistant Runs and Message Management

## Overview

The current implementation allows creating assistants, threads, and messages. Phase 2 will implement the Run resource to enable actual assistant execution and response handling in a Terraform-compatible way.

## Implementation Details

### 1. openai_run Resource

The run resource will need to:
- Accept thread_id and assistant_id as required parameters
- Support optional overrides for model, instructions, and tools
- Handle the asynchronous nature of OpenAI runs in a Terraform-friendly way
- Poll the run status until completion or failure

```hcl
resource "openai_run" "example" {
  thread_id    = openai_thread.thread.id
  assistant_id = openai_assistant.assistant.id
  
  # Optional overrides
  model        = "gpt-4-turbo-preview"  # Optional
  instructions = "Override instructions" # Optional
  tools        = ["code_interpreter"]    # Optional
  
  # Terraform-specific settings
  wait_for_completion = true            # Default true
  polling_interval    = "5s"            # Default 5s
  timeout            = "10m"            # Default 10m
}
```

### 2. Run Status Handling

The run resource must handle all possible run statuses:
- queued -> in_progress -> completed (happy path)
- requires_action (function calls - future enhancement)
- expired
- cancelling/cancelled
- failed
- incomplete

### 3. Run Steps Tracking

Run steps will be exposed as computed attributes:
```hcl
output "run_steps" {
  value = openai_run.example.steps
}
```

Each step will include:
- step_id
- status
- type (message_creation or tool_calls)
- details specific to the step type

### 4. Message Integration

After a run completes:
- New messages created by the assistant will be available
- Messages should be queryable via a data source
- Consider implementing message outputs on the run resource

### 5. Error Handling

Special attention needed for:
- Timeout handling
- API rate limits and retry logic
- Proper error propagation to Terraform
- Graceful cancellation during terraform destroy

### 6. State Management

The resource will need to handle:
- Proper state storage and cleanup
- Recreate conditions (when config changes)
- Import functionality
- State drift detection

### 7. Resource Updates Needed

Based on OpenAI's assistants documentation:

#### openai_assistant Resource
- Add support for up to 128 tools per assistant
- Add support for Vision models via the model field
- Add support for `max_completion_tokens` and `max_prompt_tokens`
- Add metadata field validation
- Add truncation_strategy field for thread management:
  ```hcl
  truncation_strategy {
    type = "last_messages"
    messages_count = 10  # Optional, only for last_messages type
  }
  ```

#### openai_thread Resource
- Add support for 100,000 messages per thread limit
- Add context window management options:
  - max_prompt_tokens 
  - max_completion_tokens
  - truncation_strategy configuration
- Add thread locking logic during active runs
- Add message annotation support with file citations and file paths
- Add support for image inputs in messages with detail levels:
  ```hcl
  message {
    content = [
      {
        type = "image_url"
        image_url = {
          url = "https://example.com/image.png"
          detail = "high"  # low, high, or auto
        }
      }
    ]
  }
  ```

#### openai_file Resource
- Add size validation (max 512MB per file)
- Add token validation (max 5,000,000 tokens per file)
- Add purpose validation for "vision" type
- Add default project limit of 100GB total storage
- Add file permission management for tools:
  - code_interpreter: max 20 files
  - file_search: max 10,000 files

#### New Data Source: openai_messages
Add a data source to query messages in a thread:
```hcl
data "openai_messages" "example" {
  thread_id = openai_thread.example.id
  
  # Optional filters
  limit = 10
  order = "desc"  # asc or desc
  after = "msg_123"  # pagination
  before = "msg_789"  # pagination
}
```

### 8. Change Log Updates Needed

- Document changes in file size limits
- Document thread message limits
- Document tool resource limits
- Document image input support
- Document context window management
- Document thread locking behavior
- Document annotation support

### 9. Performance Considerations

- Implement efficient polling with exponential backoff
- Add configurable timeouts for long-running operations
- Support parallel run creation for different threads
- Implement efficient file handling for large files
- Add support for connection pooling

## Future Enhancements

1. Function Calling Support
- Handle requires_action status
- Function registration and response handling
- Timeout management for function calls

2. Advanced Options
- Streaming support via websockets
- Additional run configuration options
- Enhanced error recovery options

3. Tools Integration
- Better code interpreter integration
- File search improvements
- Custom tool support

## Implementation Order

1. Basic Run Resource
   - Create/Read/Delete operations
   - Status polling
   - Simple attribute support

2. Enhanced Features
   - Run steps tracking
   - Message integration
   - Override support

3. Error Handling & State Management
   - Comprehensive error handling
   - State management improvements
   - Import support

4. Testing & Documentation
   - Unit tests
   - Integration tests
   - Usage documentation
   - Example configurations