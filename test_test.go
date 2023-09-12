package templit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type MockGitClient struct{}

func (m *MockGitClient) Clone(host, owner, repo, dest string) error {
	src := filepath.Join(host, owner, repo)
	return copyDir(src, dest)
}

func (m *MockGitClient) Checkout(path, branch string) error {
	// Mocked function, doesn't need to do anything for this test.
	return nil
}

// copyDir copies a directory recursively
func copyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
			err = copyDir(sourcePath, destPath)
			if err != nil {
				return err
			}
			continue
		}

		if err := copyFile(sourcePath, destPath); err != nil {
			return err
		}
	}
	return nil
}

// copyFile copies a file
func copyFile(src, dst string) error {
	bytes, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, bytes, 0644)
}

// compareDirs recursively compares the contents of two directories.
// dir1 is the expected directory, while dir2 is the actual directory.
func compareDirs(dir1, dir2 string) error {
	entries1, err := os.ReadDir(dir1)
	if err != nil {
		return err
	}

	entries2, err := os.ReadDir(dir2)
	if err != nil {
		return err
	}

	entryMap2 := make(map[string]os.DirEntry)
	for _, entry := range entries2 {
		entryMap2[entry.Name()] = entry
	}

	for _, entry1 := range entries1 {
		entry2, exists := entryMap2[entry1.Name()]

		if !exists && strings.HasPrefix(entry1.Name(), "-") {
			continue
		}

		if !exists {
			entryType := "file"
			if entry1.IsDir() {
				entryType = "directory"
			}

			return fmt.Errorf("missing %s in %s: %s", entryType, dir2, entry1.Name())
		}

		if entry1.IsDir() {
			if !entry2.IsDir() {
				return fmt.Errorf("expected file, but found directory in %s: %s", dir2, entry1.Name())
			}

			if err := compareDirs(filepath.Join(dir1, entry1.Name()), filepath.Join(dir2, entry2.Name())); err != nil {
				return err
			}
		}

		if entry1.Type().IsRegular() && entry2.Type().IsRegular() {
			bytes1, err := os.ReadFile(filepath.Join(dir1, entry1.Name()))
			if err != nil {
				return err
			}

			bytes2, err := os.ReadFile(filepath.Join(dir2, entry2.Name()))
			if err != nil {
				return err
			}

			if diff := cmp.Diff(string(bytes1), string(bytes2)); diff != "" {
				return fmt.Errorf("mismatch in files %s vs %s\n\t(-want +got):\n%s ", filepath.Join(dir1, entry1.Name()), filepath.Join(dir2, entry2.Name()), diff)
			}
			continue
		}

		if entry1.Type().IsRegular() || entry2.Type().IsRegular() {
			return fmt.Errorf("file type mismatch for %s in %s vs %s", entry1.Name(), dir1, dir2)
		}
	}

	// Check for any extra entries in dir2 that are not present in dir1
	for _, entry2 := range entries2 {
		if _, exists := entryMap2[entry2.Name()]; !exists {
			return fmt.Errorf("unexpected entry in %s: %s", dir2, entry2.Name())
		}
	}

	return nil
}
