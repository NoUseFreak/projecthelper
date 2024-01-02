package repo

import (
	"fmt"
	"os"
	"strings"
)

func GetRepoPathsAsync(baseDir string, result *[]string) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

    subdirs := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if entry.Name() == ".git" {
			*result = append(*result, baseDir)
			return nil
		}
        
        subdirs = append(subdirs, entry.Name())
	}

    for _, subdir := range subdirs {
        err := GetRepoPathsAsync(fmt.Sprintf("%s/%s", baseDir, subdir), result)
        if err != nil {
            return err
        }
    }

	return nil
}

func GetRepoPaths(baseDir string) ([]string, error) {
	result := []string{}
	err := GetRepoPathsAsync(baseDir, &result)

	return result, err
}

func FilterRepoPaths(paths []string, filter string) []string {
	result := []string{}
	for _, path := range paths {
		if strings.Contains(path, filter) {
			result = append(result, path)
		}
	}
	return result
}
