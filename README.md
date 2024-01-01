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

```bash
# Run setup (manages `.config/projecthelper/config.yaml`    
ph setup

# Clone to `${basedir}/github.com/nousefreak/projecthelper`
ph clone https://github.com/nousefreak/projecthelper

# Clones all repos (set `GITHUB_TOKEN` to include private)
ph org github.com/nousefreak

# Open a fuzzyfinder that will `cd` to the repo 
ph go [search]

# Shorthand for `ph go`
ph [search]

# Run a `git fetch` on all repos
ph update                                              

# Show commits made to any repository in the last 2 days
ph wdid 2 days
```

## Badges

[![Releases](https://img.shields.io/github/v/release/nousefreak/projecthelper?style=for-the-badge)](https://github.com/NoUseFreak/projecthelper/releases)
[![Build status](https://img.shields.io/github/actions/workflow/status/nousefreak/projecthelper/ci.yml?style=for-the-badge)](ihttps://github.com/NoUseFreak/projecthelper/actions/workflows/ci.yml)
[![GitHub License](https://img.shields.io/github/license/nousefreak/projecthelper?style=for-the-badge)](https://github.com/NoUseFreak/projecthelper/blob/main/LICENSE)
[![Static Badge](https://img.shields.io/badge/powered--by-stenic.io-blue?style=for-the-badge&logoColor=blue)](https://stenic.io)


