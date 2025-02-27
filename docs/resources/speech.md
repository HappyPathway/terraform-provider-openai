---
page_title: "openai_speech Resource - terraform-provider-openai"
description: |-
  Generates audio from text using OpenAI's Text-to-Speech (TTS) models.
---

# openai_speech (Resource)

This resource allows you to generate audio from text using OpenAI's Text-to-Speech (TTS) models. The generated audio is returned in base64 encoded format for easy storage and manipulation.

~> **Note:** This resource creates new audio each time it is created and cannot be updated. Changes to any of the arguments will result in new audio being generated.

## Example Usage

```terraform
# Generate speech with default settings
resource "openai_speech" "basic" {
  input = "Hello, this is a test of the OpenAI text to speech system."
  voice = "alloy"
}

# Generate high-quality speech with custom settings
resource "openai_speech" "custom" {
  input           = "Welcome to our application! We're excited to have you here."
  model           = "tts-1-hd"
  voice           = "nova"
  response_format = "mp3"
  speed           = 1.2
}

# Output the base64-encoded audio content
output "audio_content" {
  value = openai_speech.custom.audio_content
}
```

## Argument Reference

The following arguments are supported:

* `input` - (Required, Forces new resource) The text to generate audio for. Maximum length is 4096 characters.
* `model` - (Optional, Forces new resource) The TTS model to use. One of `tts-1` or `tts-1-hd`. Defaults to `tts-1`.
* `voice` - (Required, Forces new resource) The voice to use for the audio. Must be one of: `alloy`, `ash`, `coral`, `echo`, `fable`, `onyx`, `nova`, `sage`, or `shimmer`.
* `response_format` - (Optional, Forces new resource) The format to generate audio in. Must be one of: `mp3`, `opus`, `aac`, `flac`, `wav`, or `pcm`. Defaults to `mp3`.
* `speed` - (Optional, Forces new resource) The speed of the generated audio. Must be between `0.25` and `4.0`. Defaults to `1.0`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique identifier for this speech resource.
* `audio_content` - The generated audio content in base64 encoded format.

## Import

This resource cannot be imported as it represents a one-time audio generation.