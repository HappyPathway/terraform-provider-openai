package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	defaultBaseURL     = "https://api.openai.com/v1"
	defaultRetryMax    = 3
	defaultRetryDelay  = 1 * time.Second
	defaultHTTPTimeout = 30 * time.Second
)

type ClientConfig struct {
	BaseURL    string
	APIKey     string
	RetryMax   int
	RetryDelay time.Duration
	Timeout    time.Duration
}

// Client handles communication with the OpenAI API
type Client struct {
	baseURL    string
	apiKey     string
	retryMax   int
	retryDelay time.Duration
	httpClient *http.Client
}

// Model represents an OpenAI model
type Model struct {
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
}

// Permission represents the permissions for a model
type Permission struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	AllowCreateEngine bool   `json:"allow_create_engine"`
	AllowSampling     bool   `json:"allow_sampling"`
	AllowFineTuning   bool   `json:"allow_fine_tuning"`
}

type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// File represents an OpenAI file
type File struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// FileUploadRequest represents the parameters for uploading a file
type FileUploadRequest struct {
	File    []byte
	Purpose string
}

// Assistant represents an OpenAI assistant
type Assistant struct {
	ID           string            `json:"id"`
	Object       string            `json:"object"`
	CreatedAt    int64             `json:"created_at"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Model        string            `json:"model"`
	Instructions string            `json:"instructions"`
	Tools        []AssistantTool   `json:"tools"`
	FileIDs      []string          `json:"file_ids"`
	Metadata     map[string]string `json:"metadata"`
}

// AssistantTool represents a tool that can be used by an assistant
type AssistantTool struct {
	Type     string              `json:"type"`
	Function *FunctionDefinition `json:"function,omitempty"`
}

// FunctionDefinition represents the definition of a function tool
type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  string `json:"parameters"`
}

// CreateAssistantRequest represents the parameters for creating an assistant
type CreateAssistantRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Model        string            `json:"model"`
	Instructions string            `json:"instructions,omitempty"`
	Tools        []AssistantTool   `json:"tools,omitempty"`
	FileIDs      []string          `json:"file_ids,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// FineTuningJob represents an OpenAI fine-tuning job
type FineTuningJob struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	Model          string `json:"model"`
	CreatedAt      int64  `json:"created_at"`
	FinishedAt     int64  `json:"finished_at,omitempty"`
	Status         string `json:"status"`
	TrainingFile   string `json:"training_file"`
	ValidationFile string `json:"validation_file,omitempty"`
	Error          *Error `json:"error,omitempty"`
}

// Error represents an error from the OpenAI API
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateFineTuningJobRequest represents the parameters for creating a fine-tuning job
type CreateFineTuningJobRequest struct {
	Model           string `json:"model"`
	TrainingFile    string `json:"training_file"`
	ValidationFile  string `json:"validation_file,omitempty"`
	Hyperparameters struct {
		NEpochs int `json:"n_epochs,omitempty"`
	} `json:"hyperparameters,omitempty"`
	Suffix string `json:"suffix,omitempty"`
}

