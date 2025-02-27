package testutil

import (
	"context"
	"fmt"

	"github.com/HappyPathway/terraform-provider-openai/openai"
)

// MockClient is a mock implementation of the OpenAI client
type MockClient struct {
	Models     map[string]*openai.Model
	Files      map[string]*openai.File
	Assistants map[string]*openai.Assistant
	FineTuning map[string]*openai.FineTuningJob
}

// NewMockClient creates a new mock client with default test data
func NewMockClient() *MockClient {
	return &MockClient{
		Models: map[string]*openai.Model{
			"gpt-3.5-turbo": {
				ID:      "gpt-3.5-turbo",
				Object:  "model",
				OwnedBy: "openai",
				Permission: []openai.Permission{
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
		},
		Files:      make(map[string]*openai.File),
		Assistants: make(map[string]*openai.Assistant),
		FineTuning: make(map[string]*openai.FineTuningJob),
	}
}

// GetModel implements the Model retrieval
func (m *MockClient) GetModel(ctx context.Context, modelID string) (*openai.Model, error) {
	if model, ok := m.Models[modelID]; ok {
		return model, nil
	}
	return nil, fmt.Errorf("model %s not found", modelID)
}

// ListModels implements the Models listing
func (m *MockClient) ListModels(ctx context.Context) ([]openai.Model, error) {
	models := make([]openai.Model, 0, len(m.Models))
	for _, model := range m.Models {
		models = append(models, *model)
	}
	return models, nil
}

// UploadFile implements file upload
func (m *MockClient) UploadFile(ctx context.Context, req *openai.FileUploadRequest) (*openai.File, error) {
	fileID := fmt.Sprintf("file-%d", len(m.Files)+1)
	file := &openai.File{
		ID:        fileID,
		Object:    "file",
		Bytes:     len(req.File),
		CreatedAt: 1629298613,
		Filename:  "test.jsonl",
		Purpose:   req.Purpose,
	}
	m.Files[fileID] = file
	return file, nil
}

// GetFile implements file retrieval
func (m *MockClient) GetFile(ctx context.Context, fileID string) (*openai.File, error) {
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
func (m *MockClient) CreateAssistant(ctx context.Context, req *openai.CreateAssistantRequest) (*openai.Assistant, error) {
	assistantID := fmt.Sprintf("asst-%d", len(m.Assistants)+1)
	assistant := &openai.Assistant{
		ID:           assistantID,
		Object:       "assistant",
		CreatedAt:    1629298613,
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
func (m *MockClient) GetAssistant(ctx context.Context, assistantID string) (*openai.Assistant, error) {
	if assistant, ok := m.Assistants[assistantID]; ok {
		return assistant, nil
	}
	return nil, fmt.Errorf("assistant %s not found", assistantID)
}

// UpdateAssistant implements assistant update
func (m *MockClient) UpdateAssistant(ctx context.Context, assistantID string, req *openai.CreateAssistantRequest) (*openai.Assistant, error) {
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
func (m *MockClient) CreateFineTuningJob(ctx context.Context, req *openai.CreateFineTuningJobRequest) (*openai.FineTuningJob, error) {
	jobID := fmt.Sprintf("ftjob-%d", len(m.FineTuning)+1)
	job := &openai.FineTuningJob{
		ID:             jobID,
		Object:         "fine_tuning.job",
		Model:          req.Model,
		CreatedAt:      1629298613,
		FinishedAt:     0,
		Status:         "created",
		TrainingFile:   req.TrainingFile,
		ValidationFile: req.ValidationFile,
	}
	m.FineTuning[jobID] = job
	return job, nil
}

// GetFineTuningJob implements fine-tuning job retrieval
func (m *MockClient) GetFineTuningJob(ctx context.Context, jobID string) (*openai.FineTuningJob, error) {
	if job, ok := m.FineTuning[jobID]; ok {
		return job, nil
	}
	return nil, fmt.Errorf("fine-tuning job %s not found", jobID)
}

// CancelFineTuningJob implements fine-tuning job cancellation
func (m *MockClient) CancelFineTuningJob(ctx context.Context, jobID string) (*openai.FineTuningJob, error) {
	if job, ok := m.FineTuning[jobID]; ok {
		job.Status = "cancelled"
		job.FinishedAt = 1629298614
		return job, nil
	}
	return nil, fmt.Errorf("fine-tuning job %s not found", jobID)
}

// CreateCompletion implements completion creation
func (m *MockClient) CreateCompletion(ctx context.Context, req *openai.CreateCompletionRequest) (*openai.Completion, error) {
	return &openai.Completion{
		ID:      "cmpl-123",
		Object:  "completion",
		Created: 1629298613,
		Model:   req.Model,
		Choices: []openai.CompletionChoice{
			{
				Text:         "This is a test completion",
				Index:        0,
				FinishReason: "stop",
			},
		},
		Usage: openai.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

// CreateEmbedding implements embedding creation
func (m *MockClient) CreateEmbedding(ctx context.Context, req *openai.CreateEmbeddingRequest) (*openai.Embedding, error) {
	return &openai.Embedding{
		Object: "embedding",
		Data: []openai.EmbeddingData{
			{
				Object:    "embedding",
				Embedding: []float64{0.1, 0.2, 0.3},
				Index:     0,
			},
		},
		Model: req.Model,
		Usage: openai.EmbeddingUsage{
			PromptTokens: 10,
			TotalTokens:  10,
		},
	}, nil
}

// CreateImage implements image generation
func (m *MockClient) CreateImage(ctx context.Context, req *openai.CreateImageRequest) (*openai.ImageResponse, error) {
	return &openai.ImageResponse{
		Created: 1629298613,
		Data: []openai.ImageData{
			{
				URL:           "https://example.com/image.png",
				B64JSON:       "",
				RevisedPrompt: req.Prompt,
			},
		},
	}, nil
}

// CreateModeration implements content moderation
func (m *MockClient) CreateModeration(ctx context.Context, req *openai.CreateModerationRequest) (*openai.ModerationResponse, error) {
	return &openai.ModerationResponse{
		ID:    "modr-123",
		Model: req.Model,
		Results: []openai.ModerationResult{
			{
				Flagged:        false,
				Categories:     openai.ModerationCategories{},
				CategoryScores: openai.ModerationCategoryScores{},
			},
		},
	}, nil
}

// CreateSpeech implements text-to-speech
func (m *MockClient) CreateSpeech(ctx context.Context, req *openai.CreateSpeechRequest) (string, error) {
	return "base64-encoded-audio-data", nil
}

// CreateTranscription implements audio transcription
func (m *MockClient) CreateTranscription(ctx context.Context, req *openai.TranscriptionRequest) (*openai.TranscriptionResponse, error) {
	return &openai.TranscriptionResponse{
		Text: "This is a test transcription",
	}, nil
}

// CreateTranslation implements audio translation
func (m *MockClient) CreateTranslation(ctx context.Context, req *openai.TranslationRequest) (*openai.TranslationResponse, error) {
	return &openai.TranslationResponse{
		Text: "This is a test translation",
	}, nil
}
