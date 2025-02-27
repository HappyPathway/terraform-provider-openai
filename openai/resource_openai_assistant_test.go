package openai

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOpenAIAssistant_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + testAccResourceOpenAIAssistantConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "name", "Test Assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "description", "Test assistant for acceptance tests"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "instructions", "You are a test assistant."),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.#", "1"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.0.type", "code_interpreter"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "metadata.test", "true"),
				),
			},
			{
				Config: testAccProviderConfig + testAccResourceOpenAIAssistantConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "name", "Updated Test Assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "description", "Updated test assistant"),
					resource.TestCheckResourceAttr(
						"openai_assistant.test", "tools.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceAssistant_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAssistantConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test", "name", "Test Assistant"),
					resource.TestCheckResourceAttr("openai_assistant.test", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_assistant.test", "instructions", "You are a helpful assistant."),
					resource.TestCheckResourceAttrSet("openai_assistant.test", "id"),
				),
			},
			{
				ResourceName:      "openai_assistant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAssistant_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAssistantConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test_full", "name", "Full Test Assistant"),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "model", "gpt-4"),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "description", "A test assistant with full configuration"),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "instructions", "You are a specialized test assistant."),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "tools.#", "2"),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "tools.0.type", "code_interpreter"),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "tools.1.type", "retrieval"),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "metadata.environment", "test"),
					resource.TestCheckResourceAttr("openai_assistant.test_full", "metadata.version", "1.0"),
					resource.TestCheckResourceAttrSet("openai_assistant.test_full", "id"),
				),
			},
		},
	})
}

func TestAccResourceAssistant_withFiles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAssistantConfig_withFiles(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test_files", "name", "Assistant With Files"),
					resource.TestCheckResourceAttr("openai_assistant.test_files", "tools.#", "2"),
					resource.TestCheckResourceAttr("openai_assistant.test_files", "tools.0.type", "code_interpreter"),
					resource.TestCheckResourceAttr("openai_assistant.test_files", "tools.1.type", "retrieval"),
					resource.TestCheckResourceAttr("openai_assistant.test_files", "file_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceAssistant_withFunction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAssistantConfig_withFunction(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test_function", "name", "Assistant With Function"),
					resource.TestCheckResourceAttr("openai_assistant.test_function", "tools.#", "1"),
					resource.TestCheckResourceAttr("openai_assistant.test_function", "tools.0.type", "function"),
					resource.TestCheckResourceAttr("openai_assistant.test_function", "tools.0.function.name", "get_weather"),
					resource.TestCheckResourceAttr("openai_assistant.test_function", "tools.0.function.description", "Get the current weather for a location"),
				),
			},
		},
	})
}

func TestAccResourceAssistant_withAllTools(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAssistantConfig_withAllTools(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openai_assistant.test_all_tools", "name", "Assistant With All Tools"),
					resource.TestCheckResourceAttr("openai_assistant.test_all_tools", "tools.#", "3"),
					resource.TestCheckResourceAttr("openai_assistant.test_all_tools", "tools.0.type", "code_interpreter"),
					resource.TestCheckResourceAttr("openai_assistant.test_all_tools", "tools.1.type", "retrieval"),
					resource.TestCheckResourceAttr("openai_assistant.test_all_tools", "tools.2.type", "function"),
				),
			},
		},
	})
}

func TestAccResourceAssistant_invalidModel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAssistantConfig_invalidModel(),
				ExpectError: regexp.MustCompile(`Invalid model ID`),
			},
		},
	})
}

func TestAccResourceAssistant_invalidToolType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAssistantConfig_invalidToolType(),
				ExpectError: regexp.MustCompile(`Invalid tool type`),
			},
		},
	})
}

func TestAccResourceAssistant_invalidFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAssistantConfig_invalidFile(),
				ExpectError: regexp.MustCompile(`Invalid file purpose`),
			},
		},
	})
}

func TestAccResourceAssistant_invalidFunction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAssistantConfig_invalidFunction(),
				ExpectError: regexp.MustCompile(`Invalid function parameters`),
			},
		},
	})
}

func TestAccResourceAssistant_tooManyTools(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAssistantConfig_tooManyTools(),
				ExpectError: regexp.MustCompile(`exceeds maximum allowed tools \(128\)`),
			},
		},
	})
}

func testAccResourceOpenAIAssistantConfig() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test" {
  name         = "Test Assistant"
  description  = "Test assistant for acceptance tests"
  model        = "gpt-3.5-turbo"
  instructions = "You are a test assistant."

  tools {
    type = "code_interpreter"
  }

  metadata = {
    test = "true"
  }
}
`)
}

func testAccResourceOpenAIAssistantConfigUpdated() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test" {
  name         = "Updated Test Assistant"
  description  = "Updated test assistant"
  model        = "gpt-3.5-turbo"
  instructions = "You are a test assistant."

  metadata = {
    test = "true"
  }
}
`)
}

