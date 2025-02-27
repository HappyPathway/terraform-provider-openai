---
page_title: "openai_translation Resource - terraform-provider-openai"
description: |-
  Translates audio into English text using OpenAI's Whisper model.
---

# openai_translation (Resource)

This resource allows you to translate audio in any language into English text using OpenAI's Whisper model. It accepts audio input in various formats and can provide the translated text in multiple output formats.

~> **Note:** This resource creates a new translation each time it is created and cannot be updated. Changes to any of the arguments will result in a new translation being generated.

## Example Usage

```terraform
# Basic translation
resource "openai_translation" "basic" {
  audio_content = filebase64("audio/spanish_speech.mp3")
}

# Advanced translation with custom settings
resource "openai_translation" "custom" {
  audio_content    = filebase64("audio/mandarin_meeting.wav")
  model           = "whisper-1"
  prompt          = "This is a business meeting discussing quarterly results."
  response_format = "verbose_json"
  temperature     = 0.3
}

# Output the translated text
output "translation" {
  value = openai_translation.custom.text
}
```

## Argument Reference

The following arguments are supported:

* `audio_content` - (Required, Forces new resource) The base64-encoded audio file content to translate. Supported formats: flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.
* `model` - (Optional, Forces new resource) ID of the model to use. Currently only `whisper-1` is supported. Defaults to `whisper-1`.
* `prompt` - (Optional, Forces new resource) Optional text to guide the model's style or continue a previous audio segment. Should be in English.
* `response_format` - (Optional, Forces new resource) The format of the translation output. Must be one of: `json`, `text`, `srt`, `verbose_json`, or `vtt`. Defaults to `text`.
* `temperature` - (Optional, Forces new resource) The sampling temperature, between 0 and 1. Higher values make the output more random, lower values make it more focused and deterministic. Defaults to 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique identifier for this translation resource.
* `text` - The translated text output in English.

## Import

This resource cannot be imported as it represents a one-time translation.