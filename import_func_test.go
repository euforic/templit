package templit_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/euforic/templit"
)

// TestImportFunc tests the ImportFunc function.
func TestImportFunc(t *testing.T) {
	tests := []struct {
		name          string
		repoAndTag    string
		data          interface{}
		expectedError error
	}{
		{
			// Add your test cases here
			name:       "Valid template processing",
			repoAndTag: "https://test_data/templates/basic_test@main",
			data: map[string]interface{}{
				"Name":        "John",
				"Title":       "Project",
				"Description": "This is a test project.",
				"Detail":      "more info here.",
			},
		},
	}

	// Mocked GitClient
	client := &MockGitClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			destPath, err := os.MkdirTemp("", "test_output_"+strings.ReplaceAll(strings.ToLower(tt.name), " ", "_"))
			if err != nil {
				t.Fatal(fmt.Errorf("failed to create temp dir: %w", err))
			}

			defer os.RemoveAll(destPath) // Cleanup

			executor := templit.NewExecutor(client)
			fn := executor.ImportFunc(destPath)
			if _, err := fn(tt.repoAndTag, "./", tt.data); err != nil {
				if tt.expectedError == nil || errors.Is(err, tt.expectedError) {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			// Compare generated files with expected outputs
			repo := strings.Split(tt.repoAndTag, "@")[0]
			expectedOutputPath := strings.Replace(repo, "templates", "outputs", 1)

			err = compareDirs(destPath, strings.TrimPrefix(expectedOutputPath, "https://"))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
