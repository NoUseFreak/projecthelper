package command

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")	
	viper.AddConfigPath(home + "/.config/projecthelper")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn(err)
        cfgFile = home + "/.config/projecthelper/config.yaml"
		if err := os.MkdirAll(filepath.Dir(cfgFile), 0755); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		if err := viper.SafeWriteConfigAs(cfgFile); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		logrus.Warn("Created new config")
	}
}

