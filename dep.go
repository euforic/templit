package templit

import (
	"fmt"
	"net/url"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// DepInfo contains information about an embed URL.
type DepInfo struct {
	Host  string
	Owner string
	Repo  string
	Path  string
	Block string
	Tag   string
}

// ParseDepURL is a parsed embed URL.
func ParseDepURL(rawURL string) (*DepInfo, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Split path into components
	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid path format in embed URL")
	}

	owner := pathParts[0]
	repo := pathParts[1]
	path := ""
	if len(pathParts) > 2 {
		path = strings.Join(pathParts[2:], "/")
	}

	// Extract block and version
	block, tag := extractBlockAndTag(u.Fragment)

	if tag == "" {
		if strings.Contains(path, "@") {
			path, tag = splitAtSign(path)
		} else if strings.Contains(repo, "@") {
			repo, tag = splitAtSign(repo)
		}
	}

	return &DepInfo{
		Host:  u.Host,
		Owner: owner,
		Repo:  repo,
		Path:  path,
		Block: block,
		Tag:   tag,
	}, nil
}

// splitAtSign splits the given string at the '@' sign and returns both parts.
func splitAtSign(s string) (string, string) {
	parts := strings.Split(s, "@")
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

// extractBlockAndTag extracts the block and tag from a fragment.
func extractBlockAndTag(fragment string) (string, string) {
	parts := strings.Split(fragment, "@")
	if len(parts) > 1 {
		return parts[0], parts[1]
	} else if len(parts) == 1 {
		return parts[0], ""
	}
	return "", ""
}

// GitClient is an interface that abstracts Git operations.
type GitClient interface {
	Clone(host, owner, repo, dest string) error
	Checkout(path, branch string) error
}

// DefaultGitClient provides a default implementation for the GitClient interface.
type DefaultGitClient struct {
	Token   string
	BaseURL string
}

// NewDefaultGitClient creates a new DefaultGitClient with the given token.
func NewDefaultGitClient(token string) *DefaultGitClient {
	return &DefaultGitClient{
		Token: token,
	}
}

// Clone clones a Git repository to the given destination.
func (d *DefaultGitClient) Clone(host, owner, repo, dest string) error {
	repoURL := fmt.Sprintf("%s/%s/%s.git", host, owner, repo)
	if !strings.HasPrefix(repoURL, "https://") {
		repoURL = fmt.Sprintf("https://%s", repoURL)
	}

	var auth *http.BasicAuth
	if d.Token != "" {
		auth = &http.BasicAuth{
			Username: "username", // this can be anything except an empty string
			Password: d.Token,
		}
	}

	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:  repoURL,
		Auth: auth,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repo %s: %w", repoURL, err)
	}

	return nil
}

// Checkout checks out a branch in a Git repository.
func (d *DefaultGitClient) Checkout(path, branch string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
}
