package chat_completion_test

import (
	"fmt"
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccChatCompletionResource_Basic tests the chat completion resource with live API calls
func TestAccChatCompletionResource_Basic(t *testing.T) {
	acctest.SkipIfNotAcceptanceTest(t)
	resourceName := "openai_chat_completion.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChatCompletionResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "model", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(resourceName, "response_role", "assistant"),
					resource.TestCheckResourceAttrSet(resourceName, "response_content.0"),
				),
			},
		},
	})
}

func testAccChatCompletionResourceConfig() string {
	return fmt.Sprintf(`
%s

resource "openai_chat_completion" "test" {
  model = "gpt-3.5-turbo"
  
  messages = [
    {
      role    = "system"
      content = "You are a helpful assistant."
    },
    {
      role    = "user"
      content = "What is the capital of France?"
    }
  ]
  
  temperature = 0.7
  max_tokens = 100
}
`, acctest.ProviderConfig())
}
