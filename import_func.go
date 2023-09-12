package templit

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// ImportFunc returns a function that can be used as a template function to import and process a template from a remote git repository.
// ImportFunc allows embedding content from a remote repository into a Go template.
//
// Steps to use:
//  1. Add the function to the FuncMap.
//  2. Use the following syntax within your template:
//     `{{ import "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" "<path_to_genrate_files>" . }}`
//
// Placeholders:
//   - `<host>`: Repository hosting service (e.g., "github.com").
//   - `<owner>`: Repository owner or organization.
//   - `<repo>`: Repository name.
//   - `<path>`: Path to the desired file or directory within the repository.
//   - `<tag_or_hash_or_branch>`: Specific Git reference (tag, commit hash, or branch name).
func ImportFunc(client GitClient) func(repoAndTag, destPath string, data interface{}, funcMap template.FuncMap) (string, error) {
	return func(repoAndTag, destPath string, data interface{}, funcMap template.FuncMap) (string, error) {
		const tempDirPrefix = "temp_clone_"

		depInfo, err := ParseDepURL(repoAndTag)
		if err != nil {
			return "", fmt.Errorf("failed to parse embed URL: %w", err)
		}

		if depInfo.Tag == "" {
			depInfo.Tag = "main"
		}

		tempDir, err := os.MkdirTemp("", tempDirPrefix)
		if err != nil {
			return "", fmt.Errorf("failed to create temp dir: %w", err)
		}
		defer os.RemoveAll(tempDir) // Cleanup

		if err := client.Clone(depInfo.Host, depInfo.Owner, depInfo.Repo, tempDir); err != nil {
			return "", fmt.Errorf("failed to clone repo: %w", err)
		}

		err = client.Checkout(tempDir, depInfo.Tag)
		if err != nil {
			return "", fmt.Errorf("failed to checkout branch: %w", err)
		}

		sourcePath := filepath.Join(tempDir, depInfo.Path)
		err = WalkAndProcessDir(sourcePath, destPath, funcMap, data)
		if err != nil {
			return "", fmt.Errorf("failed to process template: %w", err)
		}

		return "", nil
	}
}
