package templit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Executor is a wrapper around the template.Template type
type Executor struct {
	*template.Template
}

// New returns a new Executor
func NewExecutor() *Executor {
	return &Executor{
		Template: template.New("main").Funcs(DefaultFuncMap),
	}
}

// ParsePath parses the given path
func (e *Executor) ParsePath(inputPath string) error {
	// check if input is a directory
	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("failed to stat input: %w", err)
	}

	if !info.IsDir() {
		content, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		if _, err := e.New(inputPath).Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}
		return nil
	}

	err = filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		// Read, parse, and execute template only if it's a file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		if _, err := e.New(path).Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	return nil
}

// Render executes the template with the given data
func (e Executor) Render(name string, data interface{}) (string, error) {
	var buf strings.Builder
	if err := e.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}
	return buf.String(), nil
}

// StringRender renders the given template string with the given data
func (e Executor) StringRender(templateString string, data interface{}) (string, error) {
	if _, err := e.New("temp").Parse(templateString); err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf strings.Builder
	if err := e.ExecuteTemplate(&buf, "temp", data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}
