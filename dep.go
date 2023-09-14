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

// String returns the string representation of a DepInfo.
func (d DepInfo) String() string {
	var builder strings.Builder

	builder.WriteString(d.Host)
	builder.WriteRune('/')
	builder.WriteString(d.Owner)
	builder.WriteRune('/')
	builder.WriteString(d.Repo)

	if d.Path != "" {
		builder.WriteRune('/')
		builder.WriteString(d.Path)
	}

	if d.Block != "" {
		builder.WriteRune('#')
		builder.WriteString(d.Block)
	}

	if d.Tag != "" {
		builder.WriteRune('@')
		builder.WriteString(d.Tag)
	}

	return builder.String()
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

	// Extract the tag if it exists before splitting the path
	fullPath := strings.Trim(u.Path, "/")
	tag := ""
	if idx := strings.Index(fullPath, "@"); idx != -1 {
		fullPath, tag = splitAtSign(fullPath)
	}

	// Split path into components
	pathParts := strings.Split(fullPath, "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid path format in embed URL")
	}

	owner := pathParts[0]
	repo := pathParts[1]
	path := ""
	if len(pathParts) > 2 {
		path = strings.Join(pathParts[2:], "/")
	}

	// If fragment contains the tag, then prioritize it over the tag in the path
	block, fragmentTag := extractBlockAndTag(u.Fragment)
	if fragmentTag != "" {
		tag = fragmentTag
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
	DefaultBranch() string
}

// DefaultGitClient provides a default implementation for the GitClient interface.
type DefaultGitClient struct {
	Token         string
	defaultBranch string
}

// NewDefaultGitClient creates a new DefaultGitClient with the given token.
func NewDefaultGitClient(defaultBranch string, token string) *DefaultGitClient {
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	return &DefaultGitClient{
		Token:         token,
		defaultBranch: defaultBranch,
	}
}

// DefaultBranch returns the default branch name.
func (d *DefaultGitClient) DefaultBranch() string {
	return d.defaultBranch
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

// Checkout checks out a branch, tag or commit hash in a Git repository.
func (d *DefaultGitClient) Checkout(path, ref string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Try to checkout branch
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/remotes/origin/" + ref),
	})
	if err == nil {
		return nil // Branch checked out successfully
	}

	// If branch checkout fails, try tag
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + ref),
	})
	if err == nil {
		return nil // Tag checked out successfully
	}

	// If tag checkout also fails, try using it as a commit hash
	commitHash := plumbing.NewHash(ref)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: commitHash,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout reference %s: %w", ref, err)
	}

	return nil
}
