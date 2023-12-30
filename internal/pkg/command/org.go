package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cloneForks bool

func getRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "repo",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parts := strings.Split(args[0], "/")

			if len(parts) < 2 {
				logrus.Fatal("Invalid repo")
			}

			switch parts[0] {
			case "github.com":
				if err := cloneGithubOrg(parts[1]); err != nil {
					logrus.Fatal(err)
				}

			default:
				logrus.Fatal("Unsupported platform")
			}

		},
	}

	cmd.Flags().BoolVarP(&cloneForks, "forks", "f", false, "Clone forks")

	return cmd
}

func getRepoFunc(org string) func() ([]*github.Repository, bool, error) {
	client := github.NewClient(nil)
	if os.Getenv("GITHUB_TOKEN") == "" {
		opts := &github.RepositoryListByUserOptions{
			Type: "all",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}
		return func() ([]*github.Repository, bool, error) {
			repos, resp, err := client.Repositories.ListByUser(context.Background(), org, opts)
			opts.Page = resp.NextPage
			return repos, resp.NextPage == 0, err
		}
	}

	logrus.Info("Using GITHUB_TOKEN")
	client = client.WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		Type: "all",
		Sort: "updated",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	return func() ([]*github.Repository, bool, error) {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(context.Background(), opts)
		opts.Page = resp.NextPage
		return repos, resp.NextPage == 0, err
	}
}

func cloneGithubOrg(org string) error {
	logrus.Infof("Cloning repos from github.com/%s", org)

	logrus.Info("Fetch repos")
	repos := []*github.Repository{}

	repoFunc := getRepoFunc(org)
	for {
		reposPage, done, err := repoFunc()
		if err != nil {
			return fmt.Errorf("Error fetching repos: %w", err)
		}
		repos = append(repos, reposPage...)

		if done {
			break
		}
	}

	orgRepos := []*github.Repository{}
	for _, repo := range repos {
		if strings.EqualFold(repo.Owner.GetLogin(), org) {
			orgRepos = append(orgRepos, repo)
		}
	}

	logrus.Infof("Found %d/%d repos", len(orgRepos), len(repos))

	for _, repo := range orgRepos {
		if *repo.Archived {
			logrus.Debugf("Skipping archived repo: %s", *repo.Name)
			continue
		}
		if cloneForks && *repo.Fork {
			logrus.Debugf("Skipping fork: %s", *repo.Name)
			continue
		}
		if _, err := cloneRepo(repo.GetSSHURL()); err != nil {
			if err == ErrDirectoryAlreadyExists || errors.Unwrap(err) == ErrDirectoryAlreadyExists {
				logrus.Debugf("Skipping existing repo: %s", *repo.Name)
			} else {
				logrus.Warn(err)
			}
		} else {
            fmt.Fprint(os.Stdout, " ; ")
        }
	}

	return nil
}
