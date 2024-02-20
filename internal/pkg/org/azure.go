package org

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

type AzureProvider struct {
	Org string
}

func (AzureProvider) FromURL(url string, typeHint string) (OrgProvider, error) {
	parts := strings.Split(url, "/")

	if parts[0] == "dev.azure.com" {
		if len(parts) < 2 {
			return nil, fmt.Errorf("Invalid repo")
		}
		return &AzureProvider{
			Org: parts[1],
		}, nil
	}

	return nil, nil
}

func (a *AzureProvider) GetRepos() ([]*Repo, bool, error) {

	connection := azuredevops.NewPatConnection("https://dev.azure.com/"+a.Org, os.Getenv("AZURE_TOKEN"))
	ctx := context.Background()
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create core client: %w", err)
	}
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create core client: %w", err)
	}

	projects, err := getAzureProjects(ctx, coreClient)
	if err != nil {
		return nil, false, fmt.Errorf("failed to list projects: %w", err)
	}

	var repos []*Repo
	for _, project := range projects {
		r, err := gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{
			Project: &project,
		})
		if err != nil {
			return nil, false, fmt.Errorf("failed to list repos: %w", err)
		}

		for _, repo := range *r {
			repos = append(repos, &Repo{
				Name:     *repo.Name,
				URL:      *repo.RemoteUrl,
				CloneURL: *repo.RemoteUrl,
				SSHURL:   *repo.SshUrl,
			})
		}
	}

	return repos, false, nil
}

func getAzureProjects(ctx context.Context, client core.Client) ([]string, error) {
	resp, err := client.GetProjects(ctx, core.GetProjectsArgs{})
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	var projects []string
	for {
		if resp == nil {
			break
		}
		for _, project := range (*resp).Value {
			projects = append(projects, *project.Name)
		}
		if resp.ContinuationToken != "" {
			ct, err := strconv.Atoi(resp.ContinuationToken)
			if err != nil {
				return nil, fmt.Errorf("failed to parse continuation token: %w", err)
			}
			resp, err = client.GetProjects(ctx, core.GetProjectsArgs{
				ContinuationToken: &ct,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to list projects: %w", err)
			}
		} else {
			resp = nil
		}
	}

	return projects, nil
}
