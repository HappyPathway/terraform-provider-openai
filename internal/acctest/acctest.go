package acctest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/darnold/terraform-provider-openai/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ProviderConfig returns a configuration string with the test credentials
func ProviderConfig() string {
	return fmt.Sprintf(`
provider "openai" {
  api_key = %q
}
`, os.Getenv("OPENAI_API_KEY"))
}

// PreCheck verifies that the required environment variables are set for acceptance tests
func PreCheck(t *testing.T) {
	if v := os.Getenv("OPENAI_API_KEY"); v == "" {
		t.Fatal("OPENAI_API_KEY environment variable must be set for acceptance tests")
	}
}

// ProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for each Terraform CLI command executed to create
// a provider server to which the CLI can reattach.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"openai": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// SkipIfEmptyEnv skips a test if the specified environment variable is not set
func SkipIfEmptyEnv(t *testing.T, envVar string) {
	if os.Getenv(envVar) == "" {
		t.Skipf("Environment variable %s is not set", envVar)
	}
}

// CheckDestroyOpenAIResource is a helper for checking the resource is destroyed
func CheckDestroyOpenAIResource(resourceType string, checkFunc func(s *terraform.State) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			// Check the resource is truly gone
			if err := checkFunc(s); err != nil {
				return err
			}
		}

		return nil
	}
}

// SkipIfNotAcceptanceTest skips a test unless it's an acceptance test
func SkipIfNotAcceptanceTest(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}
}

// TestContext returns a context for testing
func TestContext() context.Context {
	return context.Background()
}

// RandomName returns a random name with a prefix for test resources
func RandomName(prefix string) string {
	// Use the TestName function to get the test name
	return fmt.Sprintf("%s-%s", prefix, resource.UniqueId())
}
