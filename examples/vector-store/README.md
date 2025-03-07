# Vector Store Example

This example demonstrates how to use OpenAI vector stores to create a searchable knowledge base for documentation using the OpenAI Terraform provider.

## Overview

This example sets up:

- A vector store for storing documentation embeddings
- Multiple documentation files that are added to the vector store
- An assistant that can search and reference the documentation
- Output configuration to monitor vector store usage and status

## Usage

1. Place your documentation files in the `files/` directory:

   - `product-documentation.pdf`
   - `api-documentation.md`
   - `user-guides.pdf`

2. Configure your OpenAI credentials:

   ```hcl
   export OPENAI_API_KEY="your-api-key"
   ```

3. Initialize and apply the Terraform configuration:
   ```bash
   terraform init
   terraform apply
   ```

## Features Demonstrated

- Vector store creation with metadata
- File management and association with vector stores
- Expiration policy configuration
- Integration with OpenAI assistants
- Usage monitoring and statistics
- Data source usage for retrieving vector store information

## Architecture

The configuration creates a knowledge base system where:

1. Documentation files are uploaded to OpenAI
2. A vector store is created to store and index the files
3. Files are added to the vector store for semantic search
4. An assistant is configured to use the vector store for answering questions
5. Monitoring is set up to track usage and status

## Outputs

The example outputs useful statistics about the vector store:

- Total number of files
- Number of successfully processed files
- Total storage usage in MB
- Current status
- Expiration timestamp

## Notes

- Make sure your documentation files exist in the `files/` directory before applying
- The vector store is configured to expire after 1 year
- The assistant is configured to use the vector store for semantic search
- File processing may take some time to complete
