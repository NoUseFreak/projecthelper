package org

import (
	"fmt"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitlabProvider struct {
    User string
    Host string
}

func (GitlabProvider) FromURL(url string, typeHint string) (OrgProvider, error) {
    parts := strings.Split(url, "/")

    if parts[0] == "gitlab.com" || typeHint == "gitlab" {
        if len(parts) < 2 {
            return nil, fmt.Errorf("Invalid repo")
        }
        return &GitlabProvider{
            User: strings.Join(parts[1:], "/"),
            Host: fmt.Sprintf("https://%s", parts[0]),
        }, nil
    }

    return nil, nil
}

func (g *GitlabProvider) GetRepos() ([]*Repo, bool, error) {
    gl, err := gitlab.NewClient(
        os.Getenv("GITLAB_TOKEN"),
        gitlab.WithBaseURL(g.Host),
    )
    if err != nil {
        return nil, false, fmt.Errorf("failed to create gitlab client: %w", err)
    }

    opts := &gitlab.ListGroupProjectsOptions{
        IncludeSubGroups: gitlab.Ptr(true),
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

