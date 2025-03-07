# OpenAI Assistants API Migration Plan: v1 to v2

## Overview

This document outlines the step-by-step plan for migrating from OpenAI Assistants API v1 to v2. The migration deadline is December 18, 2024, after which v1 will no longer be accessible.

## Prerequisites

- Review the current API integration points in your codebase
- Identify all uses of AssistantFile and MessageFile objects
- Ensure you have the latest version of the OpenAI SDK installed

## Migration Steps

### 1. Update SDK Dependencies

If using Python:

```bash
pip install --upgrade openai
```

If using Node.js:

```bash
npm install openai@latest
```

### 2. File Management Changes

#### 2.1 Assistant Files Migration

- Replace direct file_ids on Assistants with tool_resources
- For code_interpreter files:
  - Move files to tool_resources.code_interpreter.file_ids
- For retrieval files:
  - Migrate to vector_stores under tool_resources.file_search
  - Ensure vector store ingestion is complete before creating runs

#### 2.2 Message Files Migration

- Replace file_ids with attachments in Messages
- Update file attachment logic to use the new attachments format
- Verify files are properly associated with Thread tool_resources

### 3. Tool Updates

- Rename all instances of 'retrieval' tool to 'file_search'
- Update tool configurations in Assistant creation/updates
- Modify any tool-specific logic to handle the new structure

### 4. API Version Header Updates

For direct API calls:

```
Header: OpenAI-Beta: assistants=v2
```

For SDK usage:

- Remove any explicit v1 version headers
- Update to use default v2 endpoints

### 5. Testing Strategy

1. Create test Assistants with both code_interpreter and file_search tools
2. Verify file attachments work in both Assistant and Thread contexts
3. Test vector store creation and file search functionality
4. Validate existing functionality works with the new API version

### 6. Monitoring and Rollout

1. Monitor vector store usage (billing starts 2025)
2. Clean up unused v1 resources
3. Document any v2-specific configurations
4. Keep track of vector stores created before April 17, 2024 (free tier)

## Important Considerations

- Vector stores created before April 17, 2024, are free until end of 2024
- Unused vector stores from v1 will be deleted if not used in a Run
- File deletions in v2 don't propagate to v1
- Consider cleaning up v1 files using v1 endpoints or direct file deletion

## Timeline

1. Phase 1: Development and Testing (2-3 weeks)

   - Update dependencies
   - Implement code changes
   - Create test suite

2. Phase 2: Validation (1 week)

   - Run integration tests
   - Verify all features work as expected
   - Document any issues or concerns

3. Phase 3: Production Rollout (1-2 weeks)
   - Gradual rollout to production
   - Monitor for any issues
   - Complete migration by December 18, 2024

## Rollback Plan

1. Keep v1 code in a separate branch
2. Maintain ability to switch between v1 and v2 headers
3. Document all v2-specific changes for quick reversion if needed
