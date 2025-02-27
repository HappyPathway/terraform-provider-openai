package testutil

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
	CreatedAt int    `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// Assistant represents an OpenAI assistant
type Assistant struct {
	ID           string            `json:"id"`
	Object       string            `json:"object"`
	CreatedAt    int              `json:"created_at"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Model        string           `json:"model"`
	Instructions string           `json:"instructions"`
	Tools        []AssistantTool  `json:"tools"`
	FileIDs      []string         `json:"file_ids"`
	Metadata     map[string]string `json:"metadata"`
}

// AssistantTool represents a tool that can be used by an assistant
type AssistantTool struct {
	Type string `json:"type"`
}

// FineTuningJob represents a fine-tuning job
type FineTuningJob struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	Model          string `json:"model"`
	CreatedAt      int    `json:"created_at"`
	FinishedAt     int    `json:"finished_at"`
	Status         string `json:"status"`
	TrainingFile   string `json:"training_file"`
	ValidationFile string `json:"validation_file"`
}

// CreateAssistantRequest represents the request to create an assistant
type CreateAssistantRequest struct {
	Name         string            `json:"name"`
	Description  string           `json:"description"`
	Model        string           `json:"model"`
	Instructions string           `json:"instructions"`
	Tools        []AssistantTool  `json:"tools"`
	FileIDs      []string         `json:"file_ids"`
	Metadata     map[string]string `json:"metadata"`
}

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	File    []byte `json:"file"`
	Purpose string `json:"purpose"`
}

// CreateFineTuningJobRequest represents a request to create a fine-tuning job
type CreateFineTuningJobRequest struct {
	Model          string `json:"model"`
	TrainingFile   string `json:"training_file"`
	ValidationFile string `json:"validation_file"`
	Hyperparameters struct {
		NEpochs int `json:"n_epochs"`
	} `json:"hyperparameters"`
	Suffix string `json:"suffix"`
}