package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
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

	args := Args{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	flag.BoolVar(&args.IgnoreHidden, "hidden", true, "Ignore hidden files (starts with '.').")
	flag.IntVar(&args.MaxDepth, "max-depth", 20, "Maximum watch depth. Uses a default for better performance.")
	flag.Var(&args.IgnorePatterns, "ignore", "File path pattern to ignore.")
	flag.Parse()
	args.RootDir = flag.Arg(0)
	args.FilePatterns = strings.Split(flag.Arg(1), ",")
	args.CommandArgs = skipStrings(flag.Args(), 2)

	err := args.Run(ctx)
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
	CommandArgs    []string
	FilePatterns   []string
	IgnoreHidden   bool
	IgnorePatterns stringSliceVar
	MaxDepth       int
	RootDir        string
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func (a Args) Run(ctx context.Context) error {
	if a.RootDir == "" {
		return errors.New("provide a root directory to watch")
	}
	var err error
	a.RootDir, err = filepath.Abs(a.RootDir)
	if err != nil {
		return err
	}
	if len(a.FilePatterns) == 0 {
		return errors.New("provide at least one file extension to filter by, or more than one separated by commas ','")
	}
	if len(a.CommandArgs) == 0 {
		return errors.New("must provide at least one command arg")
	}
	var env []string
	if contents, err := os.ReadFile(filepath.Join(a.RootDir, ".env")); err == nil {
		env, err = parseEnvFile(contents)
		if err != nil {
			return errors.Wrap(err, "failed to parse .env file in root directory")
		}
	}

	ranOnce := false
	return a.runWatch(ctx, func(filePaths []string) error {
		if ranOnce && !matchesFilePatterns(a.FilePatterns, filePaths) {
			return nil
		}
		shouldIgnore, err := a.ShouldIgnoreAll(filePaths)
		if err == nil && shouldIgnore {
			return nil
		}

		ranOnce = true
		clearScreen()
		fmt.Fprintln(a.Stdout, "### Running command:", strings.Join(a.CommandArgs, " "))

		cmd := exec.CommandContext(ctx, a.CommandArgs[0], a.CommandArgs[1:]...)
		cmd.Stdin = a.Stdin
		cmd.Stdout = a.Stdout
		cmd.Stderr = a.Stderr
		if len(env) > 0 {
			cmd.Env = append(os.Environ(), env...)
			fmt.Fprintln(a.Stdout, "### Including environment from .env")
		}
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(a.Stderr, err)
		}
		fmt.Fprintf(a.Stdout, "### Waiting for changes... (patterns: %s)\n", strings.Join(a.FilePatterns, ", "))
		return nil
	})
}

type callbackFunc func(filePaths []string) error

// runWatch runs a recursive file watcher starting at 'rootDir', which calls 'callback' on every file change.
// Only returns when 'ctx' is canceled.
func (a Args) runWatch(ctx context.Context, callback callbackFunc) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = filepath.WalkDir(a.RootDir, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dirEntry.IsDir() {
			if shouldIgnore, err := a.ShouldIgnore(path); err != nil || shouldIgnore {
				return filepath.SkipDir
			}
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
	var lastEvents []fsnotify.Event
	for {
		select {
		case <-timer.C:
			var filePaths []string
			for _, event := range lastEvents {
				filePaths = append(filePaths, event.Name)
			}
			lastEvents = nil
			err := callback(filePaths)
			if err != nil {
				fmt.Fprintln(a.Stderr, "Error running watch call:", err)
			}
		case <-ctx.Done():
			return nil
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Remove|fsnotify.Create) != 0 {
				info, err := os.Stat(event.Name)
				isDir := err == nil && info.IsDir()
				switch {
				case isDir && event.Op&fsnotify.Remove != 0:
					_ = watcher.Remove(event.Name)
				case isDir && event.Op&fsnotify.Create != 0:
					_ = watcher.Add(event.Name)
				case !isDir && event.Op&fsnotify.Create != 0:
					_ = watcher.Add(filepath.Dir(event.Name))
				}
			}
			if event.Op&(fsnotify.Write|fsnotify.Remove|fsnotify.Create|fsnotify.Rename) != 0 {
				lastEvents = append(lastEvents, event)
				timer.Reset(debounce)
			}
		case err := <-watcher.Errors:
			var pathErr *os.PathError
			if errors.As(err, &pathErr) && errors.Is(err, os.ErrNotExist) {
				_ = watcher.Remove(pathErr.Path)
			} else {
				fmt.Fprintln(a.Stderr, "Watch error:", err)
			}
		}
	}
}

func matchesFilePatterns(filePatterns []string, filePaths []string) bool {
	for _, filePath := range filePaths {
		for _, pattern := range filePatterns {
			if strings.HasSuffix(filePath, "."+pattern) {
				return true
			}
		}
	}
	return false
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (a Args) ShouldIgnoreAll(paths []string) (shouldIgnore bool, err error) {
	for _, path := range paths {
		shouldIgnore, err = a.ShouldIgnore(path)
		if err != nil || !shouldIgnore {
			return
		}
	}
	return
}

func (a Args) ShouldIgnore(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	if a.IgnoreHidden && strings.HasPrefix(filepath.Base(path), ".") {
		return true, nil
	}
	relPath, err := filepath.Rel(a.RootDir, path)
	if err != nil {
		return false, err
	}
	depth := strings.Count(relPath, string(filepath.Separator)) + 1
	if depth > a.MaxDepth {
		return true, nil
	}
	for _, ignorePattern := range a.IgnorePatterns {
		matched, err := filepath.Match(ignorePattern, relPath)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}
