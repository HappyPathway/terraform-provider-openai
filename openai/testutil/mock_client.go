package testutil

import (
	"context"
	"fmt"
	"time"
)

// MockClient is a mock implementation of the ClientInterface
type MockClient struct {
	Models          map[string]*Model
	Files           map[string]*File
	Assistants      map[string]*Assistant
	ChatCompletions map[string]*ChatCompletionResponse
	FineTuning      map[string]*FineTuningJob
	CompletionCount int
}

// NewMockClient creates a new mock client with default test data
func NewMockClient() *MockClient {
	return &MockClient{
		Models: map[string]*Model{
			"gpt-3.5-turbo": {
				ID:      "gpt-3.5-turbo",
				Object:  "model",
				OwnedBy: "openai",
				Permission: []Permission{
					{
						ID:                "perm-123",
						Object:            "permission",
						Created:           1629298613,
						AllowCreateEngine: false,
						AllowSampling:     true,
						AllowFineTuning:   true,
					},
				},
			},
			"gpt-4": {
				ID:      "gpt-4",
				Object:  "model",
				OwnedBy: "openai",
				Permission: []Permission{
					{
						ID:                "perm-456",
						Object:            "permission",
						Created:           1629298613,
						AllowCreateEngine: false,
						AllowSampling:     true,
						AllowFineTuning:   false,
					},
				},
			},
		},
		Files:           make(map[string]*File),
		Assistants:      make(map[string]*Assistant),
		ChatCompletions: make(map[string]*ChatCompletionResponse),
		FineTuning:      make(map[string]*FineTuningJob),
	}
}

// GetModel implements the Model retrieval
func (m *MockClient) GetModel(ctx context.Context, modelID string) (*Model, error) {
	if model, ok := m.Models[modelID]; ok {
		return model, nil
	}
	return nil, fmt.Errorf("model %s not found", modelID)
}

// ListModels implements the Models listing
func (m *MockClient) ListModels(ctx context.Context) ([]Model, error) {
	models := make([]Model, 0, len(m.Models))
	for _, model := range m.Models {
		models = append(models, *model)
	}
	return models, nil
}

// UploadFile implements file upload
func (m *MockClient) UploadFile(ctx context.Context, req *FileUploadRequest) (*File, error) {
	fileID := fmt.Sprintf("file-%d", len(m.Files)+1)
	file := &File{
		ID:        fileID,
		Object:    "file",
		Bytes:     len(req.File),
		CreatedAt: time.Now().Unix(),
		Filename:  "test.jsonl",
		Purpose:   req.Purpose,
	}
	m.Files[fileID] = file
	return file, nil
}

// GetFile implements file retrieval
func (m *MockClient) GetFile(ctx context.Context, fileID string) (*File, error) {
	if file, ok := m.Files[fileID]; ok {
		return file, nil
	}
	return nil, fmt.Errorf("file %s not found", fileID)
}

// DeleteFile implements file deletion
func (m *MockClient) DeleteFile(ctx context.Context, fileID string) error {
	if _, ok := m.Files[fileID]; !ok {
		return fmt.Errorf("file %s not found", fileID)
	}
	delete(m.Files, fileID)
	return nil
}

// CreateAssistant implements assistant creation
func (m *MockClient) CreateAssistant(ctx context.Context, req *CreateAssistantRequest) (*Assistant, error) {
	assistantID := fmt.Sprintf("asst-%d", len(m.Assistants)+1)
	assistant := &Assistant{
		ID:           assistantID,
		Object:       "assistant",
		CreatedAt:    time.Now().Unix(),
		Name:         req.Name,
		Description:  req.Description,
		Model:        req.Model,
		Instructions: req.Instructions,
		Tools:        req.Tools,
		FileIDs:      req.FileIDs,
		Metadata:     req.Metadata,
	}
	m.Assistants[assistantID] = assistant
	return assistant, nil
}

// GetAssistant implements assistant retrieval
func (m *MockClient) GetAssistant(ctx context.Context, assistantID string) (*Assistant, error) {
	if assistant, ok := m.Assistants[assistantID]; ok {
		return assistant, nil
	}
	return nil, fmt.Errorf("assistant %s not found", assistantID)
}

// UpdateAssistant implements assistant update
func (m *MockClient) UpdateAssistant(ctx context.Context, assistantID string, req *CreateAssistantRequest) (*Assistant, error) {
	if assistant, ok := m.Assistants[assistantID]; ok {
		assistant.Name = req.Name
		assistant.Description = req.Description
		assistant.Model = req.Model
		assistant.Instructions = req.Instructions
		assistant.Tools = req.Tools
		assistant.FileIDs = req.FileIDs
		assistant.Metadata = req.Metadata
		return assistant, nil
	}
	return nil, fmt.Errorf("assistant %s not found", assistantID)
}

// DeleteAssistant implements assistant deletion
func (m *MockClient) DeleteAssistant(ctx context.Context, assistantID string) error {
	if _, ok := m.Assistants[assistantID]; !ok {
		return fmt.Errorf("assistant %s not found", assistantID)
	}
	delete(m.Assistants, assistantID)
	return nil
}

