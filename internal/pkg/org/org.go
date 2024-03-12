package org

type OrgProvider interface {
	GetRepos() ([]*Repo, bool, error)
	FromURL(url string, typeHint string) (OrgProvider, error)
}

type Repo struct {
	Name     string
	URL      string
	CloneURL string
	SSHURL   string
}

func GetProviderFromURL(providers []OrgProvider, url string, typeHint string) (OrgProvider, error) {
	for _, provider := range providers {
		p, err := provider.FromURL(url, typeHint)
		if err != nil {
			return nil, err
		}
		if p != nil {
			return p, nil
		}
	}

	return nil, nil
}
