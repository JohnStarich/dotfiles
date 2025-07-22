package main

import (
	"bytes"
	"context"
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
)

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
	return nil
}

func recentCommits(_ context.Context, start, end time.Time, authors []string) ([]string, error) {
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
