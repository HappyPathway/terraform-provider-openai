# OpenAI Provider Resources Implementation Guide

## Assistants API Implementation Plan

This document outlines the implementation plan for the OpenAI Assistants API resources in the Terraform provider.

### Current Status

- [x] Assistant Resource (`openai_assistant`)
- [x] Assistant Data Source (`data.openai_assistant`)
- [ ] Thread Resource (`openai_thread`)
- [ ] Message Resource (`openai_thread_message`)
- [ ] Run Resource (`openai_thread_run`)

### Implementation Plan

#### 1. Thread Resource (`openai_thread`)

The Thread resource will be implemented first as it's the foundation for messages and runs.

##### Required Components:
- Schema Definition:
  - `id` (Computed)
  - `metadata` (Optional, Map)
  - `messages` (Optional, Block List)
    - `content` (Required)
    - `role` (Required, "user" or "assistant")
    - `file_ids` (Optional, List)
    - `metadata` (Optional, Map)
  - `tool_resources` (Optional, Block)
    - `code_interpreter` (Optional, Block)
      - `file_ids` (Optional, List)
    - `file_search` (Optional, Block)
      - `vector_store_ids` (Optional, List)

##### Implementation Files:
- `openai/resource_openai_thread.go`
- `openai/resource_openai_thread_test.go`

#### 2. Message Resource (`openai_thread_message`)

Messages are created within threads and represent the conversation.

##### Required Components:
- Schema Definition:
  - `id` (Computed)
  - `thread_id` (Required)
  - `role` (Required, "user" or "assistant")
  - `content` (Required, Block List)
    - `type` (Required, "text" or "image_url")
    - `text` (Optional)
    - `image_url` (Optional, Block)
      - `url` (Required)
      - `detail` (Optional)
  - `file_ids` (Optional, List)
  - `metadata` (Optional, Map)
  - `assistant_id` (Computed)
  - `run_id` (Computed)
  - `created_at` (Computed)

##### Implementation Files:
- `openai/resource_openai_thread_message.go`
- `openai/resource_openai_thread_message_test.go`

#### 3. Run Resource (`openai_thread_run`)

Runs execute assistant actions on threads.

##### Required Components:
- Schema Definition:
  - `id` (Computed)
  - `thread_id` (Required)
  - `assistant_id` (Required)
  - `model` (Optional, overrides assistant's model)
  - `instructions` (Optional, overrides assistant's instructions)
  - `tools` (Optional, List)
    - `type` (Required)
    - `function` (Optional, Block)
      - `name` (Required)
      - `description` (Optional)
      - `parameters` (Required, JSON string)
  - `metadata` (Optional, Map)
  - `status` (Computed)
  - `started_at` (Computed)
  - `completed_at` (Computed)
  - `last_error` (Computed, Block)
    - `code` (Computed)
    - `message` (Computed)
  - `expires_at` (Computed)
  - `required_action` (Computed, Block)
    - `type` (Computed)
    - `submit_tool_outputs` (Computed, Block)
      - `tool_calls` (Computed, List)
  - `usage` (Computed, Block)
    - `completion_tokens` (Computed)
    - `prompt_tokens` (Computed)
    - `total_tokens` (Computed)

##### Implementation Files:
- `openai/resource_openai_thread_run.go`
- `openai/resource_openai_thread_run_test.go`

### Implementation Order and Dependencies

1. Thread Resource
   - Base implementation
   - CRUD operations
   - Basic validation
   - Acceptance tests

2. Message Resource
   - Depends on Thread resource
   - CRUD operations
   - Content validation
   - File attachment handling
   - Acceptance tests

3. Run Resource
   - Depends on Thread and Assistant resources
   - CRUD operations
   - Status tracking
   - Tool output handling
   - Required action handling
   - Acceptance tests

### Testing Strategy

For each resource:
1. Unit tests for validation functions
2. Acceptance tests covering:
   - Basic CRUD operations
   - Error cases
   - Update scenarios
   - Import functionality
   - Resource dependencies

### Documentation Requirements

For each resource:
1. Resource documentation in `website/docs/r/[resource_name].html.markdown`
2. Example configurations in `examples/`
3. Guide updates for complex scenarios
4. Update provider documentation to include new resources

### Additional Considerations

1. Rate Limiting
   - Implement proper retry logic
   - Handle API quotas appropriately

2. Error Handling
   - Proper error messages
   - Graceful failure handling
   - Status code handling

3. State Management
   - Proper state tracking
   - Handle incomplete or failed operations
   - Resource cleanup

4. Data Validation
   - Input validation
   - Schema validation
   - Type checking

### Future Enhancements

1. Data Sources
   - Thread data source
   - Message data source
   - Run data source

2. Complex Operations
   - Batch operations
   - Multi-step workflows
   - Custom tool implementations

This implementation plan will be updated as needed during development.