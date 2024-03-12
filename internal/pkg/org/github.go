package org

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/sirupsen/logrus"
)

type GithubProvider struct {
	Org string
}

func (GithubProvider) FromURL(url string, typeHint string) (OrgProvider, error) {
	parts := strings.Split(url, "/")

	if parts[0] == "github.com" {
		if len(parts) < 2 {
			return nil, fmt.Errorf("Invalid repo")
		}
		return &GithubProvider{
			Org: parts[1],
		}, nil
	}

	return nil, nil
}

func (g *GithubProvider) GetRepos() ([]*Repo, bool, error) {
	repos := []*github.Repository{}
	repoFunc := getRepoFunc(g.Org)
	for {
		reposPage, done, err := repoFunc()
		if err != nil {
			return nil, false, fmt.Errorf("failed to list repositories: %w", err)
		}
		repos = append(repos, reposPage...)

		if done {
			break
		}
	}

	result := []*Repo{}
	for _, repo := range repos {
		if strings.EqualFold(repo.Owner.GetLogin(), g.Org) {
			result = append(result, &Repo{
				Name:     repo.GetName(),
				URL:      repo.GetHTMLURL(),
				CloneURL: repo.GetCloneURL(),
				SSHURL:   repo.GetSSHURL(),
			})
		}
	}

	return result, true, nil
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
