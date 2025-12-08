package walker

import (
	"path/filepath"
	"strings"
)

// Walk traverses the directory tree starting from dir, applying include and exclude patterns
// It returns a list of file paths that match the criteria
func Walk(dir string, include []string, exclude []string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d filepath.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// Check exclude patterns first
		if matchesAny(relPath, exclude) {
			return nil
		}

		// Check include patterns
		if len(include) == 0 || matchesAny(relPath, include) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// matchesAny checks if the path matches any of the patterns
func matchesAny(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
		// TODO: Support glob patterns like **/*.go
		if strings.Contains(pattern, "**") {
			// Simple glob support
			if strings.HasPrefix(pattern, "**/") {
				suffix := strings.TrimPrefix(pattern, "**/")
				if strings.HasSuffix(path, suffix) {
					return true
				}
			}
		}
	}
	return false
}