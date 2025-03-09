package client

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/time/rate"
)

const (
	// RetryMaxAttempts is the maximum number of retry attempts
	RetryMaxAttempts = 5
	// RetryInitialBackoffMs is the initial backoff in milliseconds
	RetryInitialBackoffMs = 1000
	// RetryMaxBackoffMs is the maximum backoff in milliseconds
	RetryMaxBackoffMs = 30000
)

// Config holds the configuration for the OpenAI client
type Config struct {
	APIKey       string
	BaseURL      string
	Organization string
}

// Client wraps an OpenAI API client for use with the Terraform provider
type Client struct {
	OpenAI      *openai.Client
	debug       bool
	rateLimiter *rate.Limiter
	config      Config
}

// NewClient creates a new OpenAI API client
func NewClient(ctx context.Context, config Config) (*Client, error) {
	client := &Client{
		config: config,
	}

	// Configure OpenAI client
	openaiConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		openaiConfig.BaseURL = config.BaseURL
	}
	if config.Organization != "" {
		openaiConfig.OrgID = config.Organization
	}

	// Set the Assistants API version to v2
	openaiConfig.BaseURL = strings.TrimSuffix(openaiConfig.BaseURL, "/")
	openaiConfig.HTTPClient = &http.Client{
		Transport: &headerTransport{
			base: http.DefaultTransport,
			headers: map[string]string{
				"OpenAI-Beta": "assistants=v2",
			},
		},
	}

	client.OpenAI = openai.NewClientWithConfig(openaiConfig)
	return client, nil
}

// headerTransport adds custom headers to requests
type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}
	return t.base.RoundTrip(req)
}

// LogDebug outputs debug messages if debug mode is enabled
func (c *Client) LogDebug(msg string, additionalFields ...map[string]interface{}) {
	if !c.debug {
		return
	}
	ctx := context.Background()
	fields := make(map[string]interface{})
	// Add additional fields if provided
	if len(additionalFields) > 0 && additionalFields[0] != nil {
		for k, v := range additionalFields[0] {
			fields[k] = v
		}
	}
	tflog.Debug(ctx, msg, fields)
}

// HandleError processes an API error and returns a formatted error message
func (c *Client) HandleError(err error) error {
	if err == nil {
		return nil
	}

	// Check if this is an OpenAI API error
	if apiErr, ok := err.(*openai.APIError); ok {
		// Special handling for rate limit errors
		if apiErr.HTTPStatusCode == 429 {
			return fmt.Errorf("OpenAI API rate limit exceeded: %s. Please retry after a short delay", apiErr.Message)
		}
		// Format other API errors
		return fmt.Errorf("OpenAI API error (Type: %s, Code: %s, Status: %d): %s",
			apiErr.Type, apiErr.Code, apiErr.HTTPStatusCode, apiErr.Message)
	}

	// Generic error
	return fmt.Errorf("error communicating with OpenAI API: %s", err.Error())
}

// ExecuteWithRetry executes a function with retry logic
func (c *Client) ExecuteWithRetry(operation func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error
	var attempt int

	for attempt = 0; attempt < RetryMaxAttempts; attempt++ {
		// Apply rate limiting - wait for our turn
		ctx := context.Background()
		if err := c.rateLimiter.Wait(ctx); err != nil {
			c.LogDebug(fmt.Sprintf("Rate limiter wait failed: %v", err))
			// Continue anyway
		}

		// Execute the operation
		result, err = operation()

		// If successful or non-retryable error, return immediately
		if err == nil || !c.isRetryableError(err) {
			return result, err
		}

		// Calculate backoff with exponential increase and jitter
		backoff := c.calculateBackoff(attempt)
		c.LogDebug(fmt.Sprintf("Retrying after error: %v (attempt %d of %d, waiting %d ms)",
			err, attempt+1, RetryMaxAttempts, backoff/time.Millisecond))
		time.Sleep(backoff)
	}

	return result, fmt.Errorf("operation failed after %d attempts: %v", attempt, err)
}

