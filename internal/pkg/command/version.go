package command

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)


var (
	Version string
	Commit  string
	Date    string
)

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
				Version = info.Main.Version
				Commit = info.Main.Sum
			}
			logrus.Infof("version: \t %s", Version)
			logrus.Infof("commit: \t %s", Commit)
			logrus.Infof("date: \t %s", Date)
		},
	}
}

