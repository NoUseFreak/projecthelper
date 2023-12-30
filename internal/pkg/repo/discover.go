package repo

import (
	"fmt"
	"os"
)

func GetRepoPathsAsync(baseDir string, result *[]string) error {
    entries, err := os.ReadDir(baseDir)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }

        if entry.Name() == ".git" {
            *result = append(*result, baseDir)
            return nil
        }

        err := GetRepoPathsAsync(fmt.Sprintf("%s/%s", baseDir, entry.Name()), result)
        if err != nil {
            return err
        }
    }

    return nil
}

func GetRepoPaths(baseDir string) ([]string, error) {
   result := []string{}

    entries, err := os.ReadDir(baseDir)
    if err != nil {
        return result, err
    }

    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }

        if entry.Name() == ".git" {
            return []string{baseDir}, nil
        }

        paths, err := GetRepoPaths(fmt.Sprintf("%s/%s", baseDir, entry.Name()))
        if err != nil {
            return result, err
        }

        result = append(result, paths...)
    }

    return result, nil
}
