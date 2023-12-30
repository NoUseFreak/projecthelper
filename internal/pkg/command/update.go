package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nousefreak/projecthelper/internal/pkg/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getUpdateCmd() *cobra.Command {
    updateCmd := &cobra.Command{
        Use:   "update",
        Short: "update command",
        Long:  `update command`,
        Run: func(cmd *cobra.Command, args []string) {
            baseDir := getBaseDir()
            repoPaths, err := repo.GetRepoPaths(baseDir)
            if err != nil {
                logrus.Fatal(fmt.Errorf("failed to get repo paths: %w", err))
            }

            cmds := []string{}
            for _, repoPath := range repoPaths {
                cmds = append(cmds, fmt.Sprintf("(echo \"Updating %s\" && cd %s && git fetch -q)", repoPath, filepath.Join(baseDir, repoPath)))
            }

            fmt.Fprint(CmdOutput, strings.Join(cmds, " ; "))
        },   
    }
    return updateCmd
}

