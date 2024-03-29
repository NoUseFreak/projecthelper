package repo

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home + path[1:]
	}
	return os.ExpandEnv(path)
}

func extraDirs() []string {
	paths := viper.GetStringSlice("extraDirs")
	for i, path := range paths {
		paths[i] = ExpandPath(path)
	}

	return paths
}

func GetRepoPathsChan(basedir string, includeExtras bool) <-chan string {
	basedir = ExpandPath(basedir)
	out := make(chan string)
	go func() {
		defer close(out)
		if includeExtras {
			for _, p := range extraDirs() {
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
			if inExcludeList(basedir) {
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

func inExcludeList(name string) bool {
	excludes := viper.GetStringSlice("excludeDirs")
	for _, exclude := range excludes {
		if name == ExpandPath(exclude) {
			return true
		}
	}
	return false
}

func GetRepoPathsAsync(baseDir string, result *[]string) error {
	ch := GetRepoPathsChan(baseDir, true)
	for p := range ch {
		*result = append(*result, p)
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
