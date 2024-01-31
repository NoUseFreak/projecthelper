package org

type OrgProvider interface {
    GetRepos() ([]*Repo, bool, error)
}

type Repo struct {
    Name string
    URL string
    CloneURL string
    SSHURL string
}