func testAccResourceAssistantConfig_basic() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test" {
  name         = "Test Assistant"
  model        = "gpt-4"
  instructions = "You are a helpful assistant."
}
`)
}

func testAccResourceAssistantConfig_full() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test_full" {
  name        = "Full Test Assistant"
  model       = "gpt-4"
  description = "A test assistant with full configuration"
  instructions = "You are a specialized test assistant."
  tools = [
    {
      type = "code_interpreter"
    },
    {
      type = "retrieval"
    }
  ]
  metadata = {
    environment = "test"
    version     = "1.0"
  }
}
`)
}

func testAccResourceAssistantConfig_withFiles() string {
	return fmt.Sprintf(`
resource "openai_file" "knowledge" {
  content  = "This is example knowledge base content."
  filename = "knowledge.txt"
  purpose  = "assistants"
}

resource "openai_assistant" "test_files" {
  name         = "Assistant With Files"
  model        = "gpt-4"
  description  = "An assistant that uses files and code interpreter"
  instructions = "You are a helpful assistant with access to files and code capabilities."
  
  tools = [
    {
      type = "code_interpreter"
    },
    {
      type = "retrieval"
    }
  ]
  
  file_ids = [openai_file.knowledge.id]
}
`)
}

func testAccResourceAssistantConfig_withFunction() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test_function" {
  name         = "Assistant With Function"
  model        = "gpt-4"
  description  = "An assistant that uses function calling"
  instructions = "You are a helpful assistant that can check the weather."
  
  tools = [
    {
      type = "function"
      function = {
        name = "get_weather"
        description = "Get the current weather for a location"
        parameters = jsonencode({
          type = "object"
          properties = {
            location = {
              type = "string"
              description = "The location to get weather for, e.g. 'San Francisco, CA'"
            }
            unit = {
              type = "string"
              enum = ["celsius", "fahrenheit"]
              description = "The unit of temperature to return"
            }
          }
          required = ["location"]
        })
      }
    }
  ]
}
`)
}

func testAccResourceAssistantConfig_withAllTools() string {
	return fmt.Sprintf(`
resource "openai_file" "knowledge_base" {
  content  = "This is example knowledge base content for all tools test."
  filename = "knowledge_base.txt"
  purpose  = "assistants"
}

resource "openai_assistant" "test_all_tools" {
  name         = "Assistant With All Tools"
  model        = "gpt-4"
  description  = "An assistant that uses all available tools"
  instructions = "You are a helpful assistant with access to all tools."
  
  tools = [
    {
      type = "code_interpreter"
    },
    {
      type = "retrieval"
    },
    {
      type = "function"
      function = {
        name = "get_current_time"
        description = "Get the current time in a given timezone"
        parameters = jsonencode({
          type = "object"
          properties = {
            timezone = {
              type = "string"
              description = "The timezone to get current time for, e.g. 'America/New_York'"
            }
          }
          required = ["timezone"]
        })
      }
    }
  ]
  
  file_ids = [openai_file.knowledge_base.id]
}
`)
}

func testAccResourceAssistantConfig_invalidModel() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test_invalid_model" {
    name  = "Invalid Model Assistant"
    model = "invalid-model-id"
}
`)
}

func testAccResourceAssistantConfig_invalidToolType() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test_invalid_tool" {
    name  = "Invalid Tool Assistant"
    model = "gpt-4"
    tools = [
        {
            type = "invalid_tool_type"
        }
    ]
}
`)
}

func testAccResourceAssistantConfig_invalidFile() string {
	return fmt.Sprintf(`
resource "openai_file" "invalid_purpose" {
    content  = "Test content"
    filename = "test.txt"
    purpose  = "fine-tune"  # Should be "assistants"
}

resource "openai_assistant" "test_invalid_file" {
    name     = "Invalid File Assistant"
    model    = "gpt-4"
    tools    = [
        {
            type = "retrieval"
        }
    ]
    file_ids = [openai_file.invalid_purpose.id]
}
`)
}

func testAccResourceAssistantConfig_invalidFunction() string {
	return fmt.Sprintf(`
resource "openai_assistant" "test_invalid_function" {
    name  = "Invalid Function Assistant"
    model = "gpt-4"
    tools = [
        {
            type = "function"
            function = {
                name = "invalid_function"
                description = "Invalid function definition"
                parameters = jsonencode({
                    type = "invalid_type",  # Should be "object"
                    properties = {
                        param = {
                            type = "invalid_type"  # Invalid type
                        }
                    }
                })
            }
        }
    ]
}
`)
}

func testAccResourceAssistantConfig_tooManyTools() string {
	// Generate 129 tools (exceeding the 128 limit)
	var tools strings.Builder
	for i := 0; i < 129; i++ {
		tools.WriteString(`
        {
            type = "function"
            function = {
                name = "function_${i}"
                description = "Test function ${i}"
                parameters = jsonencode({
                    type = "object"
                    properties = {
                        param = {
                            type = "string"
                            description = "Test parameter"
                        }
                    }
                })
            }
        },`)
	}

	return fmt.Sprintf(`
resource "openai_assistant" "test_too_many_tools" {
    name  = "Too Many Tools Assistant"
    model = "gpt-4"
    tools = [
        %s
    ]
}
`, tools.String())
}
