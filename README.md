# Project Helper

Project Helper helps you structure your projects on your filesystem.

## Install

```bash
# Download the binary
go install github.com/nousefreak/projecthelper@latest
# Install the ph alias
projecthelper install
```

## Commands

| command                                                | description                                               |
| ------------------------------------------------------ | --------------------------------------------------------- | 
| `ph setup`                                             | Run setup (manages `.config/projecthelper/config.yaml`    |
| `ph clone https://github.com/nousefreak/projecthelper` | Clone to `${basedir}/github.com/nousefreak/projecthelper` |
| `ph org github.com/nousefreak`                         | Clones all repos (set `GITHUB_TOKEN` to include private)  |
| `ph go [search]`                                       | Open a fuzzyfinder that will `cd` to the repo             |
| `ph [search]`                                          | Shorthand for `ph go`                                     |
| `ph update`                                            | Run a `git fetch` on all repos                            |


