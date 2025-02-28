package openai

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenAIMessage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpenAIMessageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenAIMessageBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenAIMessageExists("openai_message.test"),
					resource.TestCheckResourceAttr("openai_message.test", "role", "user"),
					resource.TestCheckResourceAttrSet("openai_message.test", "created_at"),
					resource.TestCheckResourceAttrSet("openai_message.test", "thread_id"),
				),
			},
			{
				ResourceName:      "openai_message.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOpenAIMessage_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpenAIMessageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenAIMessageWithMetadata(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenAIMessageExists("openai_message.test"),
					resource.TestCheckResourceAttr("openai_message.test", "metadata.test", "value"),
				),
			},
		},
	})
}

func testAccCheckOpenAIMessageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Message ID is set")
		}

		client := testAccProvider.Meta().(*Config).Client
		threadID := rs.Primary.Attributes["thread_id"]

		_, err := client.Beta.Threads.Messages.Get(context.Background(), threadID, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckOpenAIMessageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openai_message" {
			continue
		}

		threadID := rs.Primary.Attributes["thread_id"]
		_, err := client.Beta.Threads.Messages.Get(context.Background(), threadID, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Message still exists")
		}
	}

	return nil
}

func testAccOpenAIMessageBasic() string {
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
}`
}

func testAccOpenAIMessageWithMetadata() string {
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

  metadata = {
    test = "value"
  }
}`
}
