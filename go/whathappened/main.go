package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/markusmobius/go-dateparser"
	"github.com/pkg/errors"
)

/*

Useful queries:

sqlite3 ~/Library/Calendars/Calendar.sqlitedb "select summary, start_date from CalendarItem where start_date>($(gdate +%s -d 'yesterday 8am')-978307200) AND end_date<($(gdate +%s -d today)-978307200) limit 100;"
sqlite3 ~/Library/Calendars/Calendar.sqlitedb '.schema'
sqlite3 ~/Library/Calendars/Calendar.sqlitedb "select calendar.title, item.summary, item.start_date from CalendarItem as item left join Calendar as calendar on (calendar.rowid = item.calendar_id) where item.start_date>($(gdate +%s -d 'yesterday 8am')-978307200) AND item.end_date<($(gdate +%s -d today)-978307200) limit 100;"

*/

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := run(ctx, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	set := flag.NewFlagSet("whathappened", flag.ContinueOnError)
	since := set.String("since", "yesterday", "The relative or absolute start time to look up what happened")
	emailAddresses := newStringSliceFlag(availableGitEmailAddresses(ctx))
	set.Var(emailAddresses, "author-email", "The author email to pull git commits for. Can be repeated for multiple email addresses. Always adds global and local git config email.")
	if err := set.Parse(args); err != nil {
		return err
	}

	year, month, day := time.Now().Date()
	now := time.Date(year, month, day, 6, 0, 0, 0, time.Local)
	sinceTime, err := dateparser.Parse(&dateparser.Configuration{
		CurrentTime: now,
	}, *since)
	if err != nil {
		return err
	}

	commits, err := recentCommits(ctx, sinceTime.Time, now, emailAddresses.Values())
	if err != nil {
		return err
	}
	for _, commit := range commits {
		fmt.Println("-", commit)
	}
	/* Non-functional on macOS Sonoma.
	events, err := recentEvents(ctx, yesterday, now, map[string][]string{
		workEmail: {"Calendar"},
	})
	if err != nil {
		return err
	}
	for _, event := range events {
		fmt.Println("-", event)
	}
	*/
	return nil
}

func recentEvents(ctx context.Context, start, end time.Time, calendars map[string][]string) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "sqlite3", filepath.Join(home, "Library/Calendars/Calendar.sqlitedb"), fmt.Sprintf("select calendar.owner_identity_email, calendar.title, item.summary, item.start_date from CalendarItem as item left join Calendar as calendar on (calendar.rowid = item.calendar_id) where item.start_date>%d AND item.end_date<%d limit 100;", timeToMacOSCocoaCoreData(start), timeToMacOSCocoaCoreData(end)), "--csv")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.WithMessage(err, string(output))
	}
	reader := csv.NewReader(bytes.NewReader(output))
	columns, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var messages []string
	for _, columns := range columns {
		ownerEmail, calendarTitle, itemSummary := columns[0], columns[1], columns[2]
		if slices.Contains(calendars[ownerEmail], calendarTitle) {
			messages = append(messages, "meeting: "+itemSummary)
		}
	}
	return messages, nil
}

func timeToMacOSCocoaCoreData(t time.Time) int64 {
	const appleCocoaCoreDataTimestampOffsetSecondsToUnixEpoch = 978307200 // https://www.epochconverter.com/coredata
	return t.Unix() - appleCocoaCoreDataTimestampOffsetSecondsToUnixEpoch
}

func recentCommits(ctx context.Context, start, end time.Time, authors []string) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	matches, err := filepath.Glob(filepath.Join(home, "projects", "*", ".git"))
	if err != nil {
		return nil, err
	}
	var messages []string
	for _, match := range matches {
		projectDir := filepath.Dir(match)
		projectName := filepath.Base(projectDir)
		repo, err := git.PlainOpen(projectDir)
		if err != nil {
			return nil, err
		}
		commitRef, err := repo.Head()
		if err != nil {
			return nil, err
		}
		commit, err := repo.CommitObject(commitRef.Hash())
		if err != nil {
			return nil, err
		}
		for !commit.Author.When.Before(start) && !commit.Author.When.After(end) {
			if slices.Contains(authors, commit.Author.Email) {
				messages = append(messages, fmt.Sprintf("%s: %s", projectName, firstLine(commit.Message)))
			}
			commit, err = commit.Parent(1)
			if err != nil {
				break
			}
		}
	}
	return messages, nil
}

func firstLine(s string) string {
	index := strings.IndexAny(s, "\r\n")
	if index != -1 {
		return s[:index]
	}
	return s
}

func availableGitEmailAddresses(ctx context.Context) []string {
	emailAddresses := make(map[string]struct{})

	localEmail, err := exec.CommandContext(ctx, "git", "config", "user.email").CombinedOutput()
	if err == nil {
		localEmail = bytes.TrimSpace(localEmail)
		emailAddresses[string(localEmail)] = struct{}{}
	}

	globalEmail, err := exec.CommandContext(ctx, "git", "config", "--global", "user.email").CombinedOutput()
	if err == nil {
		globalEmail = bytes.TrimSpace(globalEmail)
		emailAddresses[string(globalEmail)] = struct{}{}
	}
	return slices.Collect(maps.Keys(emailAddresses))
}
