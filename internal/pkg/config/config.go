package config

import (
	"os"
	"path/filepath"

	"github.com/nousefreak/projecthelper/internal/pkg/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ConfigFile string

func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home + "/.config/projecthelper")

	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn(err)
		ConfigFile = home + "/.config/projecthelper/config.yaml"
		if err := os.MkdirAll(filepath.Dir(ConfigFile), 0755); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		if err := viper.SafeWriteConfigAs(ConfigFile); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		logrus.Warn("Created new config")
	}
}

func GetBaseDir() string {
	baseDir := viper.GetString("basedir")
	if baseDir == "" {
		logrus.Fatal("Basedir not set. Run `ph setup` to set it.")
	}

	return repo.ExpandPath(baseDir)
}

