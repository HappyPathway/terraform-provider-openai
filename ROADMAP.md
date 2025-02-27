# OpenAI Provider Roadmap

## Phase 1: Core Resources (Initial Release)

### Authentication & Configuration
- [x] Basic provider configuration with API key
- [ ] Organization ID support
- [x] API request retry and timeout configurations

### Resources
- [ ] `openai_completion` - Manage text completions
- [ ] `openai_chat_completion` - Manage chat completions
- [x] `openai_fine_tuning_job` - Manage fine-tuning jobs (Implementation complete, testing in progress)
- [x] `openai_assistant` - Manage OpenAI assistants
- [x] `openai_file` - Manage files for fine-tuning and assistants

### Data Sources
- [x] `openai_model` - Information about available models
- [x] `openai_models` - List of available models
- [ ] `openai_file` - Information about uploaded files

## Current Focus: Testing & Stability

### Testing Requirements
- [x] Basic acceptance test structure
- [ ] Test environment variable configuration
  - Required variables: TF_VAR_project_prompt, TF_VAR_repo_org, TF_VAR_project_name
- [ ] Test data preparation
- [ ] Mock API responses for tests that require paid features

### Immediate Next Steps
1. Update test configurations to work without GitHub Pro features
2. Implement mock responses for fine-tuning job tests
3. Document required test environment variables
4. Add validation for required test configuration

## Phase 2: Extended Resources

### Resources
- [ ] `openai_assistant_file` - Manage files attached to assistants
- [ ] `openai_thread` - Manage persistent threads for conversations
- [ ] `openai_deployment` - Manage model deployments (Azure OpenAI)
- [ ] `openai_embedding` - Generate and manage embeddings

### Data Sources
- [ ] `openai_fine_tuning_job` - Information about fine-tuning jobs
- [ ] `openai_assistant` - Information about existing assistants
- [ ] `openai_files` - List of uploaded files

## Phase 3: Advanced Features

### Resources
- [ ] `openai_fine_tuned_model` - Manage custom fine-tuned models
- [ ] `openai_assistant_tool` - Manage tools for assistants
- [ ] `openai_message` - Manage messages in threads

### Data Sources
- [ ] `openai_deployment` - Information about model deployments
- [ ] `openai_deployments` - List of model deployments
- [ ] `openai_usage` - Usage statistics and quotas

## Future Considerations

### Potential Features
- Azure OpenAI Service support
- Rate limiting and quota management
- Batch operations support
- Cost estimation and management
- Monitoring and logging integrations

### Integration Points
- [ ] Integration with Azure OpenAI Service
- [ ] Integration with AWS Bedrock
- [ ] Support for organization management
- [ ] Support for team-level access controls

## Implementation Notes

### Priority Order
1. [x] Core authentication and configuration
2. [x] Basic model and file management
3. [ ] Stabilize existing resource tests
4. [ ] Completion and Chat APIs support
5. [ ] Assistant and fine-tuning support
6. [ ] Extended features and integrations

### Development Guidelines
- Each resource/data source will include:
  - Complete documentation
  - Example configurations
  - Acceptance tests that work without paid features
  - Import support where applicable
  - Proper error handling
  - Rate limiting consideration

### Testing Strategy
- Unit tests for all resources and data sources
- Acceptance tests using minimal API features
- Mock responses for paid features
- Documentation examples as tests
- Environment variable configuration guide

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to the provider.