package wdid

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nousefreak/projecthelper/internal/pkg/config"
	"github.com/sirupsen/logrus"

	parallel "github.com/NoUseFreak/go-parallel"
)

func GetWDIDReport(window string, repoPaths []string, groupByDate bool) (map[string][]RepoLog, error) {
	input := parallel.Input{}
	for _, repo := range repoPaths {
		input = append(input, repo)
	}

	author, err := getGitAuthor()
	if err != nil {
		logrus.Fatal(err)
	}

	p := parallel.Processor{Threads: 10}
	result := p.Process(input, func(i interface{}) interface{} {
		project := i.(string)

		logrus.Debug("Checking", project, "for commits by", author)
		cmd := exec.Command("git", "log", "--reverse", "--no-merges", "--pretty=%ct %s", "--since='"+window+"'", "--author="+author, "--branches=*")
		cmd.Dir = project
		out, err := cmd.Output()
		if err != nil {
			logrus.Error(fmt.Errorf("error running git log: %w", err))
			return nil
		}
		s := strings.TrimSpace(string(out))
		if s == "" {
			return nil
		}

		var res []RepoLog
		for _, l := range strings.Split(s, "\n") {
			res = append(res, RepoLog{
				line:    l,
				project: project,
			})
		}
		return res
	})

	var full []RepoLog
	for _, r := range result {
		if l, ok := r.([]RepoLog); ok {
			full = append(full, l...)
		}
	}

	sort.SliceStable(full, func(i, j int) bool {
		return full[i].Timestamp() < full[j].Timestamp()
	})

	if groupByDate {
		return intoDateGroups(unique(full)), nil
	}
	return intoGroups(unique(full)), nil
}

func getGitAuthor() (string, error) {
	out, err := exec.Command("git", "config", "--get", "user.email").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func unique(slice []RepoLog) []RepoLog {
	keys := make(map[string]bool)
	var list []RepoLog
	for _, entry := range slice {
		k := fmt.Sprintf("%s-%s", entry.ChangeLine(), entry.project)
		if _, value := keys[k]; !value {
			keys[k] = true
			list = append(list, entry)
		}
	}

	return list
}

type RepoLog struct {
	line    string
	project string
}

func (r *RepoLog) String() string {
	return fmt.Sprintf("%s (%s)", r.ChangeLine(), r.ShortRepo())
}

func (r *RepoLog) ShortRepo() string {
	return strings.TrimPrefix(strings.TrimPrefix(r.project, config.GetBaseDir()), "/")
}
func (r *RepoLog) ChangeLine() string {
	return r.line[11:]
}
func (r *RepoLog) Timestamp() int64 {
	if t, err := strconv.ParseInt(r.line[0:10], 10, 64); err == nil {
		return t
	}
	return 0
}

func (r *RepoLog) IsRepoSeparator() bool {
	return strings.Contains(r.line, "REPO_SEPARATOR:")
}

func (r *RepoLog) GetSeparatorRepo() string {
	if r.IsRepoSeparator() {
		parts := strings.Split(r.line, "REPO_SEPARATOR:")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return ""
}

func intoGroups(slice []RepoLog) map[string][]RepoLog {
	groups := make(map[string][]RepoLog)
	for _, entry := range slice {
		groups[entry.ShortRepo()] = append(groups[entry.ShortRepo()], entry)
	}
	return groups
}

func intoDateGroups(slice []RepoLog) map[string][]RepoLog {
	dateGroups := make(map[string]map[string][]RepoLog)

	// First group by date, then by repo
	for _, entry := range slice {
		t := time.Unix(entry.Timestamp(), 0)
		dateKey := t.Format("2006-01-02")
		repoKey := entry.ShortRepo()

		if dateGroups[dateKey] == nil {
			dateGroups[dateKey] = make(map[string][]RepoLog)
		}
		dateGroups[dateKey][repoKey] = append(dateGroups[dateKey][repoKey], entry)
	}

	// Flatten into the expected format with repo subheadings
	result := make(map[string][]RepoLog)
	for dateKey, repoGroups := range dateGroups {
		var allEntries []RepoLog
		for repoKey, entries := range repoGroups {
			// Add a separator entry for the repo name
			separator := RepoLog{
				line:    fmt.Sprintf("0000000000 REPO_SEPARATOR:%s", repoKey),
				project: "",
			}
			allEntries = append(allEntries, separator)
			allEntries = append(allEntries, entries...)
		}
		result[dateKey] = allEntries
	}

	return result
}
