package templit

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

// EmbedFunc returns a template function that can be used to process and embed a template from a remote git repository.
// EmbedFunc allows embedding content from a remote repository directly into a Go template.
//
// Steps to use:
//  1. Add the function to the FuncMap.
//  2. Use the following syntax within your template:
//     `{{ embed "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" . }}`
//     or
//     `{{ embed "<host>/<owner>/<repo>#<block>@<tag_or_hash_or_branch>" . }}`
//
// Placeholders:
//   - `<host>`: Repository hosting service (e.g., "github.com").
//   - `<owner>`: Repository owner or organization.
//   - `<repo>`: Repository name.
//   - `<path>`: Path to the desired file or directory within the repository.
//   - `<block>`: Specific template block name.
//   - `<tag_or_hash_or_branch>`: Specific Git reference (tag, commit hash, or branch name).
func EmbedFunc(client GitClient) func(remotePath string, data interface{}, funcMap template.FuncMap) (string, error) {
	return func(remotePath string, data interface{}, funcMap template.FuncMap) (string, error) {
		embedInfo, err := ParseDepURL(remotePath)
		if err != nil {
			return "", err
		}

		const tempDirPrefix = "temp_clone_"
		tempDir, err := os.MkdirTemp("", tempDirPrefix)
		if err != nil {
			return "", fmt.Errorf("failed to create temp dir: %w", err)
		}
		defer os.RemoveAll(tempDir) // Cleanup

		err = client.Clone(embedInfo.Host, embedInfo.Owner, embedInfo.Repo, tempDir)
		if err != nil {
			return "", fmt.Errorf("failed to clone repo: %w", err)
		}

		err = client.Checkout(tempDir, embedInfo.Tag)
		if err != nil {
			return "", fmt.Errorf("failed to checkout branch: %w", err)
		}

		executor, err := NewExecutor(tempDir, funcMap)
		if err != nil {
			return "", fmt.Errorf("failed to create executor: %w", err)
		}

		if embedInfo.Block != "" {
			return executor.Render(embedInfo.Block, data)
		}

		return executor.Render(path.Join(tempDir, embedInfo.Path), data)
	}
}