// CreateFineTuningJob implements fine-tuning job creation
func (m *MockClient) CreateFineTuningJob(ctx context.Context, req *CreateFineTuningJobRequest) (*FineTuningJob, error) {
	jobID := fmt.Sprintf("ftjob-%d", len(m.FineTuning)+1)
	job := &FineTuningJob{
		ID:             jobID,
		Object:         "fine_tuning.job",
		Model:          req.Model,
		CreatedAt:      time.Now().Unix(),
		FinishedAt:     time.Now().Unix(),
		Status:         "created",
		TrainingFile:   req.TrainingFile,
		ValidationFile: req.ValidationFile,
	}
	m.FineTuning[jobID] = job
	return job, nil
}

// GetFineTuningJob implements fine-tuning job retrieval
func (m *MockClient) GetFineTuningJob(ctx context.Context, jobID string) (*FineTuningJob, error) {
	if job, ok := m.FineTuning[jobID]; ok {
		return job, nil
	}
	return nil, fmt.Errorf("fine-tuning job %s not found", jobID)
}

// CancelFineTuningJob implements fine-tuning job cancellation
func (m *MockClient) CancelFineTuningJob(ctx context.Context, jobID string) (*FineTuningJob, error) {
	if job, ok := m.FineTuning[jobID]; ok {
		job.Status = "cancelled"
		job.FinishedAt = time.Now().Unix()
		return job, nil
	}
	return nil, fmt.Errorf("fine-tuning job %s not found", jobID)
}

// CreateCompletion implements completion creation
func (m *MockClient) CreateCompletion(ctx context.Context, req *CreateCompletionRequest) (*Completion, error) {
	return &Completion{
		ID:      "cmpl-123",
		Object:  "completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []CompletionChoice{
			{
				Text:         "This is a test completion",
				Index:        0,
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

// CreateEmbedding implements embedding creation
func (m *MockClient) CreateEmbedding(ctx context.Context, req *CreateEmbeddingRequest) (*Embedding, error) {
	return &Embedding{
		Object: "embedding",
		Data: []EmbeddingData{
			{
				Object:    "embedding",
				Embedding: []float64{0.1, 0.2, 0.3},
				Index:     0,
			},
		},
		Model: req.Model,
		Usage: EmbeddingUsage{
			PromptTokens: 10,
			TotalTokens:  10,
		},
	}, nil
}

// CreateImage implements image generation
func (m *MockClient) CreateImage(ctx context.Context, req *CreateImageRequest) (*ImageResponse, error) {
	return &ImageResponse{
		Created: time.Now().Unix(),
		Data: []ImageData{
			{
				URL:           "https://example.com/image.png",
				B64JSON:       "",
				RevisedPrompt: req.Prompt,
			},
		},
	}, nil
}

// CreateModeration implements content moderation
func (m *MockClient) CreateModeration(ctx context.Context, req *CreateModerationRequest) (*ModerationResponse, error) {
	return &ModerationResponse{
		ID:    "modr-123",
		Model: req.Model,
		Results: []ModerationResult{
			{
				Flagged:        false,
				Categories:     ModerationCategories{},
				CategoryScores: ModerationCategoryScores{},
			},
		},
	}, nil
}

// CreateSpeech implements text-to-speech
func (m *MockClient) CreateSpeech(ctx context.Context, req *CreateSpeechRequest) (string, error) {
	return "base64-encoded-audio-data", nil
}

// CreateTranscription implements audio transcription
func (m *MockClient) CreateTranscription(ctx context.Context, req *TranscriptionRequest) (*TranscriptionResponse, error) {
	return &TranscriptionResponse{
		Text: "This is a test transcription",
	}, nil
}

// CreateTranslation implements audio translation
func (m *MockClient) CreateTranslation(ctx context.Context, req *TranslationRequest) (*TranslationResponse, error) {
	return &TranslationResponse{
		Text: "This is a test translation",
	}, nil
}

// CreateChatCompletion implements chat completion
func (m *MockClient) CreateChatCompletion(ctx context.Context, req *CreateChatCompletionRequest) (*ChatCompletionResponse, error) {
	if req.Model == "invalid-model" {
		return nil, fmt.Errorf("invalid model specified")
	}

	m.CompletionCount++
	id := fmt.Sprintf("chatcmpl-%d", m.CompletionCount)

	resp := &ChatCompletionResponse{
		ID:      id,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []ChatCompletionChoice{
			{
				Message: ChatCompletionMessage{
					Role:    "assistant",
					Content: "This is a mock response for testing purposes.",
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
		},
	}

	if req.ResponseFormat != nil && req.ResponseFormat["type"] == "json_object" {
		resp.Choices[0].Message.Content = `{"title":"The Hobbit","author":"J.R.R. Tolkien","year_published":1937,"genre":["fantasy","children","adventure"]}`
	}

	m.ChatCompletions[id] = resp
	return resp, nil
}
