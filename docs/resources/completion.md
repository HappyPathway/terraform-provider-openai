---
page_title: "openai_completion Resource - terraform-provider-openai"
description: |-
  Creates completions with OpenAI models.
---

# openai_completion (Resource)

Creates completions with OpenAI models. Given a prompt, the model will generate one or more predicted completions.

~> **Note:** This resource creates a completion each time it is created and cannot be updated. Changes to any of the arguments will result in a new completion being created.

## Example Usage

```terraform
resource "openai_completion" "example" {
  model     = "gpt-3.5-turbo-instruct"
  prompt    = "Write a haiku about Terraform"
  max_tokens = 50
  temperature = 0.7
  n = 1
}

output "completion_text" {
  value = openai_completion.example.choices[0].text
}
```

## Argument Reference

The following arguments are supported:

* `model` - (Required, Forces new resource) ID of the model to use for completion.
* `prompt` - (Required, Forces new resource) The prompt to generate completions for.
* `max_tokens` - (Optional, Forces new resource) The maximum number of tokens to generate.
* `temperature` - (Optional, Forces new resource) Sampling temperature to use, between 0 and 2. Higher values mean more random completions. Defaults to 1.
* `top_p` - (Optional, Forces new resource) Alternative to temperature for nucleus sampling. Defaults to 1.
* `n` - (Optional, Forces new resource) How many completion choices to generate. Defaults to 1.
* `best_of` - (Optional, Forces new resource) Number of completions to generate server-side and return the best n.
* `frequency_penalty` - (Optional, Forces new resource) Penalize new tokens based on their frequency in the text so far. Between -2.0 and 2.0.
* `presence_penalty` - (Optional, Forces new resource) Penalize new tokens based on whether they appear in the text so far. Between -2.0 and 2.0.
* `stop` - (Optional, Forces new resource) List of sequences where the API will stop generating tokens.
* `suffix` - (Optional, Forces new resource) Text to append to the completion.
* `logit_bias` - (Optional, Forces new resource) Map of token biases to use for completion.
* `user` - (Optional, Forces new resource) A unique identifier representing your end-user.
* `seed` - (Optional, Forces new resource) Integer seed for deterministic completions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique identifier for the completion.
* `object` - The object type, which is always "text_completion".
* `created` - The Unix timestamp (in seconds) of when the completion was created.
* `choices` - A list of completion choices. Each contains:
  * `text` - The generated completion text.
  * `index` - The index of this completion in the list.
  * `finish_reason` - The reason why the completion ended.
* `usage` - Information about token usage, contains:
  * `prompt_tokens` - Number of tokens in the prompt.
  * `completion_tokens` - Number of tokens in the generated completion.
  * `total_tokens` - Total number of tokens used.