package command

import (
	"os"

	"github.com/nousefreak/projecthelper/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdName = "ph"
)

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   cmdName,
		Short: "project helper",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getGoCmd().Run(cmd, args)
		},
	}
	return rootCmd
}

func Execute() {
	rootCmd := getRootCmd()
	rootCmd.DisableSuggestions = true

	rootCmd.AddCommand(getCloneCmd())
	rootCmd.AddCommand(getGoCmd())
	rootCmd.AddCommand(getInstallCmd())
	rootCmd.AddCommand(getSetupCmd())
    rootCmd.AddCommand(getUpdateCmd())
    rootCmd.AddCommand(getVersionCmd())
    rootCmd.AddCommand(getOrgCmd())
    rootCmd.AddCommand(getWDIDCmd())

	cobra.OnInitialize(config.InitConfig)

	rootCmd.SetOut(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
