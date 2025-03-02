package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/acctest"
	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock OpenAI client
type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(openai.ChatCompletionResponse), args.Error(1)
}

func TestChatCompletionResource_Metadata(t *testing.T) {
	r := &ChatCompletionResource{}
	req := resource.MetadataRequest{
		ProviderTypeName: "openai",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)
	assert.Equal(t, "openai_chat_completion", resp.TypeName)
}

func TestChatCompletionResource_Schema(t *testing.T) {
	r := &ChatCompletionResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	// Verify schema structure
	assert.NotNil(t, resp.Schema)
	assert.NotNil(t, resp.Schema.Attributes["id"])
	assert.NotNil(t, resp.Schema.Attributes["model"])
	assert.NotNil(t, resp.Schema.Attributes["messages"])
	assert.NotNil(t, resp.Schema.Attributes["response_content"])
}

func TestChatCompletionResource_Create(t *testing.T) {
	// Set up mock client
	mockClient := &MockOpenAIClient{}
	client := &client.Client{
		OpenAI: mockClient,
	}

	// Set up resource
	r := &ChatCompletionResource{
		client: client,
	}

	// Mock API response
	created := int64(1699573200)
	mockResponse := openai.ChatCompletionResponse{
		ID:      "chatcmpl-123456789",
		Created: created,
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Role:    "assistant",
					Content: "This is a test response",
				},
			},
		},
	}

	// Set up expected API call
	mockClient.On("CreateChatCompletion", mock.Anything, mock.MatchedBy(func(req openai.ChatCompletionRequest) bool {
		return req.Model == "gpt-4" && len(req.Messages) == 1 && req.Messages[0].Role == "user"
	})).Return(mockResponse, nil)

	// Create plan
	messageElements := []attr.Value{
		types.ObjectValue(
			map[string]attr.Type{
				"role":    types.StringType,
				"content": types.StringType,
			},
			map[string]attr.Value{
				"role":    types.StringValue("user"),
				"content": types.StringValue("Hello, how are you?"),
			},
		),
	}

	messagesList, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"role":    types.StringType,
				"content": types.StringType,
			},
		},
		messageElements,
	)

	plan := ChatCompletionResourceModel{
		Model:    types.StringValue("gpt-4"),
		Messages: messagesList,
	}

	// Create request and response
	schema := schema.Schema{}
	req := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Schema: schema,
		},
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: schema,
		},
	}

	// Set the plan in the request
	req.Plan.Set(context.Background(), plan)

	// Call the method under test
	r.Create(context.Background(), req, resp)

	// Assertions
	assert.False(t, resp.Diagnostics.HasError())

	var resultState ChatCompletionResourceModel
	resp.State.Get(context.Background(), &resultState)

	assert.Equal(t, fmt.Sprintf("chat-gpt-4-%d", created), resultState.ID.ValueString())
	assert.Equal(t, "assistant", resultState.ResponseRole.ValueString())

	// Verify the response content
	var contents []string
	resultState.ResponseContent.ElementsAs(context.Background(), &contents, false)
	assert.Equal(t, 1, len(contents))
	assert.Equal(t, "This is a test response", contents[0])

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

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