// assistantRequestBody represents the request body for creating/updating an assistant
type assistantRequestBody struct {
	Model        string            `json:"model"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Instructions string            `json:"instructions,omitempty"`
	Tools        []AssistantTool   `json:"tools,omitempty"`
	FileIDs      []string          `json:"file_ids,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Chat completion types
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionChoice struct {
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletion struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Usage   CompletionUsage        `json:"usage"`
	Choices []ChatCompletionChoice `json:"choices"`
}

type CreateChatCompletionRequest struct {
	Model          string                  `json:"model"`
	Messages       []ChatCompletionMessage `json:"messages"`
	Temperature    float32                 `json:"temperature,omitempty"`
	ResponseFormat *ResponseFormat         `json:"response_format,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

// NewClient creates a new OpenAI API client
func NewClient(apiKey string) *Client {
	return NewClientWithConfig(ClientConfig{
		APIKey: apiKey,
	})
}

// NewClientWithConfig creates a new OpenAI API client with custom configuration
func NewClientWithConfig(config ClientConfig) *Client {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	if config.RetryMax == 0 {
		config.RetryMax = defaultRetryMax
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = defaultRetryDelay
	}
	if config.Timeout == 0 {
		config.Timeout = defaultHTTPTimeout
	}

	return &Client{
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
		retryMax:   config.RetryMax,
		retryDelay: config.RetryDelay,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// doRequest performs an HTTP request with retries
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.retryMax; i++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if i == c.retryMax {
				return nil, fmt.Errorf("max retries reached: %v", err)
			}
			time.Sleep(c.retryDelay)
			continue
		}

		// Retry on rate limit errors
		if resp.StatusCode == http.StatusTooManyRequests {
			if i == c.retryMax {
				return resp, nil
			}
			time.Sleep(c.retryDelay)
			continue
		}

		return resp, nil
	}

	return resp, err
}

// GetModel retrieves information about a specific model
func (c *Client) GetModel(ctx context.Context, modelID string) (*Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/models/%s", c.baseURL, modelID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var model Model
	if err := json.Unmarshal(body, &model); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &model, nil
}

// ListModels retrieves the list of available models
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/models", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var modelsResp ModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return modelsResp.Data, nil
}

// UploadFile uploads a file to OpenAI
func (c *Client) UploadFile(ctx context.Context, req *FileUploadRequest) (*File, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}
	_, err = part.Write(req.File)
	if err != nil {
		return nil, fmt.Errorf("error writing file data: %v", err)
	}

	err = writer.WriteField("purpose", req.Purpose)
	if err != nil {
		return nil, fmt.Errorf("error writing purpose field: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/files", c.baseURL), body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, respBody)
	}

	var file File
	if err := json.Unmarshal(respBody, &file); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &file, nil
}

// GetFile retrieves a file by ID
func (c *Client) GetFile(ctx context.Context, fileID string) (*File, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/files/%s", c.baseURL, fileID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var file File
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &file, nil
}

// DeleteFile deletes a file by ID
func (c *Client) DeleteFile(ctx context.Context, fileID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/files/%s", c.baseURL, fileID), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	return nil
}

// CreateAssistant creates a new assistant
func (c *Client) CreateAssistant(ctx context.Context, request *CreateAssistantRequest) (*Assistant, error) {
	reqBody := assistantRequestBody{
		Model:        request.Model,
		Name:         request.Name,
		Description:  request.Description,
		Instructions: request.Instructions,
		Tools:        request.Tools,
		FileIDs:      request.FileIDs,
		Metadata:     request.Metadata,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/assistants", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, respBody)
	}

	var assistant Assistant
	if err := json.Unmarshal(respBody, &assistant); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &assistant, nil
}

// GetAssistant retrieves an assistant by ID
func (c *Client) GetAssistant(ctx context.Context, assistantID string) (*Assistant, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/assistants/%s", c.baseURL, assistantID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var assistant Assistant
	if err := json.Unmarshal(body, &assistant); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &assistant, nil
}

// UpdateAssistant updates an existing assistant
func (c *Client) UpdateAssistant(ctx context.Context, assistantID string, req *CreateAssistantRequest) (*Assistant, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/assistants/%s", c.baseURL, assistantID), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, respBody)
	}

	var assistant Assistant
	if err := json.Unmarshal(respBody, &assistant); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &assistant, nil
}

// DeleteAssistant deletes an assistant by ID
func (c *Client) DeleteAssistant(ctx context.Context, assistantID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/assistants/%s", c.baseURL, assistantID), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	return nil
}

// CreateFineTuningJob creates a new fine-tuning job
func (c *Client) CreateFineTuningJob(ctx context.Context, req *CreateFineTuningJobRequest) (*FineTuningJob, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/fine_tuning/jobs", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, respBody)
	}

	var job FineTuningJob
	if err := json.Unmarshal(respBody, &job); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &job, nil
}

// GetFineTuningJob retrieves a fine-tuning job by ID
func (c *Client) GetFineTuningJob(ctx context.Context, jobID string) (*FineTuningJob, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/fine_tuning/jobs/%s", c.baseURL, jobID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var job FineTuningJob
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &job, nil
}

// CancelFineTuningJob cancels a fine-tuning job
func (c *Client) CancelFineTuningJob(ctx context.Context, jobID string) (*FineTuningJob, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/fine_tuning/jobs/%s/cancel", c.baseURL, jobID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var job FineTuningJob
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &job, nil
}

// CreateChatCompletion creates a new chat completion
func (c *Client) CreateChatCompletion(ctx context.Context, req *CreateChatCompletionRequest) (*ChatCompletion, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/chat/completions", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, respBody)
	}

	var completion ChatCompletion
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &completion, nil
}
