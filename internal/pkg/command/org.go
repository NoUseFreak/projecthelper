package command

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nousefreak/projecthelper/internal/pkg/org"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cloneForks bool
	typeHint   string
)

func getOrgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org ORANIZATION_URL",
		Short: "Clone repositories from an organization into the project structure",
		Long: `Clone repositories from an organization into the project structure.

The ORGANIZATION_URL is the URL of the organization on the platform. For example:
- github.com/nousefreak
- gitlab.com/nousefreak
- dev.azure.com/nousefreak
- bitbucket.org/nousefreak

In addidtion to the ORGANIZATION_URL you should also specify credentials.
The credentials are read from the environment variables:
- GITHUB_TOKEN
- GITLAB_TOKEN
- AZURE_TOKEN
- BITBUCKET_USERNAME, BITBUCKET_PASSWORD
        `,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parts := strings.Split(args[0], "/")

			if len(parts) < 2 {
				logrus.Fatal("Invalid repo")
			}

			providers := []org.OrgProvider{
				&org.AzureProvider{},
				&org.GithubProvider{},
				&org.GitlabProvider{},
				&org.BitbucketProvider{},
			}

			provider, err := org.GetProviderFromURL(providers, args[0],	typeHint)
			if err != nil {
				logrus.Fatal(err)
			}

			if provider == nil {
				logrus.Fatal("Unsupported platform")
			}

			repos, found, err := provider.GetRepos()
			if err != nil {
				logrus.Fatal(fmt.Errorf("failed to get repos: %w", err))
			}

			if !found {
				logrus.Warn("No repos found")
			}

			logrus.Infof("Found %d repos", len(repos))

			for _, repo := range repos {
                var cloneURL string
                if cloneProtocol == "https" {
                    cloneURL = repo.CloneURL
                } else {
                    cloneURL = repo.SSHURL
                }

				if _, err := cloneRepo(cloneURL); err != nil {
					if err == ErrDirectoryAlreadyExists || errors.Unwrap(err) == ErrDirectoryAlreadyExists {
						logrus.Debugf("Skipping existing repo: %s", repo.Name)
					} else {
						logrus.Warn(err)
					}
				} else {
					fmt.Fprint(os.Stdout, " ; ")
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&cloneForks, "forks", "f", false, "Clone forks")
	cmd.Flags().StringVarP(&typeHint, "type-hint", "t", "", "Add a type hint to the URL to force a specific provider")
    cmd.Flags().StringVarP(&cloneProtocol, "clone-protocol", "p", "ssh", "Clone protocol (ssh, https)")

	return cmd
}
