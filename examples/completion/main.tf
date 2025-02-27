terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

resource "openai_completion" "story_completion" {
  model       = "gpt-3.5-turbo-instruct"
  prompt      = "Once upon a time in a digital forest"
  max_tokens  = 150
  temperature = 0.7
  n           = 1  # Number of completions to generate

  # Optional parameters
  presence_penalty  = 0.6
  frequency_penalty = 0.2
  stop              = ["\n\n", "THE END"]
}

output "story" {
  value = openai_completion.story_completion.choices[0].text
}