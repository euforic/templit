package templit_test

import (
	"testing"

	"github.com/euforic/templit"
	"github.com/google/go-cmp/cmp"
)

// TestToCamelCase tests the ToCamelCase function.
func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"hello world", "helloWorld"},
		{"Hello World", "helloWorld"},
		{"HELLO_WORLD", "helloWorld"},
		{"XML HTTP_request2_a-b", "xmlHttpRequest2AB"},
	}

	for _, test := range tests {
		result := templit.ToCamelCase(test.input)
		if diff := cmp.Diff(test.output, result); diff != "" {
			t.Errorf("toCamelCase(%s) mismatch (-want +got):\n%s", test.input, diff)
		}
	}
}

// TestToSnakeCase tests the ToSnakeCase function.
func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"hello world", "hello_world"},
		{"Hello World", "hello_world"},
		{"HELLO_WORLD", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"HelloWorld Today", "hello_world_today"},
		{"HelloWorld-today", "hello_world_today"},
		{"XML HTTP_request2_a-b", "xml_http_request2_a_b"},
	}

	for _, test := range tests {
		result := templit.ToSnakeCase(test.input)
		if diff := cmp.Diff(test.output, result); diff != "" {
			t.Errorf("ToSnakeCase(%s) mismatch (-want +got):\n%s", test.input, diff)
		}
	}
}

// TestToKebabCase tests the ToKebabCase function.
func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"hello world", "hello-world"},
		{"Hello World", "hello-world"},
		{"HELLO_WORLD", "hello-world"},
		{"HelloWorld", "hello-world"},
		{"HelloWorld Today", "hello-world-today"},
		{"HelloWorld-today", "hello-world-today"},
		{"XML HTTP_request2_a-b", "xml-http-request2-a-b"},
	}

	for _, test := range tests {
		result := templit.ToKebabCase(test.input)
		if diff := cmp.Diff(test.output, result); diff != "" {
			t.Errorf("ToKebabCase(%s) mismatch (-want +got):\n%s", test.input, diff)
		}
	}
}

// TestToPascalCase tests the ToPascalCase function.
func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"hello world", "HelloWorld"},
		{"Hello World", "HelloWorld"},
		{"HELLO_WORLD", "HelloWorld"},
		{"XML HTTP_request2_a-b", "XmlHttpRequest2AB"},
	}

	for _, test := range tests {
		result := templit.ToPascalCase(test.input)
		if diff := cmp.Diff(test.output, result); diff != "" {
			t.Errorf("ToPascalCase(%s) mismatch (-want +got):\n%s", test.input, diff)
		}
	}
}
