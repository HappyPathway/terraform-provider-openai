package testutil

import (
	"os"
	"testing"
)

// PreCheck verifies required test env vars are set
func PreCheck(t *testing.T) {
	if os.Getenv("OPENAI_MOCK") == "" {
		// Only check API key when not using mocks
		if v := os.Getenv("OPENAI_API_KEY"); v == "" {
			t.Fatal("OPENAI_API_KEY must be set for acceptance tests")
		}
	}
}

// ProviderTestPreCheck verifies testing variables are set
func ProviderTestPreCheck(t *testing.T) {
	if os.Getenv("OPENAI_MOCK") == "" {
		// Only check project variables when not using mocks
		if v := os.Getenv("TF_VAR_project_prompt"); v == "" {
			t.Skip("TF_VAR_project_prompt not set")
		}
		if v := os.Getenv("TF_VAR_repo_org"); v == "" {
			t.Skip("TF_VAR_repo_org not set")
		}
		if v := os.Getenv("TF_VAR_project_name"); v == "" {
			t.Skip("TF_VAR_project_name not set")
		}
	}
}
