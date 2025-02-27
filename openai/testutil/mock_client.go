package testutil

import (
	"context"
	"fmt"
)

// MockClient is a mock implementation of the OpenAI client
type MockClient struct {
	Models      map[string]*Model
	Files       map[string]*File
	Assistants  map[string]*Assistant
	FineTuning  map[string]*FineTuningJob
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
		},
		Files:      make(map[string]*File),
		Assistants: make(map[string]*Assistant),
		FineTuning: make(map[string]*FineTuningJob),
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
		CreatedAt: 1629298613,
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
		job.FinishedAt = 1629298614
		return job, nil
	}
	return nil, fmt.Errorf("fine-tuning job %s not found", jobID)
}