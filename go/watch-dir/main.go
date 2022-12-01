package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	increaseNoFile()

	flag.Parse()
	args := Args{
		RootDir:      flag.Arg(0),
		FilePatterns: strings.Split(flag.Arg(1), ","),
		CommandArgs:  skipStrings(flag.Args(), 2),
		Stdin:        os.Stdin,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
	}

	err := run(ctx, args)
	if err != nil {
		panic(err)
	}
}

func skipStrings(str []string, skip int) []string {
	if len(str) < skip {
		return nil
	}
	return str[skip:]
}

type Args struct {
	RootDir        string
	FilePatterns   []string
	CommandArgs    []string
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func run(ctx context.Context, args Args) error {
	if args.RootDir == "" {
		return errors.New("provide a root directory to watch")
	}
	if len(args.FilePatterns) == 0 {
		return errors.New("provide at least one file extension to filter by, or more than one separated by commas ','")
	}
	if len(args.CommandArgs) == 0 {
		return errors.New("must provide at least one command arg")
	}

	ranOnce := false
	return runWatch(ctx, args.RootDir, args.Stderr, func(filePath string) error {
		if ranOnce && !matchesFilePatterns(args.FilePatterns, filePath) {
			return nil
		}
		ranOnce = true
		clearScreen()
		fmt.Fprintln(args.Stdout, "### Running command:", strings.Join(args.CommandArgs, " "))

		cmd := exec.CommandContext(ctx, args.CommandArgs[0], args.CommandArgs[1:]...)
		cmd.Stdin = args.Stdin
		cmd.Stdout = args.Stdout
		cmd.Stderr = args.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(args.Stderr, err)
		}
		fmt.Fprintf(args.Stdout, "### Waiting for changes... (patterns: %s)\n", strings.Join(args.FilePatterns, ", "))
		return nil
	})
}

type callbackFunc func(filePath string) error

// runWatch runs a recursive file watcher starting at 'rootDir', which calls 'callback' on every file change.
// Only returns when 'ctx' is canceled.
func runWatch(ctx context.Context, rootDir string, errWriter io.Writer, callback callbackFunc) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	timer := time.NewTimer(0) // fire watch right away
	defer timer.Stop()

	const debounce = 2 * time.Second
	var lastEvent fsnotify.Event
	for {
		select {
		case <-timer.C:
			err := callback(lastEvent.Name)
			if err != nil {
				fmt.Fprintln(errWriter, "Error running watch call:", err)
			}
		case <-ctx.Done():
			return nil
		case event := <-watcher.Events:
			switch {
			case event.Op&fsnotify.Write == fsnotify.Write,
				event.Op&fsnotify.Create == fsnotify.Create:
				lastEvent = event
				timer.Reset(debounce)
			}
		case err := <-watcher.Errors:
			fmt.Fprintln(errWriter, "Watch error:", err)
		}
	}
}

func matchesFilePatterns(filePatterns []string, filePath string) bool {
	for _, pattern := range filePatterns {
		if strings.HasSuffix(filePath, "."+pattern) {
			return true
		}
	}
	return false
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