// isRetryableError determines if an error should trigger a retry
func (c *Client) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check if this is an OpenAI API error
	if apiErr, ok := err.(*openai.APIError); ok {
		// Retry on rate limit errors
		if apiErr.HTTPStatusCode == 429 {
			return true
		}

		// Retry on server errors (5xx)
		if apiErr.HTTPStatusCode >= 500 && apiErr.HTTPStatusCode <= 599 {
			return true
		}

		// Retry on specific error types that might be transient
		if apiErr.Type == "server_error" || apiErr.Type == "timeout" {
			return true
		}
	}

	// Don't retry other errors
	return false
}

// calculateBackoff calculates exponential backoff with jitter
func (c *Client) calculateBackoff(attempt int) time.Duration {
	// Base delay with exponential increase: initialBackoff * 2^attempt
	backoffMs := float64(RetryInitialBackoffMs) * math.Pow(2, float64(attempt))

	// Apply jitter: random value between 0.8 and 1.2 of the base backoff
	jitter := 0.8 + (0.4 * (float64(time.Now().UnixNano()%1000) / 1000.0))
	backoffWithJitterMs := backoffMs * jitter

	// Cap to max backoff
	if backoffWithJitterMs > float64(RetryMaxBackoffMs) {
		backoffWithJitterMs = float64(RetryMaxBackoffMs)
	}

	return time.Duration(backoffWithJitterMs) * time.Millisecond
}

// CreateRunRequest is our internal run creation request type
type CreateRunRequest struct {
	ThreadID            string
	AssistantID         string
	Model               string
	Instructions        string
	Tools               []openai.AssistantTool
	Metadata            map[string]interface{}
	MaxPromptTokens     int
	MaxCompletionTokens int
}

// CreateRun creates a new run for a thread
func (c *Client) CreateRun(ctx context.Context, req *CreateRunRequest) (*openai.Run, error) {
	runRequest := openai.RunRequest{
		AssistantID:  req.AssistantID,
		Model:        req.Model,
		Instructions: req.Instructions,
	}

	if req.Metadata != nil {
		runRequest.Metadata = req.Metadata
	}

	// Set token control parameters if provided
	if req.MaxPromptTokens > 0 {
		runRequest.MaxPromptTokens = req.MaxPromptTokens
	}
	if req.MaxCompletionTokens > 0 {
		runRequest.MaxCompletionTokens = req.MaxCompletionTokens
	}

	// Convert AssistantTool to Tool
	var tools []openai.Tool
	for _, tool := range req.Tools {
		var toolType openai.ToolType
		switch tool.Type {
		case openai.AssistantToolTypeCodeInterpreter, openai.AssistantToolTypeRetrieval, openai.AssistantToolTypeFunction:
			toolType = openai.ToolTypeFunction
		default:
			return nil, fmt.Errorf("unsupported tool type: %s", tool.Type)
		}
		tools = append(tools, openai.Tool{
			Type: toolType,
		})
	}
	runRequest.Tools = tools

	run, err := c.OpenAI.CreateRun(ctx, req.ThreadID, runRequest)
	if err != nil {
		return nil, fmt.Errorf("error creating run: %v", err)
	}
	return &run, nil
}

// GetRun retrieves a run by ID and thread ID
func (c *Client) GetRun(ctx context.Context, id string, threadID string) (*openai.Run, error) {
	run, err := c.OpenAI.RetrieveRun(ctx, threadID, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving run: %v", err)
	}
	return &run, nil
}

// CancelRun cancels a run
func (c *Client) CancelRun(ctx context.Context, id string, threadID string) error {
	// First check the run's status
	run, err := c.GetRun(ctx, id, threadID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil // Run doesn't exist, nothing to cancel
		}
		return fmt.Errorf("error getting run status: %v", err)
	}

	// If the run is already in a terminal state, just return
	switch run.Status {
	case openai.RunStatusCompleted, openai.RunStatusFailed, openai.RunStatusCancelled, openai.RunStatusExpired:
		return nil
	}

	// Only try to cancel if the run is in a cancellable state
	_, err = c.OpenAI.CancelRun(ctx, threadID, id)
	if err != nil {
		return fmt.Errorf("error cancelling run: %v", err)
	}
	return nil
}
