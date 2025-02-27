package openai

import "testing"

// skipUnimplementedTest skips tests for features that are not yet implemented
func skipUnimplementedTest(t *testing.T, feature string) {
	t.Skipf("Skipping test for unimplemented feature: %s", feature)
}

// List of currently implemented features
var implementedFeatures = map[string]bool{
	"model":             true, // data source
	"models":            true, // data source
	"file":              true, // resource
	"assistant":         true, // resource
	"content_generator": true, // resource
}

// isFeatureImplemented checks if a feature is implemented
func isFeatureImplemented(feature string) bool {
	return implementedFeatures[feature]
}
