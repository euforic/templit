package templit

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestParseDepURL tests the ParseDepURL function.
func TestParseDepURL(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		expected *DepInfo
		wantErr  bool
	}{
		{
			name:   "Basic URL with path and fragment",
			rawURL: "https://github.com/owner/repo/path/to/file#block@v1.2.3",
			expected: &DepInfo{
				Host:  "github.com",
				Owner: "owner",
				Repo:  "repo",
				Path:  "path/to/file",
				Block: "block",
				Tag:   "v1.2.3",
			},
			wantErr: false,
		},
		{
			name:   "URL without path",
			rawURL: "https://github.com/owner/repo",
			expected: &DepInfo{
				Host:  "github.com",
				Owner: "owner",
				Repo:  "repo",
				Path:  "",
				Block: "",
				Tag:   "",
			},
			wantErr: false,
		},
		{
			name:   "URL with tag in path",
			rawURL: "https://github.com/owner/repo/path/to/file@v1.2.3",
			expected: &DepInfo{
				Host:  "github.com",
				Owner: "owner",
				Repo:  "repo",
				Path:  "path/to/file",
				Block: "",
				Tag:   "v1.2.3",
			},
			wantErr: false,
		},
		{
			name:   "URL with tag in repo",
			rawURL: "https://github.com/owner/repo@v1.2.3",
			expected: &DepInfo{
				Host:  "github.com",
				Owner: "owner",
				Repo:  "repo",
				Path:  "",
				Block: "",
				Tag:   "v1.2.3",
			},
			wantErr: false,
		},
		{
			name:   "URL without protocol",
			rawURL: "https://github.com/owner/repo/some/path#test_block@v1.2.3",
			expected: &DepInfo{
				Host:  "github.com",
				Owner: "owner",
				Repo:  "repo",
				Path:  "some/path",
				Block: "test_block",
				Tag:   "v1.2.3",
			},
			wantErr: false,
		},
		{
			name:     "Invalid URL missing repo name",
			rawURL:   "https://github.com/owner",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDepURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if err == nil {
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("result mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
