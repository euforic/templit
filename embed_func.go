package templit

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// EmbedFunc returns a template function that can be used to process and embed a template from a remote git repository.
// EmbedFunc allows embedding content from a remote repository directly into a Go template.
//
// Steps to use:
//  1. Add the function to the FuncMap.
//  2. Use the following syntax within your template:
//     ```
//     {{ embed "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" . }}
//     {{ embed "<host>/<owner>/<repo>#<block>@<tag_or_hash_or_branch>" . }}
//     ```
//
// Placeholders:
//   - `<host>`: Repository hosting service (e.g., "github.com").
//   - `<owner>`: Repository owner or organization.
//   - `<repo>`: Repository name.
//   - `<path>`: Path to the desired file or directory within the repository.
//   - `<block>`: Specific template block name.
//   - `<tag_or_hash_or_branch>`: Specific Git reference (tag, commit hash, or branch name).
func (e *Executor) EmbedFunc(remotePath string, data interface{}) (string, error) {
	const tempDirPrefix = "templit_clone_"

	depInfo, err := ParseDepURL(remotePath)
	if err != nil {
		return "", err
	}

	if depInfo.Tag == "" {
		depInfo.Tag = e.git.DefaultBranch()
	}

	tempDir, err := os.MkdirTemp("", tempDirPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir) // Cleanup

	err = e.git.Clone(depInfo.Host, depInfo.Owner, depInfo.Repo, tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to clone repo: %w", err)
	}

	if depInfo.Tag != "" && depInfo.Tag != e.git.DefaultBranch() {
		err = e.git.Checkout(tempDir, depInfo.Tag)
		if err != nil {
			return "", fmt.Errorf("failed to checkout branch: %w", err)
		}
	}

	// templatePath is the path to the template file or directory
	templatePath := path.Join(tempDir, depInfo.Path)

	if err := e.ParsePath(filepath.Dir(templatePath)); err != nil {
		return "", fmt.Errorf("failed to create executor: %w", err)
	}

	if depInfo.Block != "" {
		return e.Render(depInfo.Block, data)
	}

	return e.Render(path.Join(tempDir, depInfo.Path), data)
}
