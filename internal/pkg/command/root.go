package command

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdName = "ph"
	cfgFile string
)

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ph",
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

	cobra.OnInitialize(initConfig)

	rootCmd.SetOut(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
