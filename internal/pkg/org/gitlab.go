package org

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
)

type GitlabProvider struct {
    User string
}

func (g *GitlabProvider) GetRepos() ([]*Repo, bool, error) {
    gl, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"))
    if err != nil {
        return nil, false, fmt.Errorf("failed to create gitlab client: %w", err)
    }

    opts := &gitlab.ListGroupProjectsOptions{
        IncludeSubGroups: gitlab.Bool(true),
        ListOptions: gitlab.ListOptions{
            PerPage: 100,
        },
    }

    var repos []*Repo
    for {
        projects, resp, err := gl.Groups.ListGroupProjects(g.User, opts)
        if err != nil {
            return nil, false, fmt.Errorf("failed to list projects: %w", err)
        }

        for _, project := range projects {
            repos = append(repos, &Repo{
                Name: project.Name,
                URL: project.WebURL,
                CloneURL: project.HTTPURLToRepo,
                SSHURL: project.SSHURLToRepo,
            })

            if resp.CurrentPage >= resp.TotalPages {
                return repos, true, nil
            }

            opts.Page = resp.NextPage
        }
    }
}

