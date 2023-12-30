package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/nousefreak/projecthelper/internal/pkg/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getGoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go",
		Short: "go to a project",
		Run: func(cmd *cobra.Command, args []string) {
			baseDir := getBaseDir()
			if baseDir == "" {
				logrus.Fatal("Basedir not set. Run `ph setup` to set it.")
			}
			paths, err := repo.GetRepoPaths(baseDir)
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

