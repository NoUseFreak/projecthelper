package repo

import (
	"fmt"
	"os"
)

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
