package openai

import (
	"bytes"
	"context"
	"encoding/base64"
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
	Type string `json:"type"`
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

// Completion represents an OpenAI completion response
type Completion struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// CompletionChoice represents a choice in a completion response
type CompletionChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CreateCompletionRequest represents parameters for creating a completion
type CreateCompletionRequest struct {
	Model            string             `json:"model"`
	Prompt           string             `json:"prompt"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	N                int                `json:"n,omitempty"`
	BestOf           int                `json:"best_of,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	Echo             bool               `json:"echo,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
	Seed             int                `json:"seed,omitempty"`
	Suffix           string             `json:"suffix,omitempty"`
}

// Embedding represents an OpenAI embedding response
type Embedding struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  EmbeddingUsage  `json:"usage"`
}

// EmbeddingData represents the embedding vector data
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingUsage represents token usage for embeddings
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// CreateEmbeddingRequest represents parameters for creating embeddings
type CreateEmbeddingRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	User           string `json:"user,omitempty"`
	Dimensions     int    `json:"dimensions,omitempty"`
	EncodingFormat string `json:"encoding_format,omitempty"`
}

// CreateImageRequest represents parameters for generating images
type CreateImageRequest struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	N              int    `json:"n,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Size           string `json:"size,omitempty"`
	Style          string `json:"style,omitempty"`
	User           string `json:"user,omitempty"`
}

// ImageResponse represents a response from the image generation API
type ImageResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

// ImageData represents generated image data
type ImageData struct {
	B64JSON       string `json:"b64_json,omitempty"`
	URL           string `json:"url,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// ModerationCategories represents the categories of harmful content
type ModerationCategories struct {
	Harassment            bool `json:"harassment"`
	HarassmentThreatening bool `json:"harassment/threatening"`
	Hate                  bool `json:"hate"`
	HateThreatening       bool `json:"hate/threatening"`
	Illicit               bool `json:"illicit"`
	IllicitViolent        bool `json:"illicit/violent"`
	SelfHarm              bool `json:"self-harm"`
	SelfHarmInstructions  bool `json:"self-harm/instructions"`
	SelfHarmIntent        bool `json:"self-harm/intent"`
	Sexual                bool `json:"sexual"`
	SexualMinors          bool `json:"sexual/minors"`
	Violence              bool `json:"violence"`
	ViolenceGraphic       bool `json:"violence/graphic"`
}

// ModerationCategoryScores represents the scores for each category
type ModerationCategoryScores struct {
	Harassment            float64 `json:"harassment"`
	HarassmentThreatening float64 `json:"harassment/threatening"`
	Hate                  float64 `json:"hate"`
	HateThreatening       float64 `json:"hate/threatening"`
	Illicit               float64 `json:"illicit"`
	IllicitViolent        float64 `json:"illicit/violent"`
	SelfHarm              float64 `json:"self-harm"`
	SelfHarmInstructions  float64 `json:"self-harm/instructions"`
	SelfHarmIntent        float64 `json:"self-harm/intent"`
	Sexual                float64 `json:"sexual"`
	SexualMinors          float64 `json:"sexual/minors"`
	Violence              float64 `json:"violence"`
	ViolenceGraphic       float64 `json:"violence/graphic"`
}

// ModerationResult represents a single moderation result
type ModerationResult struct {
	Flagged        bool                     `json:"flagged"`
	Categories     ModerationCategories     `json:"categories"`
	CategoryScores ModerationCategoryScores `json:"category_scores"`
}

// ModerationResponse represents the response from the moderation API
type ModerationResponse struct {
	ID      string             `json:"id"`
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
}

// CreateModerationRequest represents the parameters for creating a moderation
type CreateModerationRequest struct {
	Input string `json:"input"`
	Model string `json:"model,omitempty"`
}

// CreateSpeechRequest represents parameters for generating speech
type CreateSpeechRequest struct {
	Input          string  `json:"input"`
	Model          string  `json:"model"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
}

// TranscriptionRequest represents parameters for transcribing audio
type TranscriptionRequest struct {
	File           []byte  `json:"file"`
	Model          string  `json:"model"`
	Language       string  `json:"language,omitempty"`
	Prompt         string  `json:"prompt,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
}

// TranscriptionResponse represents the response from a transcription request
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// TranslationRequest represents parameters for translating audio
type TranslationRequest struct {
	File           []byte  `json:"file"`
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
}

