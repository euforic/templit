package templit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// WalkAndProcessDirFunc is called for each file and directory when walking a directory.
type WalkAndProcessDirFunc func(path string, isDir bool, content string) error

// WalkAndProcessDir processes all files in a directory with the given data.
// If walkFunc is provided, it's called for each file and directory without writing the file to disk.
func WalkAndProcessDir(inputDir, outputDir string, funcMap template.FuncMap, data interface{}) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Master template for file processing
	masterTmpl := template.New("main").Funcs(funcMap)

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through directory: %w", err)
		}

		parsedName, err := RenderTemplate(filepath.Base(path), data, funcMap)
		if err != nil {
			return fmt.Errorf("error rendering path template: %w", err)
		}

		relPath, err := filepath.Rel(inputDir, filepath.Dir(path))
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}

		outPath := filepath.Join(outputDir, relPath, parsedName)
		parsedOutPath, err := RenderTemplate(outPath, data, funcMap)
		if err != nil {
			return fmt.Errorf("error rendering path template: %w", err)
		}

		if info.IsDir() {
			// Skip directories with empty or "-" prefixed names
			if parsedName == "" || strings.HasPrefix(parsedName, "-") {
				return filepath.SkipDir
			}

			// Skip root directory
			if filepath.Base(outPath) == filepath.Base(inputDir) {
				return nil
			}

			if err := os.MkdirAll(parsedOutPath, info.Mode()); err != nil {
				return fmt.Errorf("error creating directory: %w", err)
			}

			return nil
		}

		// Skip files with empty names or "-" prefixed
		if parsedName == "" || strings.HasPrefix(parsedName, "-") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file from templates: %w", err)
		}

		tmpl, err := masterTmpl.New(path).Parse(string(content))
		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		var buf strings.Builder
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("error executing template: %w", err)
		}

		if err := os.WriteFile(parsedOutPath, []byte(buf.String()), info.Mode()); err != nil {
			return fmt.Errorf("error writing file to output: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through directory: %w", err)
	}

	return nil
}

// RenderTemplate renders a template with provided data.
func RenderTemplate(tmpl string, data interface{}, funcMap template.FuncMap) (string, error) {
	t, err := template.New("main").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}
