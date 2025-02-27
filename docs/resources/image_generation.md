---
page_title: "openai_image_generation Resource - terraform-provider-openai"
description: |-
  Generates images using OpenAI's DALL-E models.
---

# openai_image_generation (Resource)

Generates images using OpenAI's DALL-E models. This resource allows you to create images based on text descriptions using either DALL-E 2 or DALL-E 3 models.

~> **Note:** This resource creates a new image each time it is created and cannot be updated. Changes to any of the arguments will result in a new image being generated.

## Example Usage

```terraform
# Generate a single image with DALL-E 3
resource "openai_image_generation" "example" {
  prompt = "A serene mountain landscape with a misty lake at sunrise"
  model  = "dall-e-3"
  size   = "1024x1024"
  style  = "natural"
}

# Generate multiple images with DALL-E 2
resource "openai_image_generation" "multiple" {
  prompt          = "A futuristic cityscape at night"
  model           = "dall-e-2"
  n               = 3
  size            = "1024x1024"
  response_format = "url"
}

# Output the generated image URLs
output "image_urls" {
  value = openai_image_generation.multiple.images[*].url
}
```

## Argument Reference

The following arguments are supported:

* `prompt` - (Required, Forces new resource) A text description of the desired image(s). Maximum length is 1000 characters for dall-e-2 and 4000 characters for dall-e-3.
* `model` - (Optional, Forces new resource) The model to use for image generation. Must be one of "dall-e-2" or "dall-e-3". Defaults to "dall-e-2".
* `n` - (Optional, Forces new resource) The number of images to generate. Must be between 1 and 10. For dall-e-3, only n=1 is supported. Defaults to 1.
* `quality` - (Optional, Forces new resource) The quality of the image that will be generated. Must be one of "standard" or "hd". HD creates images with finer details. Only supported for dall-e-3. Defaults to "standard".
* `response_format` - (Optional, Forces new resource) The format in which the generated images are returned. Must be one of "url" or "b64_json". Defaults to "url".
* `size` - (Optional, Forces new resource) The size of the generated images. For dall-e-2, must be one of "256x256", "512x512", or "1024x1024". For dall-e-3, must be one of "1024x1024", "1792x1024", or "1024x1792". Defaults to "1024x1024".
* `style` - (Optional, Forces new resource) The style of the generated images. Must be one of "vivid" or "natural". Only supported for dall-e-3. Defaults to "vivid".
* `user` - (Optional, Forces new resource) A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - A list of generated images. Each image contains:
  * `b64_json` - The base64-encoded JSON of the generated image (only when response_format is "b64_json").
  * `url` - The URL of the generated image (only when response_format is "url").
  * `revised_prompt` - The prompt that was used to generate the image, if there was any revision.