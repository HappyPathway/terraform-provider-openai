# Terraform OpenAI Provider Development Guide

This repository contains a Terraform provider for OpenAI services. Here's everything you need to know to get started and contribute effectively.

## Project Overview

This is a Terraform provider that integrates with OpenAI's API services, allowing infrastructure-as-code management of OpenAI resources. The provider is built using:
- Go 1.22.4 (required)
- Terraform Plugin SDK v2
- OpenAI Go API Client Library

## Development Environment

### Docker Development Container
The project includes a development container setup with:
- Go 1.22.4 via goenv
- Terraform via tfenv
- ZSH with Oh My ZSH
- Git and common development tools

To start development:
1. Build and start the container: `docker compose up -d`
2. Attach to the container: `docker compose exec dev zsh`

Required environment variables:
- OPENAI_API_KEY
- OPENAI_ORGANIZATION_ID (optional)
- TF_LOG (defaults to INFO)
- TF_LOG_PATH (defaults to /app/terraform.log)

### Local Provider Setup
The provider is configured for local development with:
- Source: "happypathway/openai"
- Local development override in .terraformrc
- Provider configuration via environment variables

## Development Configuration

### Key Project Files
- `.terraformrc` - Local provider override configuration
- `docker-compose.yml` - Development container setup
- `Dockerfile` - Container configuration with Go 1.22.4 and Terraform 1.7.4
- `test.tf` - Example configuration for testing resources

### Environment Setup
Required environment variables for development:
```
OPENAI_API_KEY=your-api-key
OPENAI_ORGANIZATION_ID=your-org-id  # Optional
TF_LOG=INFO  # Default
TF_LOG_PATH=/app/terraform.log  # Default
```

### Development Container Features
- Go 1.22.4 via goenv
- Terraform 1.7.4 via tfenv
- ZSH with Oh My ZSH
- Development tools (git, make, vim, etc.)
- VSCode integration with devcontainer support

### Local Provider Configuration
Provider source is configured as "happypathway/openai" with local development override in .terraformrc.

### Current Example Resources
The test.tf file includes examples of:
- File resources (training_data, knowledge_base)
- Vector store configuration
- Thread resources with tool integration
- Assistant configuration with multiple tool types
- Message resources with different content types

### VSCode Development Container
To use the development container in VSCode:
1. Install VS Code Remote Development extension
2. Command Palette -> "Remote-Containers: Open Folder in Container"
3. VS Code will build and start the container
4. Container includes all necessary development tools and extensions

## Project Structure

### Key Directories
- `/openai/` - Core provider implementation
  - provider.go - Main provider configuration
  - resource_*.go - Resource implementations
  - data_source_*.go - Data source implementations
- `/examples/` - Example configurations for different features
- `/testdata/` - Test fixtures and data

### Documentation Sources
Refer to these directories for implementation guidance:
- `/docs/openai-go/` - OpenAI Go API documentation
- `/docs/terraform-plugin-sdk/website/docs/plugin/sdkv2/` - Terraform Plugin SDK v2 documentation

## Development Guidelines

1. **Documentation**: 
   - Follow .md extension convention
   - Reference official documentation in openai-go and sdkv2 directories
   - Don't rely on assumptions about provider setup

2. **Resource Implementation**:
   - Follow Terraform Plugin SDK v2 patterns
   - Include proper schema validation
   - Implement CRUD operations completely
   - Add acceptance tests

3. **Code Organization**:
   - Keep resource implementations in separate files
   - Use consistent naming conventions
   - Follow Go best practices and SDK patterns

4. **Testing**:
   - Write both unit and acceptance tests
   - Use test fixtures from testdata directory
   - Test with both regular and Beta API features

5. **API Integration**:
   - Use the official OpenAI Go client library
   - Handle API versioning appropriately
   - Support both regular and Azure OpenAI endpoints

## Common Tasks

### Adding a New Resource
1. Create new resource file in /openai/
2. Define schema and CRUD functions
3. Register in provider.go
4. Add acceptance tests
5. Add example configuration
6. Update documentation

### Testing Changes
1. Build the provider: `make build`
2. Run tests: `make test`
3. Run acceptance tests: `make testacc`
4. Test with example configurations

### Debugging
- Use TF_LOG environment variable for logging
- Check terraform.log for detailed logs
- Use the debugging facilities in the SDK

## Key Implementation Details

### Provider Configuration
- Supports API key and organization ID
- Configurable retry and timeout settings
- Beta API features support

### Resource Types
- Files
- Models
- Fine-tuning jobs
- Assistants
- Threads
- Messages
- Vector stores

## Testing Guidelines

### Test Structure
1. **Unit Tests**
   - Located in `_test.go` files alongside the implementation
   - Use Go's standard testing framework
   - Mock OpenAI API responses using testutil package
   - Example: `openai/data_source_openai_fine_tuned_model_test.go`

2. **Acceptance Tests**
   - Test actual integration with OpenAI API
   - Use `resource.Test` framework from Terraform SDK
   - Implement PreCheck functions to validate environment
   - Example test pattern:
     ```go
     resource.Test(t, resource.TestCase{
       PreCheck:          func() { testAccPreCheck(t) },
       ProviderFactories: providerFactories,
       Steps: []resource.TestStep{
         {
           Config: testAccConfig,
           Check:  resource.ComposeTestCheckFunc(...),
         },
       },
     })
     ```

3. **Test Data**
   - Store test fixtures in /testdata directory
   - Use meaningful model IDs and file names in tests
   - Include test configurations for all resource types
   - Validate resource attributes and state

### Running Tests
1. **Environment Setup**
   ```
   export OPENAI_API_KEY="your-api-key"
   export OPENAI_ORGANIZATION_ID="your-org-id" # optional
   export TF_ACC=1  # for acceptance tests
   ```

2. **Test Commands**
   - Unit Tests: `make test`
   - Acceptance Tests: `make testacc`
   - Specific Test: `go test -v ./openai -run TestAccDataSource`

3. **Test Debugging**
   - Use TF_LOG=DEBUG for detailed logs
   - Check terraform.log for API interactions
   - Examine test state files in terraform.tfstate

### Test Types to Implement

1. **Basic Resource Tests**
   - Creation and deletion
   - Required fields validation
   - Optional fields handling
   - State consistency

2. **Data Source Tests**
   - Data retrieval accuracy
   - Field mapping verification
   - Error handling
   - Optional parameters

3. **Integration Tests**
   - Resource dependencies
   - Cross-resource references
   - State management
   - Update scenarios

4. **Error Cases**
   - Invalid configurations
   - API error responses
   - Resource constraints
   - Timeout handling

## Contributing Guidelines

### Pull Request Process
1. Create feature branch
2. Add tests for new functionality
3. Update documentation
4. Submit PR with detailed description

### Code Standards
1. Follow Go best practices
2. Use consistent error handling
3. Implement proper logging
4. Add doc comments for exported functions

### Documentation Requirements
1. Update provider docs
2. Include example configurations
3. Document resource attributes
4. Provide usage examples

## Troubleshooting

### Common Issues
1. API Authentication
   - Verify environment variables
   - Check API key permissions
   - Validate organization ID

2. Resource State
   - Examine terraform.tfstate
   - Check for drift
   - Verify ID formats

3. Testing Problems
   - Mock API responses
   - Test environment setup
   - Resource cleanup

### Debugging Tools
1. TF_LOG levels
2. API response inspection
3. State file examination
4. Error tracking

Remember to always validate changes with both unit and acceptance tests, and ensure documentation is updated accordingly.

