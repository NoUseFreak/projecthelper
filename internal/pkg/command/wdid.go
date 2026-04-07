package command

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/nousefreak/projecthelper/internal/pkg/color"
	"github.com/nousefreak/projecthelper/internal/pkg/config"
	"github.com/nousefreak/projecthelper/internal/pkg/repo"
	"github.com/nousefreak/projecthelper/internal/pkg/wdid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func getWDIDCmd() *cobra.Command {
	var (
		filter      string
		groupByDate bool
	)

	cmd := &cobra.Command{
		Use:   "wdid [AMOUNT UNIT]",
		Short: "Generate a WhatDidIDo report over a period of time",
		Run: func(cmd *cobra.Command, args []string) {

			amount := "1"
			unit := "day"

			if len(args) >= 1 {
				amount = args[0]
			}
			if len(args) >= 2 {
				unit = args[1]
			}

			window := fmt.Sprintf("%s %s ago", amount, unit)

			repoPaths, err := repo.GetRepoPaths(config.GetBaseDir())
			if err != nil {
				logrus.Fatal(err)
			}

			if filter != "" {
				repoPaths = repo.FilterRepoPaths(repoPaths, filter)
			}

			logrus.Infof("Looking for your commits in %d repos since %s", len(repoPaths), window)
			reports, err := wdid.GetWDIDReport(window, repoPaths, groupByDate)
			if err != nil {
				logrus.Fatal(err)
			}

			width, _, _ := term.GetSize(0)
			out := os.Stderr

			// Sort the keys (dates) if grouping by date
			var sortedKeys []string
			if groupByDate {
				for name := range reports {
					sortedKeys = append(sortedKeys, name)
				}
				sort.Strings(sortedKeys) // This will sort dates chronologically (oldest first)
			} else {
				// For repo grouping, maintain original order
				for name := range reports {
					sortedKeys = append(sortedKeys, name)
				}
			}

			for _, name := range sortedKeys {
				group := reports[name]
				fmt.Fprintln(out, "")
				fmt.Fprintf(out, "%s\n", color.Color(color.FgGreen, name))
				for _, r := range group {
					t := time.Unix(r.Timestamp(), 0)
					if groupByDate {
						// Check if this is a repository separator
						if r.IsRepoSeparator() {
							repoName := r.GetSeparatorRepo()
							fmt.Fprintf(out, "  %s\n", color.Color(color.FgYellow, repoName))
							continue
						}
						// Regular commit entry with timestamp
						fmt.Fprintf(out, "   - %-*s %s\n", width-25, r.ChangeLine(), color.Color(color.FgBlue, t.Format("2006-01-02 15:04")))
					} else {
						fmt.Fprintf(out, " - %-*s %s\n", width-22, r.ChangeLine(), color.Color(color.FgBlue, t.Format("2006-01-02 15:04")))
					}
				}
			}
		},
	}

	cmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter the projects by a substring")
	cmd.Flags().BoolVarP(&groupByDate, "group-by-date", "d", false, "Group commits by date instead of repository")

	return cmd

}
