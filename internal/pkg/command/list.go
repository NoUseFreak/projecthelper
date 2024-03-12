package command

import (
	"fmt"
	"github.com/nousefreak/projecthelper/internal/pkg/config"
	"github.com/nousefreak/projecthelper/internal/pkg/repo"
	"github.com/spf13/cobra"
)

func getListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "list_fzf",
		Short:  "List all repositories",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			baseDir := config.GetBaseDir()
			repos := repo.GetRepoPathsChan(baseDir, true)
			for repo := range repos {
				fmt.Println(repo)
			}
		},
	}
	return cmd
}
