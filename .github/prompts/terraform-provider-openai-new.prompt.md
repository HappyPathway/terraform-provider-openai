# Terraform OpenAI Provider Prompt

## Objective

Develop a **Terraform provider for OpenAI** using the **Terraform Plugin Framework**, following the structure of AWS and Google providers. This provider should use the [`sashabaranov/go-openai`](https://github.com/sashabaranov/go-openai) library and support key OpenAI functionalities.

---

## Operations And Development

After making changes to Go files, run make build and make install before testing with Terraform commands.

## Documentation Sources

Use the following documentation sources to guide implementation:

- **Terraform AWS Provider**
  - Data Sources: [`website/docs/d`](https://github.com/hashicorp/terraform-provider-aws/tree/main/website/docs/d)
  - Resources: [`website/docs/r`](https://github.com/hashicorp/terraform-provider-aws/tree/main/website/docs/r)
- **Terraform Google Provider**
  - Data Sources: [`website/docs/d`](https://github.com/hashicorp/terraform-provider-google/tree/main/website/docs/d)
  - Resources: [`website/docs/r`](https://github.com/hashicorp/terraform-provider-google/tree/main/website/docs/r)
- **OpenAI Go SDK**: [`sashabaranov/go-openai`](https://github.com/sashabaranov/go-openai.git)
- **Terraform Plugin Framework**: [`terraform-plugin-framework`](https://github.com/hashicorp/terraform-plugin-framework.git)

---

## Core Requirements

The provider should be built with **Terraform Plugin Framework** and support:

- **Provider Configuration**

  - Accept an **API Key** (via environment variables or provider block).
  - Use the **sashabaranov/go-openai** client to interact with OpenAI APIs.

- **Resources**

  - `openai_chat_completion` – Generate chat completions (GPT models).
  - `openai_embedding` – Create text embeddings.
  - `openai_file` – Upload files for fine-tuning and retrieval.
  - `openai_assistant` – Manage OpenAI Assistants.
  - `openai_thread` – Manage conversation threads.
  - `openai_message` – Send and retrieve messages in threads.

- **Data Sources**
  - `openai_model` – Retrieve OpenAI model details.
  - `openai_assistant` – Retrieve OpenAI Assistants.
  - `openai_chat_completion` – Same as resource but happens durin terraform plan.

---

## Resource and Data Source Breakdown

### 1. `openai_chat_completion`

- Accepts a model (`gpt-4`, `gpt-3.5-turbo`) and a list of messages.
- Calls OpenAI’s **Chat Completion API**.
- Returns the generated response.
- Stateless (does not need persistent storage).

### 2. `openai_embedding`

- Accepts input text and a model (e.g., `text-embedding-ada-002`).
- Calls OpenAI’s **Embeddings API**.
- Returns a vector representation of the input text.

### 3. `openai_file`

- Uploads a file to OpenAI for fine-tuning or retrieval.
- Supports `purpose` (`fine-tune`, `assistants`).
- Stores OpenAI’s assigned file ID.

### 4. `openai_fine_tune`

- Creates a fine-tuning job from an uploaded dataset.
- Tracks the fine-tuning job’s status (`queued`, `running`, `succeeded`, `failed`).
- Stores details of the resulting fine-tuned model.

### 5. `openai_assistant`

- Represents an OpenAI Assistant.
- Stores properties like `instructions`, `model`, and `tools`.
- Supports updating assistant behavior.

### 6. `openai_thread`

- Represents a conversation thread.
- Stores messages and the conversation state.
- Allows retrieving and modifying thread history.

### 7. `openai_message`

- Sends a message within a thread.
- Fetches the assistant’s response.
- Supports **both user and assistant messages**.

### 8. `openai_model` (Data Source)

- Fetches available OpenAI models.
- Allows filtering by model type (e.g., GPT, embeddings, fine-tuning).

---

## Technical Considerations

- **Use Terraform Plugin Framework**

  - Implement provider and resources using the official [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).
  - Follow AWS and Google providers' structure for **resource/schema design**.

- **Authentication**

  - Accept API Key as a provider argument.
  - Support Terraform environment variables (`TF_VAR_openai_api_key`).

- **Rate Limiting & Retries**

  - Handle OpenAI’s API rate limits (429 errors).
  - Implement retries with exponential backoff.

- **State Management**

  - Store persistent resources (`assistant`, `thread`, `message`).
  - Make **chat completions stateless** (i.e., use `plan-modify` without storing state).

- **Logging & Debugging**

  - Use `log.Printf` for internal debugging.
  - Allow enabling debug logs via provider config.

- **Testing**
  - Implement **unit tests** using Go testing framework.
  - Write **Terraform acceptance tests** to validate provider behavior.

---

## Next Steps

1. **Scaffold the provider** using Terraform Plugin Framework.
2. **Implement core resources** (`chat_completion`, `embedding`, `assistant`, `thread`, `message`).
3. **Write documentation** following AWS & Google providers’ format.
4. **Publish to Terraform Registry** after testing and validation.

---

This prompt ensures:
✅ **Terraform Plugin Framework compatibility**  
✅ **AWS/Google provider alignment**  
✅ **Comprehensive OpenAI API coverage**
