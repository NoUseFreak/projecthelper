package org

import (
	"fmt"
	"os"

	"github.com/ktrysmt/go-bitbucket"
)

type BitbucketProvider struct {
	Org string
}

func (b *BitbucketProvider) GetRepos() ([]*Repo, bool, error) {
	client := bitbucket.NewBasicAuth(os.Getenv("BITBUCKET_USERNAME"), os.Getenv("BITBUCKET_PASSWORD"))
	repos, err := client.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
		Owner: b.Org,
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to list repositories: %w", err)
	}

	var result []*Repo
	for _, repo := range repos.Items {
		result = append(result, &Repo{
			Name:     repo.Name,
			URL:      getBitbucketLink(repo.Links, "html"),
			CloneURL: getBitbucketLink(repo.Links, "clone", "https"),
			SSHURL:   getBitbucketLink(repo.Links, "clone", "ssh"),
		})
	}

	return result, false, nil
}

func getBitbucketLink(links map[string]interface{}, p ...string) string {
	for name, link := range links {
		if name != p[0] {
			continue
		}
		if l, ok := link.(map[string]interface{}); ok {
			return l["href"].(string)
		}
		if l, ok := link.([]interface{}); ok {
			for _, ll := range l {
				if llm, ok := ll.(map[string]interface{}); ok {
					if llm["name"].(string) != p[1] {
						continue
					}
					return llm["href"].(string)
				}
			}
		}
	}

	return ""
}
