package command

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/erikgeiser/promptkit/textinput"
)

func getSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup the environment",
		Long:  "Setup the environment",
		Run: func(cmd *cobra.Command, args []string) {
			home, err := os.UserHomeDir()
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("Setting up project helper")
			logrus.Info("Choose a base directory for project helper")

			input := textinput.New("Path:")
			input.InitialValue = filepath.Join(home, "src")
			input.Output = os.Stderr
			input.ResultTemplate = ""
			dir, err := input.RunPrompt()
			if err != nil {
				logrus.Fatal(err)
			}

			if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
				logrus.Infof("Creating directory %s", dir)
				if err := os.MkdirAll(dir, 0755); err != nil {
					logrus.Fatal(err)
				}
			}

			logrus.Infof("Setting up project helper in %s", dir)
			viper.Set("basedir", dir)
			if err := viper.WriteConfig(); err != nil {
				logrus.Fatal(err)
			}
		},
	}

	return cmd
}
