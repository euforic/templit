package templit

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestEmbedFunc tests the EmbedFunc function.
func TestEmbedFunc(t *testing.T) {
	tests := []struct {
		name          string
		repoAndPath   string
		ctx           interface{}
		expectedFile  string
		expectedText  string
		expectedError error
	}{
		{
			name:         "Valid template fetch and execute",
			repoAndPath:  "https://test_data/templates/basic_test/greeting.txt@main",
			ctx:          map[string]string{"Name": "John"},
			expectedFile: "test_data/outputs/basic_test/greeting.txt",
		},
		{
			name:         "Valid template fetch and execute block",
			repoAndPath:  "https://test_data/templates/basic_test/-block.txt#example_block@main",
			ctx:          map[string]string{"Greeting": "Hey"},
			expectedText: "Hey, this is an example block.",
		},
		{
			name:          "Invalid repo path format",
			repoAndPath:   "invalidpath",
			ctx:           nil,
			expectedText:  "",
			expectedError: fmt.Errorf("invalid path format in embed URL"),
		},
		{
			name:          "Invalid repo",
			repoAndPath:   "https://localhost/owner/invalidrepo/greeting.txt@main",
			ctx:           nil,
			expectedFile:  "",
			expectedError: fmt.Errorf("failed to clone repo: open localhost/owner/invalidrepo: no such file or directory"),
		},
	}

	client := &MockGitClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := EmbedFunc(client)
			result, err := fn(tt.repoAndPath, tt.ctx, nil)
			if err != nil {
				if tt.expectedError == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			var expected string
			if tt.expectedFile != "" {
				expectedBytes, _ := os.ReadFile(tt.expectedFile)
				expected = string(expectedBytes)
			} else {
				expected = tt.expectedText
			}

			if diff := cmp.Diff(string(expected), result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
