package command

import (
	"fmt"
	"os"
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
	cmd := &cobra.Command{
		Use:   "wdid",
		Short: "wdid",
		Long:  `wdid`,
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

			logrus.Infof("Looking for your commits in %d repos since %s", len(repoPaths), window)
			reports, err := wdid.GetWDIDReport(window, repoPaths)
			if err != nil {
				logrus.Fatal(err)
			}

			width, _, _ := term.GetSize(0)
			out := os.Stderr
			for name, group := range reports {
				fmt.Fprintln(out, "")
				fmt.Fprintf(out, "%s\n", color.Color(color.FgGreen, name))
				for _, r := range group {
					t := time.Unix(r.Timestamp(), 0)
					fmt.Fprintf(out, " - %-*s %s\n", width-22, r.ChangeLine(), color.Color(color.FgBlue, t.Format("2006-01-02 15:04")))
				}
			}

		},
	}
	return cmd

}
