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

// NewExecutor creates a new Executor with the given template and funcMap
func NewExecutor(inputPath string, funcMap template.FuncMap) (*Executor, error) {
	if funcMap == nil {
		funcMap = DefaultFuncMap
	}

	mainTmpl := template.New("main").Funcs(funcMap)

	// check if input is a directory
	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat input: %w", err)
	}

	if !info.IsDir() {
		content, err := os.ReadFile(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		tmpl, err := mainTmpl.New(inputPath).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}
		return &Executor{Template: tmpl}, nil
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

		if _, err := mainTmpl.New(path).Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Executor{
		Template: mainTmpl,
	}, nil
}

// Render executes the template with the given data
func (e Executor) Render(name string, data interface{}) (string, error) {
	var buf strings.Builder
	if err := e.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}
	return buf.String(), nil
}
