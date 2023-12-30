package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetRepoPaths(baseDir string) ([]string, error) {
	result := []string{}
	if err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}
		if !info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		if info.Name() != ".git" {
			return nil
		}

		result = append(result, filepath.Dir(relPath))

		return err
	}); err != nil {
		return nil, err
	}

	return result, nil
}
