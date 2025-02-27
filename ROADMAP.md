# Terraform Provider OpenAI Roadmap

This document outlines the current state and future plans for the OpenAI Terraform provider.

## Currently Implemented Resources

### Core Resources
- `openai_assistant` - Manages OpenAI assistants with various capabilities
- `openai_file` - Manages file uploads for use with various OpenAI features
- `openai_fine_tuning_job` - Manages fine-tuning jobs for customizing models
- `openai_model` (Data Source) - For retrieving available models and their capabilities
- `openai_models` (Data Source) - For listing and filtering available models
- `openai_embedding` (Resource) - For managing text embedding configurations
- `openai_image_generation` (Resource) - For managing DALL-E image generation configurations
- `openai_moderation` (Resource) - For managing content moderation configurations

### Audio Resources
- `openai_speech` (Resource) - For managing text-to-speech configurations
- `openai_transcription` (Resource) - For managing audio transcription configurations
- `openai_translation` (Resource) - For managing audio translation configurations

## In Progress Resources

### Beta API Resources
- `openai_vector_store` - For managing vector storage configurations
  - Implementation started with basic CRUD operations
  - File management and batch operations support in development
- `openai_thread` - For managing conversation threads
  - Core implementation underway with support for basic operations
- `openai_thread_message` - For managing messages within threads
  - Basic structure defined, implementation in progress
- `openai_thread_run` - For managing assistant runs on threads
  - Initial implementation started

## Upcoming Features

### Core API Enhancements
- Enhanced error handling and validation for all resources
- Improved documentation with more examples and use cases
- Better support for OpenAI API versioning
- Implementation of resource import functionality

### Integration Improvements
- Better integration with other HashiCorp products
- Support for OpenAI organization management
- Enhanced monitoring and logging capabilities

## Implementation Priority

1. Complete Beta API Resources
   - Focus on stabilizing vector store implementation
   - Complete thread-related resources implementation
   - Add comprehensive testing for Beta features

2. Core API Enhancements
   - Implement remaining planned features
   - Improve existing resource functionality
   - Update documentation and examples

## Contributing

We welcome contributions to help implement any of these planned resources. When contributing:

1. Ensure comprehensive test coverage
2. Follow HashiCorp best practices for provider development
3. Document any OpenAI API version dependencies
4. Include examples in the provider documentation

## Notes

- Implementation timeline may be adjusted based on OpenAI API changes
- Beta features are subject to API changes from OpenAI
- Breaking changes in the OpenAI API may require major version bumps
- Some features may be limited by the OpenAI API's capabilities and restrictions