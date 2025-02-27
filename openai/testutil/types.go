package testutil

import (
	"context"
)

// ClientInterface defines the interface that both real and mock clients must implement
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
	CreateChatCompletion(ctx context.Context, req *CreateChatCompletionRequest) (*ChatCompletionResponse, error)
}

// Model represents an OpenAI model
type Model struct {
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
}

// Permission represents model permissions
type Permission struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	AllowCreateEngine bool   `json:"allow_create_engine"`
	AllowSampling     bool   `json:"allow_sampling"`
	AllowFineTuning   bool   `json:"allow_fine_tuning"`
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
	Type     string             `json:"type"`
	Function *AssistantFunction `json:"function,omitempty"`
}

// AssistantFunction represents a function that can be called by an assistant
type AssistantFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// FineTuningJobError represents an error in a fine-tuning job
type FineTuningJobError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
}

// FineTuningJob represents a fine-tuning job
type FineTuningJob struct {
	ID             string              `json:"id"`
	Object         string              `json:"object"`
	Model          string              `json:"model"`
	CreatedAt      int64               `json:"created_at"`
	FinishedAt     int64               `json:"finished_at"`
	Status         string              `json:"status"`
	TrainingFile   string              `json:"training_file"`
	ValidationFile string              `json:"validation_file"`
	Error          *FineTuningJobError `json:"error,omitempty"`
}

// CreateAssistantRequest represents the request to create an assistant
type CreateAssistantRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Model        string            `json:"model"`
	Instructions string            `json:"instructions"`
	Tools        []AssistantTool   `json:"tools"`
	FileIDs      []string          `json:"file_ids"`
	Metadata     map[string]string `json:"metadata"`
}

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	File     []byte `json:"file"`
	Purpose  string `json:"purpose"`
	Filename string `json:"filename"`
}

// CreateFineTuningJobRequest represents a request to create a fine-tuning job
type CreateFineTuningJobRequest struct {
	Model           string `json:"model"`
	TrainingFile    string `json:"training_file"`
	ValidationFile  string `json:"validation_file"`
	Hyperparameters struct {
		NEpochs int `json:"n_epochs"`
	} `json:"hyperparameters"`
	Suffix string `json:"suffix"`
}

// All the other OpenAI API types

type CompletionChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Completion struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

type CreateCompletionRequest struct {
	Model            string             `json:"model"`
	Prompt           string             `json:"prompt"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	Temperature      float32            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	N                int                `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	Echo             bool               `json:"echo,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	BestOf           int                `json:"best_of,omitempty"`
	User             string             `json:"user,omitempty"`
	Seed             int                `json:"seed,omitempty"`
	Suffix           string             `json:"suffix,omitempty"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type Embedding struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  EmbeddingUsage  `json:"usage"`
}

type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type CreateEmbeddingRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	Dimensions     int    `json:"dimensions,omitempty"`
	EncodingFormat string `json:"encoding_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageData struct {
	URL           string `json:"url"`
	B64JSON       string `json:"b64_json"`
	RevisedPrompt string `json:"revised_prompt"`
}

type ImageResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

type CreateImageRequest struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	N              int    `json:"n"`
	Quality        string `json:"quality"`
	ResponseFormat string `json:"response_format"`
	Size           string `json:"size"`
	Style          string `json:"style"`
	User           string `json:"user"`
}

type ModerationCategories struct {
	Hate                  bool `json:"hate"`
	HateThreatening       bool `json:"hate/threatening"`
	Harassment            bool `json:"harassment"`
	HarassmentThreatening bool `json:"harassment/threatening"`
	SelfHarm              bool `json:"self-harm"`
	SelfHarmInstructions  bool `json:"self-harm/instructions"`
	Sexual                bool `json:"sexual"`
	SexualMinors          bool `json:"sexual/minors"`
	Violence              bool `json:"violence"`
	ViolenceGraphic       bool `json:"violence/graphic"`
	Illicit               bool `json:"illicit"`
	IllicitViolent        bool `json:"illicit/violent"`
}

type ModerationCategoryScores struct {
	Hate                  float32 `json:"hate"`
	HateThreatening       float32 `json:"hate/threatening"`
	Harassment            float32 `json:"harassment"`
	HarassmentThreatening float32 `json:"harassment/threatening"`
	SelfHarm              float32 `json:"self-harm"`
	SelfHarmInstructions  float32 `json:"self-harm/instructions"`
	Sexual                float32 `json:"sexual"`
	SexualMinors          float32 `json:"sexual/minors"`
	Violence              float32 `json:"violence"`
	ViolenceGraphic       float32 `json:"violence/graphic"`
	Illicit               float32 `json:"illicit"`
	IllicitViolent        float32 `json:"illicit/violent"`
}

type ModerationResult struct {
	Flagged        bool                     `json:"flagged"`
	Categories     ModerationCategories     `json:"categories"`
	CategoryScores ModerationCategoryScores `json:"category_scores"`
}

type ModerationResponse struct {
	ID      string             `json:"id"`
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
}

type CreateModerationRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type CreateSpeechRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
}

type TranscriptionRequest struct {
	File           []byte  `json:"file"`
	Model          string  `json:"model"`
	Language       string  `json:"language"`
	Temperature    float32 `json:"temperature"`
	Prompt         string  `json:"prompt"`
	ResponseFormat string  `json:"response_format"`
}

type TranscriptionResponse struct {
	Text string `json:"text"`
}

type TranslationRequest struct {
	File           []byte  `json:"file"`
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt"`
	ResponseFormat string  `json:"response_format"`
	Temperature    float32 `json:"temperature"`
}

type TranslationResponse struct {
	Text string `json:"text"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionChoice struct {
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

type CreateChatCompletionRequest struct {
	Model          string                  `json:"model"`
	Messages       []ChatCompletionMessage `json:"messages"`
	Temperature    float32                 `json:"temperature"`
	MaxTokens      int                     `json:"max_tokens"`
	ResponseFormat map[string]string       `json:"response_format,omitempty"`
}
