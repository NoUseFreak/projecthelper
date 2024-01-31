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

var cloneForks bool

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
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parts := strings.Split(args[0], "/")

			if len(parts) < 2 {
				logrus.Fatal("Invalid repo")
			}

			var provider org.OrgProvider
			switch parts[0] {
			case "github.com":
				provider = &org.GithubProvider{
					Org: parts[1],
				}
			case "gitlab.com":
				provider = &org.GitlabProvider{
					User: strings.Join(parts[1:], "/"),
				}
			case "dev.azure.com":
				provider = &org.AzureProvider{
					Org: parts[1],
				}
            case "bitbucket.org":
                provider = &org.BitbucketProvider{
                    Org: parts[1],
                }
			default:
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
				if _, err := cloneRepo(repo.SSHURL); err != nil {
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

	return cmd
}
