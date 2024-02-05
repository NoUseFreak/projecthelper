package repo

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func GetRepoPathsChan(basedir string, includeExtras bool) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		if includeExtras {
			for _, p := range viper.GetStringSlice("extraDirs") {
				out <- p
			}
		}

		entries, err := os.ReadDir(basedir)
		if err != nil {
			return
		}
		subdirs := []string{}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if entry.Name() == ".git" {
				out <- basedir
				return
			}
			subdirs = append(subdirs, entry.Name())
		}

		for _, subdir := range subdirs {
			subchan := GetRepoPathsChan(fmt.Sprintf("%s/%s", basedir, subdir), false)
			for p := range subchan {
				out <- p
			}
		}
	}()

	return out
}

func GetRepoPathsAsync(baseDir string, result *[]string) error {
	if len(*result) == 0 {
		*result = append(*result, viper.GetStringSlice("extraDirs")...)
	}

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
