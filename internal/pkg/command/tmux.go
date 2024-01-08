package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getTmuxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tmux [search]",
		Aliases: []string{"cd"},
		Short:   "tmux session of result",
		Run: func(cmd *cobra.Command, args []string) {
			path, err := findPath(strings.Join(args, " "))
			name := strings.ReplaceAll(filepath.Base(path), ".", "_")

			if err != nil {
				if err == fuzzyfinder.ErrAbort {
					logrus.Fatal("aborted")
				} else {
					logrus.Fatal(fmt.Errorf("failed to find repo: %w", err))
				}
				os.Exit(1)

			}

			isRunning := exec.Command("pgrep", "tmux").Run() == nil

			if !isRunning {
				logrus.Infof("Starting tmux with new session '%s' %s", name, path)
				fmt.Fprintf(CmdOutput, "tmux new-session -s %s -c %s; ", name, path)
				os.Exit(0)
			}

			hasSession := exec.Command("tmux", "has-session", "-t", name).Run() == nil

			if !hasSession {
				logrus.Infof("Create new session '%s' %s", name, path)
				fmt.Fprintf(CmdOutput, "tmux new-session -ds %s -c %s; ", name, path)
			}

			logrus.Infof("Switch to session '%s' %s", name, path)
			fmt.Fprintf(CmdOutput, "tmux switch-client -t %s; ", name)
		},
	}
	return cmd
}
