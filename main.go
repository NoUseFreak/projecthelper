package main

import (
	"github.com/nousefreak/projecthelper/internal/pkg/command"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	command.Version = version
	command.Commit = commit
	command.Date = date

	command.Execute()
}
