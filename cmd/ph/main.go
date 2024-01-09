package main

import "github.com/nousefreak/projecthelper/internal/pkg/command"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	command.Version = Version
	command.Commit = Commit
	command.Date = Date
	command.Execute()
}
