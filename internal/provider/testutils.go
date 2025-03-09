package provider

import (
	"os"
	"testing"
)

func TestAccPreCheck(t *testing.T) {
	if os.Getenv("OPENAI_API_TOKEN") == "" {
		t.Skip("OPENAI_API_TOKEN must be set for acceptance tests")
	}
}
