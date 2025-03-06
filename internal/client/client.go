package client

import (
	"context"
	"fmt"
	"math"
	"net/http"
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

// Client wraps an OpenAI API client for use with the Terraform provider
type Client struct {
	OpenAI      *openai.Client
	httpClient  *http.Client
	apiKey      string
	baseURL     string
	debug       bool
	rateLimiter *rate.Limiter
}

// NewClient creates a new OpenAI API client
func NewClient(apiKey, organization, baseURL string, debug bool) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key must not be empty")
	}

	// Configure HTTP client with reasonable defaults
	httpClient := &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// Create the base client config
	config := openai.DefaultConfig(apiKey)
	config.HTTPClient = httpClient
	config.AssistantVersion = "v2" // Set the AssistantVersion to v2

	// Add organization option if set
	if organization != "" {
		config.OrgID = organization
	}

	// Add custom base URL if set
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	// Create the OpenAI client
	client := openai.NewClientWithConfig(config)

	// Create a rate limiter that allows for 60 requests per minute (default OpenAI rate limit)
	// with a burst of 5 requests
	rateLimiter := rate.NewLimiter(rate.Limit(1), 5)

	return &Client{
		OpenAI:      client,
		httpClient:  httpClient,
		apiKey:      apiKey,
		baseURL:     baseURL,
		debug:       debug,
		rateLimiter: rateLimiter,
	}, nil
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
