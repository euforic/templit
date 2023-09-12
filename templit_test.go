package templit

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"text/template"
)

// TestRenderTemplate tests the RenderTemplate function.
func TestRenderTemplate(t *testing.T) {
	var tests = []struct {
		name     string
		tmpl     string
		data     interface{}
		funcMap  template.FuncMap
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple template without data",
			tmpl:     "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "Template with data",
			tmpl:     "Hello, {{.Name}}!",
			data:     map[string]string{"Name": "John"},
			expected: "Hello, John!",
		},
		{
			name:     "Template with function",
			tmpl:     "Hello, {{lower .Name}}!",
			data:     map[string]string{"Name": "JOHN"},
			funcMap:  template.FuncMap{"lower": func(s string) string { return strings.ToLower(s) }},
			expected: "Hello, john!",
		},
		{
			name:    "Malformed template",
			tmpl:    "Hello, {{.Name!",
			data:    map[string]string{"Name": "John"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderTemplate(tt.tmpl, tt.data, tt.funcMap)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("RenderTemplate() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// TestWalkAndProcessDir tests the WalkAndProcessDir function.
func TestWalkAndProcessDir(t *testing.T) {
	tests := []struct {
		name           string
		inputDir       string
		data           interface{}
		funcMap        template.FuncMap
		expectedOutput string
		expectedError  error
	}{
		{
			name:     "Simple directory processing",
			inputDir: "test_data/templates/basic_test",
			data: map[string]interface{}{
				"Name":        "John",
				"Title":       "Project",
				"Description": "This is a test project.",
				"Detail":      "more info here.",
			},
			funcMap:        DefaultFuncMap,
			expectedOutput: "test_data/outputs/basic_test/",
		},
		// ... (other test cases)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for output
			tempOutputDir, err := os.MkdirTemp("", "test_output_"+strings.ReplaceAll(strings.ToLower(tt.name), " ", "_"))
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			//defer os.RemoveAll(tempOutputDir) // Cleanup

			err = WalkAndProcessDir(tt.inputDir, tempOutputDir, tt.funcMap, tt.data)
			if err != nil {
				if tt.expectedError == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			// Compare generated output directory with expected directory
			fmt.Println(tt.expectedOutput, tempOutputDir)
			err = compareDirs(tt.expectedOutput, tempOutputDir)
			if err != nil {
				t.Fatalf("output directory mismatch: %v", err)
			}
		})
	}
}
