package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOpenAIMessage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenAIMessageConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.openai_message.test", "created_at"),
					resource.TestCheckResourceAttr("data.openai_message.test", "role", "user"),
					resource.TestCheckResourceAttr("data.openai_message.test", "content.0.type", "text"),
					resource.TestCheckResourceAttr("data.openai_message.test", "content.0.text", "This is a test message"),
				),
			},
		},
	})
}

func testAccDataSourceOpenAIMessageConfig() string {
	return `
resource "openai_thread" "test" {
  metadata = {
    test = "value"
  }
}

resource "openai_message" "test" {
  thread_id = openai_thread.test.id
  role      = "user"
  
  content {
    type = "text"
    text = "This is a test message"
  }
}

data "openai_message" "test" {
  thread_id  = openai_thread.test.id
  message_id = openai_message.test.id
}`
}
