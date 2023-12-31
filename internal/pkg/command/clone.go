package command

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nousefreak/projecthelper/internal/pkg/config"
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
            dir, err := cloneRepo(args[0])
            if ; err != nil {
                logrus.Fatal(err)
            }

            fmt.Fprintf(CmdOutput, "cd %s \n", dir)
		},
	}
	return cloneCmd
}

var ErrDirectoryAlreadyExists = fmt.Errorf("Directory already exists")

func cloneRepo(repo string) (string, error) {
			baseDir := config.GetBaseDir()
			repoURL, err := giturls.Parse(repo)
			if err != nil {
				return "", fmt.Errorf("Error cleaning repo URL: %w", err)
			}
			repoPath, err := makePath(repoURL)
			if err != nil {
				return "", fmt.Errorf("Error parsing repo URL: %s", err)
			}
			gitURL, err := makeURL(repoURL, viper.GetStringMapString("renameRepo"))
			if err != nil {
			    return "", fmt.Errorf("Error making git URL: %s", err)
			}
			targetDir := strings.ToLower(filepath.Join(baseDir, repoPath))

			if stat, err := os.Stat(targetDir); err == nil && stat.IsDir() {
				// check if directory is empty
				if _, err := os.ReadDir(targetDir); err == nil {
                    return "", fmt.Errorf("Directory %s already exists and is not empty: %w", targetDir, ErrDirectoryAlreadyExists)
				}
			}

			fmt.Fprintf(CmdOutput, "(echo \"\\033[0;32m*\\033[0m Cloning %s into %s\" && git clone %s %s) \n", gitURL, targetDir, gitURL, targetDir)

            return targetDir, nil
        }


func makeURL(u *url.URL, renameRepo map[string]string) (string, error) {
	for match, host := range renameRepo {
		r := regexp.MustCompile(regexp.QuoteMeta(match))
		if r.MatchString(u.String()) {
			u.Host = host
		}
	}

	return u.String(), nil
}

func makePath(u *url.URL) (string, error) {
	return fmt.Sprintf("%s/%s", u.Hostname(), strings.TrimSuffix(u.Path, ".git")), nil
}
