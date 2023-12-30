package command

import (
	"fmt"
	"strings"
	"sync"

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
            var mut sync.RWMutex
            paths := []string{}
			go func(paths *[]string) {
                err := repo.GetRepoPathsAsync(baseDir, paths)
                if err != nil {
                    logrus.Fatal(fmt.Errorf("failed to get repo paths: %w", err))
                }
            }(&paths)

			idx, err := fuzzyfinder.Find(
				&paths,
				func(i int) string {
					return paths[i]
				},
				fuzzyfinder.WithQuery(strings.Join(args, " ")),
				fuzzyfinder.WithSelectOne(),
                fuzzyfinder.WithHotReloadLock(mut.RLocker()),
			)
			switch err {
			case nil:
				logrus.Infof("Jumping to %s", paths[idx])
				fmt.Fprintf(CmdOutput, "cd %s\n", paths[idx])
			case fuzzyfinder.ErrAbort:
				logrus.Fatal("aborted")
			default:
				logrus.Fatal(fmt.Errorf("failed to find repo: %w", err))
			}
		},
	}
	return cmd
}

