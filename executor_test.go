package templit

import (
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
)

// TestNewExecutor tests the NewExecutor function.
func TestNewExecutor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		funcMap  template.FuncMap
		expected string
		err      bool
	}{
		{
			name:    "valid template file",
			input:   "test_data/templates/basic_test/greeting.txt",
			funcMap: nil,
			err:     false,
		},
		// ... add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor, err := NewExecutor(tt.input, tt.funcMap)
			if (err != nil) != tt.err {
				t.Fatalf("expected error: %v, got: %v", tt.err, err)
			}
			if executor == nil {
				t.Fatalf("expected executor to not be nil")
			}
		})
	}
}

// TestRender tests the Render function.
func TestRender(t *testing.T) {
	tests := []struct {
		name         string
		inputPath    string
		templateName string
		data         interface{}
		expected     string
		err          bool
	}{
		{
			name:         "valid template",
			inputPath:    "test_data/templates/basic_test",
			templateName: "test_data/templates/basic_test/greeting.txt",
			data:         map[string]string{"Name": "John"},
			expected:     "Hello, John!\n",
			err:          false,
		},
		{
			name:         "valid template block",
			inputPath:    "test_data/templates/basic_test",
			templateName: "example_block",
			data:         map[string]string{"greeting": "Hey"},
			expected:     "Hey, this is an example block.",
			err:          false,
		},

		// ... add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor, err := NewExecutor(tt.inputPath, nil)
			if err != nil {
				t.Fatalf("failed to create executor: %v", err)
			}
			result, err := executor.Render(tt.templateName, tt.data)
			if (err != nil) != tt.err {
				t.Fatalf("expected error: %v, got: %v", tt.err, err)
			}
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Fatalf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}
