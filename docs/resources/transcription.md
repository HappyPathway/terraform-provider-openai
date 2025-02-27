---
page_title: "openai_transcription Resource - terraform-provider-openai"
description: |-
  Transcribes audio into text using OpenAI's Whisper model.
---

# openai_transcription (Resource)

This resource allows you to transcribe audio into text using OpenAI's Whisper model. It accepts audio input in various formats and can provide the transcription in multiple output formats.

~> **Note:** This resource creates a new transcription each time it is created and cannot be updated. Changes to any of the arguments will result in a new transcription being generated.

## Example Usage

```terraform
# Basic transcription
resource "openai_transcription" "basic" {
  audio_content = filebase64("audio/recording.mp3")
}

# Advanced transcription with custom settings
resource "openai_transcription" "custom" {
  audio_content    = filebase64("audio/meeting.wav")
  model           = "whisper-1"
  language        = "en"
  prompt          = "This is a business meeting discussion about project timelines."
  response_format = "verbose_json"
  temperature     = 0.3
}

# Output the transcribed text
output "transcription" {
  value = openai_transcription.custom.text
}
```

## Argument Reference

The following arguments are supported:

* `audio_content` - (Required, Forces new resource) The base64-encoded audio file content to transcribe. Supported formats: flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.
* `model` - (Optional, Forces new resource) ID of the model to use. Currently only `whisper-1` is supported. Defaults to `whisper-1`.
* `language` - (Optional, Forces new resource) The language of the input audio in ISO-639-1 format (e.g., "en" for English). If not specified, the model will auto-detect the language.
* `prompt` - (Optional, Forces new resource) Optional text to guide the model's style or continue a previous audio segment.
* `response_format` - (Optional, Forces new resource) The format of the transcript output. Must be one of: `json`, `text`, `srt`, `verbose_json`, or `vtt`. Defaults to `text`.
* `temperature` - (Optional, Forces new resource) The sampling temperature, between 0 and 1. Higher values make the output more random, while lower values make it more focused and deterministic. Defaults to 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique identifier for this transcription resource.
* `text` - The transcribed text output.

## Import

This resource cannot be imported as it represents a one-time transcription.