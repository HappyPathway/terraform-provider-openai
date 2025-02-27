# Terraform Provider OpenAI Roadmap

This document outlines the current state and future plans for the OpenAI Terraform provider.

## Currently Implemented Resources

### Core Resources
- `openai_assistant` - Manages OpenAI assistants with various capabilities
- `openai_file` - Manages file uploads for use with various OpenAI features
- `openai_fine_tuning_job` - Manages fine-tuning jobs for customizing models
- `openai_model` (Data Source) - For retrieving available models and their capabilities
- `openai_embedding` (Resource) - For managing text embedding configurations
- `openai_image_generation` (Resource) - For managing DALL-E image generation configurations
- `openai_moderation` (Resource) - For managing content moderation configurations

### Audio Resources
- `openai_speech` (Resource) - For managing text-to-speech configurations
- `openai_transcription` (Resource) - For managing audio transcription configurations
- `openai_translation` (Resource) - For managing audio translation configurations

## Planned Resources

### Core API Resources
- `openai_deployment` (Resource) - For managing model deployments and configurations

### Beta API Resources
- `openai_thread` - For managing conversation threads
- `openai_thread_message` - For managing messages within threads
- `openai_thread_run` - For managing assistant runs on threads
- `openai_vector_store` - For managing vector storage configurations

## Implementation Priority

1. Beta API Resources
   - Implementation will follow after OpenAI stabilizes these APIs
   - May require breaking changes as APIs evolve

## Contributing

We welcome contributions to help implement any of these planned resources. When contributing:

1. Ensure comprehensive test coverage
2. Follow HashiCorp best practices for provider development
3. Document any OpenAI API version dependencies
4. Include examples in the provider documentation

## Notes

- Some features may be limited by the OpenAI API's capabilities and restrictions
- Implementation timeline may be adjusted based on OpenAI API changes and community needs
- Breaking changes in the OpenAI API may require major version bumps