// TranslationResponse represents the response from a translation request
type TranslationResponse struct {
	Text string `json:"text"`
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

// CreateCompletion creates a completion
func (c *Client) CreateCompletion(ctx context.Context, req *CreateCompletionRequest) (*Completion, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/completions", c.baseURL), bytes.NewReader(body))
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

	var completion Completion
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &completion, nil
}

// CreateEmbedding creates embeddings for the given input
func (c *Client) CreateEmbedding(ctx context.Context, req *CreateEmbeddingRequest) (*Embedding, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/embeddings", c.baseURL), bytes.NewReader(body))
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

	var embedding Embedding
	if err := json.Unmarshal(respBody, &embedding); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &embedding, nil
}

// CreateImage generates images using DALL·E models
func (c *Client) CreateImage(ctx context.Context, req *CreateImageRequest) (*ImageResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/images/generations", c.baseURL), bytes.NewReader(body))
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

	var imageResponse ImageResponse
	if err := json.Unmarshal(respBody, &imageResponse); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &imageResponse, nil
}

// CreateModeration creates a new moderation
func (c *Client) CreateModeration(ctx context.Context, req *CreateModerationRequest) (*ModerationResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/moderations", c.baseURL), bytes.NewReader(body))
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

	var moderation ModerationResponse
	if err := json.Unmarshal(respBody, &moderation); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &moderation, nil
}

// CreateSpeech generates audio from the input text
func (c *Client) CreateSpeech(ctx context.Context, req *CreateSpeechRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/audio/speech", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(httpReq)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, respBody)
	}

	// Convert audio bytes to base64
	return base64.StdEncoding.EncodeToString(respBody), nil
}

// CreateTranscription transcribes audio into text
func (c *Client) CreateTranscription(ctx context.Context, req *TranscriptionRequest) (*TranscriptionResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file
	part, err := writer.CreateFormFile("file", "audio.mp3")
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}
	if _, err := part.Write(req.File); err != nil {
		return nil, fmt.Errorf("error writing file data: %v", err)
	}

	// Add other fields
	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, fmt.Errorf("error writing model field: %v", err)
	}
	if req.Language != "" {
		if err := writer.WriteField("language", req.Language); err != nil {
			return nil, fmt.Errorf("error writing language field: %v", err)
		}
	}
	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, fmt.Errorf("error writing prompt field: %v", err)
		}
	}
	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, fmt.Errorf("error writing response_format field: %v", err)
		}
	}
	if req.Temperature != 0 {
		if err := writer.WriteField("temperature", fmt.Sprintf("%f", req.Temperature)); err != nil {
			return nil, fmt.Errorf("error writing temperature field: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/audio/transcriptions", c.baseURL), body)
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

	var transcription TranscriptionResponse
	if err := json.Unmarshal(respBody, &transcription); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &transcription, nil
}

// CreateTranslation translates audio into English text
func (c *Client) CreateTranslation(ctx context.Context, req *TranslationRequest) (*TranslationResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file
	part, err := writer.CreateFormFile("file", "audio.mp3")
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}
	if _, err := part.Write(req.File); err != nil {
		return nil, fmt.Errorf("error writing file data: %v", err)
	}

	// Add other fields
	if err := writer.WriteField("model", req.Model); err != nil {
		return nil, fmt.Errorf("error writing model field: %v", err)
	}
	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, fmt.Errorf("error writing prompt field: %v", err)
		}
	}
	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", req.ResponseFormat); err != nil {
			return nil, fmt.Errorf("error writing response_format field: %v", err)
		}
	}
	if req.Temperature != 0 {
		if err := writer.WriteField("temperature", fmt.Sprintf("%f", req.Temperature)); err != nil {
			return nil, fmt.Errorf("error writing temperature field: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/audio/translations", c.baseURL), body)
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

	var translation TranslationResponse
	if err := json.Unmarshal(respBody, &translation); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &translation, nil
}

// ClientInterface defines the interface that both our real and mock clients must implement
type ClientInterface interface {
	GetModel(ctx context.Context, modelID string) (*Model, error)
	ListModels(ctx context.Context) ([]Model, error)
	UploadFile(ctx context.Context, req *FileUploadRequest) (*File, error)
	GetFile(ctx context.Context, fileID string) (*File, error)
	DeleteFile(ctx context.Context, fileID string) error
	CreateAssistant(ctx context.Context, req *CreateAssistantRequest) (*Assistant, error)
	GetAssistant(ctx context.Context, assistantID string) (*Assistant, error)
	UpdateAssistant(ctx context.Context, assistantID string, req *CreateAssistantRequest) (*Assistant, error)
	DeleteAssistant(ctx context.Context, assistantID string) error
	CreateFineTuningJob(ctx context.Context, req *CreateFineTuningJobRequest) (*FineTuningJob, error)
	GetFineTuningJob(ctx context.Context, jobID string) (*FineTuningJob, error)
	CancelFineTuningJob(ctx context.Context, jobID string) (*FineTuningJob, error)
	CreateCompletion(ctx context.Context, req *CreateCompletionRequest) (*Completion, error)
	CreateEmbedding(ctx context.Context, req *CreateEmbeddingRequest) (*Embedding, error)
	CreateImage(ctx context.Context, req *CreateImageRequest) (*ImageResponse, error)
	CreateModeration(ctx context.Context, req *CreateModerationRequest) (*ModerationResponse, error)
	CreateSpeech(ctx context.Context, req *CreateSpeechRequest) (string, error)
	CreateTranscription(ctx context.Context, req *TranscriptionRequest) (*TranscriptionResponse, error)
	CreateTranslation(ctx context.Context, req *TranslationRequest) (*TranslationResponse, error)
}
