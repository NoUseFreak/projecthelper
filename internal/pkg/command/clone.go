package command

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
    "os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	giturls "github.com/whilp/git-urls"
)

func getCloneCmd() *cobra.Command {
	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "clone command",
		Long:  `clone command`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoURL, err := giturls.Parse(args[0])
			if err != nil {
				logrus.Fatal(fmt.Sprintf("Error cleaning repo URL: %s", err))
			}
			repoPath, err := makePath(repoURL)
			if err != nil {
				logrus.Fatal(fmt.Sprintf("Error parsing repo URL: %s", err))
			}
			gitURL, err := makeURL(repoURL)
			if err != nil {
				logrus.Fatal(fmt.Sprintf("Error making git URL: %s", err))
			}
			baseDir := viper.GetString("basedir")
			if baseDir == "" {
				logrus.Fatal("Basedir not set. Run `ph setup` to set it.")
			}

			targetDir := strings.ToLower(filepath.Join(baseDir, repoPath))
        
            if stat, err := os.Stat(targetDir); err == nil && stat.IsDir() {
                // check if directory is empty
                if _, err := os.ReadDir(targetDir); err == nil {
                    logrus.Fatalf("Directory %s already exists and is not empty", targetDir)
                }
            }

			logrus.Infof("Cloning %s into %s", gitURL, targetDir)
			fmt.Fprintf(CmdOutput, "git clone %s %s && cd %s\n", gitURL, targetDir, targetDir)
		},
	}
	return cloneCmd
}

func makeURL(u *url.URL) (string, error) {
	hostRename := [][]string{
		{"github.com/org", "gh-repo"},
	}

	for _, set := range hostRename {
		r := regexp.MustCompile(regexp.QuoteMeta(set[0]))
		if r.MatchString(u.String()) {
			u.Host = set[1]
		}
	}

	return u.String(), nil
}

func makePath(u *url.URL) (string, error) {
	return fmt.Sprintf("%s/%s", u.Hostname(), strings.TrimSuffix(u.Path, ".git")), nil
}
