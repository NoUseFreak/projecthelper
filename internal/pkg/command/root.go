package command

import (
	"io"
	"os"

	"github.com/nousefreak/projecthelper/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdName = "ph"

	quite bool
)

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   cmdName,
		Short: "Project helper",
		Long: `Project helper is a tool to help you manage and interact with your projects.

        It follows a simple convention of having all git projects cloned into a directory that reflects the git url.
        ex: git@github.com:nousefreak/projecthelper.git
            -> {basedir}/github.com/nousefreak/projecthelper

        `,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getGoCmd().Run(cmd, args)
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&quite, "quite", "q", true, "Hide output")
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
	rootCmd.AddCommand(getTmuxCmd())

	cobra.OnInitialize(config.InitConfig)

	rootCmd.SetOut(os.Stderr)

	if quite {
		logrus.SetOutput(io.Discard)
		rootCmd.SetOut(io.Discard)
	}

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
