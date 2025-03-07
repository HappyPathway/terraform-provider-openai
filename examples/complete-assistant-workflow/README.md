# Complete Assistant Workflow Example

This example demonstrates how to create a complete OpenAI assistant workflow using v2 of the Assistants API. It shows:

1. Creating a vector store for file search
2. Uploading files for code interpreter and file search
3. Creating an assistant with tool_resources
4. Creating a thread with its own tool_resources
5. Creating a message with attachments that automatically update the thread's tool_resources
6. Starting a run to process the message

## Understanding tool_resources

The example shows two ways to manage tool_resources:

1. **Assistant tool_resources**: Files and vector stores that are always available to the assistant across all threads

   ```hcl
   tool_resources = {
     code_interpreter = {
       file_ids = [openai_file.code_file.id]
     }
     file_search = {
       vector_store_ids = [openai_vector_store.knowledge_base.id]
     }
   }
   ```

2. **Thread tool_resources**: Files and vector stores that are only available for this specific thread

   ```hcl
   tool_resources = {
     code_interpreter = {
       file_ids = [openai_file.data_file.id]
     }
     file_search = {
       vector_store_ids = [openai_vector_store.knowledge_base.id]
     }
   }
   ```

3. **Message attachments**: Files attached to messages automatically get added to the thread's tool_resources
   ```hcl
   attachments = [
     {
       file_id = openai_file.data_file.id
       tools = [
         { type = "code_interpreter" },
         { type = "file_search" }
       ]
     }
   ]
   ```

## Files in this Example

- `main.tf` - The main Terraform configuration
- `data.csv` - Sample financial data for analysis
- `analysis.py` - A Python script that provides data analysis functions

## Usage

1. Configure your OpenAI provider credentials
2. Initialize Terraform: `terraform init`
3. Apply the configuration: `terraform apply`
