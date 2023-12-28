package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getGoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go",
		Short: "go to a project",
		Run: func(cmd *cobra.Command, args []string) {
			baseDir := viper.GetString("basedir")
			if baseDir == "" {
				logrus.Fatal("Basedir not set. Run `ph setup` to set it.")
			}
			paths, err := getRepoPaths(baseDir)
			if err != nil {
				logrus.Fatal(fmt.Errorf("failed to get repo paths: %w", err))
			}
			idx, err := fuzzyfinder.Find(
				paths,
				func(i int) string {
					return paths[i]
				},
				fuzzyfinder.WithQuery(strings.Join(args, " ")),
				fuzzyfinder.WithSelectOne(),
			)
			switch err {
			case nil:
				logrus.Infof("Jumping to %s", paths[idx])
				fmt.Fprintf(CmdOutput, "cd %s\n", filepath.Join(baseDir, paths[idx]))
			case fuzzyfinder.ErrAbort:
				logrus.Fatal("aborted")
			default:
				logrus.Fatal(fmt.Errorf("failed to find repo: %w", err))
			}
		},
	}
	return cmd
}

func getRepoPaths(baseDir string) ([]string, error) {

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
