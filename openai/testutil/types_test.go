package testutil

import (
	"testing"

	"github.com/HappyPathway/terraform-provider-openai/openai"
)

func TestModelTypes(t *testing.T) {
	model := openai.Model{
		ID:      "test-model",
		Object:  "model",
		OwnedBy: "openai",
		Permission: []openai.Permission{
			{
				ID:                "perm-test",
				Object:            "permission",
				Created:           1629298613,
				AllowCreateEngine: false,
				AllowSampling:     true,
				AllowFineTuning:   true,
			},
		},
	}

	if model.ID != "test-model" {
		t.Errorf("model ID = %v, want %v", model.ID, "test-model")
	}
	if model.OwnedBy != "openai" {
		t.Errorf("model OwnedBy = %v, want %v", model.OwnedBy, "openai")
	}
	if len(model.Permission) != 1 {
		t.Errorf("model Permission count = %v, want %v", len(model.Permission), 1)
	}

	perm := model.Permission[0]
	if !perm.AllowFineTuning {
		t.Errorf("permission AllowFineTuning = %v, want %v", perm.AllowFineTuning, true)
	}
}

func TestMockClient(t *testing.T) {
	client := NewMockClient()

	// Test model operations
	model, err := client.GetModel(nil, "gpt-3.5-turbo")
	if err != nil {
		t.Errorf("GetModel error = %v", err)
	}
	if model.ID != "gpt-3.5-turbo" {
		t.Errorf("model ID = %v, want %v", model.ID, "gpt-3.5-turbo")
	}

	models, err := client.ListModels(nil)
	if err != nil {
		t.Errorf("ListModels error = %v", err)
	}
	if len(models) != 1 {
		t.Errorf("models count = %v, want %v", len(models), 1)
	}

	// Test file operations
	fileReq := &openai.FileUploadRequest{
		File:    []byte("test data"),
		Purpose: "fine-tune",
	}
	file, err := client.UploadFile(nil, fileReq)
	if err != nil {
		t.Errorf("UploadFile error = %v", err)
	}
	if file.Purpose != "fine-tune" {
		t.Errorf("file Purpose = %v, want %v", file.Purpose, "fine-tune")
	}

	// Test assistant operations
	assistantReq := &openai.CreateAssistantRequest{
		Name:        "test assistant",
		Description: "test description",
		Model:       "gpt-3.5-turbo",
		Tools: []openai.AssistantTool{
			{Type: "code_interpreter"},
		},
	}
	assistant, err := client.CreateAssistant(nil, assistantReq)
	if err != nil {
		t.Errorf("CreateAssistant error = %v", err)
	}
	if assistant.Name != "test assistant" {
		t.Errorf("assistant Name = %v, want %v", assistant.Name, "test assistant")
	}
}
