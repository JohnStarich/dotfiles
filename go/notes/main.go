package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/markusmobius/go-dateparser"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

const (
	appName        = "notes"
	usage          = "Usage: notes edit <subject> [date ...]"
	noteDateFormat = "2006-01-02"
	noteExtension  = ".md"
)

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", appName, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return errors.New(usage)
	}
	switch args[0] {
	case "edit", "e":
		subject := args[1]
		day := time.Now()
		if len(args) > 2 {
			date, err := dateparser.Parse(nil, strings.Join(args[2:], " "))
			if err != nil {
				return errors.WithMessage(err, "invalid date")
			}
			day = date.Time
		}
		return edit(subject, day)
	case "search", "s":
		return search(strings.Join(args[1:], " "))
	default:
		return errors.New(usage)
	}
}

func subjectBasePath(subject string) (string, error) {
	notesBase, err := notesBase()
	if err != nil {
		return "", err
	}
	const glob = "*"
	subjectBasePattern := filepath.Join(notesBase, subject+glob)
	possibleSubjectPaths, err := filepath.Glob(subjectBasePattern)
	if err != nil {
		return "", err
	}
	if len(possibleSubjectPaths) == 0 {
		return "", errors.Errorf("subject not found matching pattern: %s", subjectBasePattern)
	}
	return possibleSubjectPaths[0], nil
}

func edit(subject string, day time.Time) error {
	subjectBasePath, err := subjectBasePath(subject)
	if err != nil {
		return err
	}
	notePath := filepath.Join(subjectBasePath, day.Format(noteDateFormat)+noteExtension)
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println(notePath)
		return nil
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func search(search string) error {
	notesBase, err := notesBase()
	if err != nil {
		return err
	}
	name, args := searchCommand(search, notesBase)
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func searchCommand(search, path string) (string, []string) {
	silverSearcherPath, err := exec.LookPath("ag")
	if err == nil {
		return silverSearcherPath, []string{search, path}
	}
	return "grep", []string{"-RE", search, path}
}

func notesBase() (string, error) {
	notesBase := os.Getenv("NOTES_BASE")
	if notesBase == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		notesBase = filepath.Join(home, "notes")
	}
	return notesBase, nil
}
