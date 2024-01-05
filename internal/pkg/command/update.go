package command

import (
	"fmt"
	"strings"

	"github.com/nousefreak/projecthelper/internal/pkg/config"
	"github.com/nousefreak/projecthelper/internal/pkg/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getUpdateCmd() *cobra.Command {
    updateCmd := &cobra.Command{
        Use:   "update",
        Short: "Run git fetch on all repos",
        Run: func(cmd *cobra.Command, args []string) {
            baseDir := config.GetBaseDir()
            repoPaths, err := repo.GetRepoPaths(baseDir)
            if err != nil {
                logrus.Fatal(fmt.Errorf("failed to get repo paths: %w", err))
            }

            if len(args) > 0 {
                repoPaths = repo.FilterRepoPaths(repoPaths, args[0])
            }

            cmds := []string{}
            for _, repoPath := range repoPaths {
                cmds = append(cmds, fmt.Sprintf("(echo \"\\033[0;32m*\\033[0m Updating %s\" && git -C %s fetch -q)", repoPath, repoPath))
            }

            fmt.Fprint(CmdOutput, strings.Join(cmds, " ; "))
        },
    }
    return updateCmd
}

