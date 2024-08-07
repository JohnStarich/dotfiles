package main

import (
	"context"
	"fmt"
	goos "os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/hack-pad/hackpadfs/os"
	"github.com/johnstarich/go/gowerline/internal/power"
	"github.com/johnstarich/go/gowerline/internal/status"
	"github.com/johnstarich/go/gowerline/internal/weather"
	"github.com/pkg/errors"
)

const appName = "gowerline"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), goos.Interrupt)
	defer cancel()
	err := run(ctx, goos.Args[1:])
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("an action is required: gowerline (status-right|tmux-setup)")
	}
	action := args[0]
	var debug bool
	if len(args) > 1 {
		debug = args[1] == "--debug"
	}
	switch action {
	case "status-right":
		return generateStatus(ctx, debug)
	case "tmux-setup":
		return setUpTmux(ctx, debug)
	default:
		return errors.Errorf("unrecognized action: %q; gowerline (status-right|tmux-setup)", args[0])
	}
}

func generateStatus(ctx context.Context, debug bool) error {
	fs := os.NewFS()
	cacheDir, err := goos.UserCacheDir()
	if err != nil {
		return err
	}
	absoluteCacheDir, err := filepath.Abs(cacheDir)
	if err != nil {
		return err
	}
	cacheSubPath, err := fs.FromOSPath(absoluteCacheDir)
	if err != nil {
		return err
	}
	appCacheSubPath := path.Join(cacheSubPath, appName)
	cacheFS, err := fs.Sub(appCacheSubPath)
	if err != nil {
		return err
	}

	statusLine := status.Line{
		SkipCache: debug,
		Segments: []status.Segment{
			{
				Font:            status.Font{Foreground: "#797aac", Background: defaultSecondaryColor},
				GenerateContent: weather.Status,
				Name:            "weather",
				SeparatorFont:   status.Font{Foreground: defaultSecondaryColor},
			},
			{
				Font:            status.Font{Foreground: "#f3e6d8", Background: defaultSecondaryColor},
				GenerateContent: power.Status,
				Name:            "battery",
				SeparatorFont:   status.Font{Foreground: "#f3e6d8", Background: defaultSecondaryColor},
			},
			{
				Font: status.Font{Foreground: "#9e9e9e", Background: "#303030"},
				GenerateContent: func(ctx status.Context) (time.Duration, error) {
					const dateFormat = "Mon Jan _2"
					fmt.Fprint(ctx.Writer, time.Now().Format(dateFormat))
					return 0, nil
				},
				Name:          "date",
				SeparatorFont: status.Font{Foreground: "#303030", Background: defaultSecondaryColor},
				Separator:     status.Separator{FullArrow: true},
			},
			{
				Font: status.Font{Foreground: defaultPrimaryColor, Background: "#303030", Bold: true},
				GenerateContent: func(ctx status.Context) (time.Duration, error) {
					const timeFormat = "3:04 PM"
					fmt.Fprint(ctx.Writer, time.Now().Format(timeFormat))
					return 0, nil
				},
				Name:          "time",
				SeparatorFont: status.Font{Foreground: "#626262", Background: "#303030"},
			},
		},
	}
	return statusLine.Status(ctx, goos.Stdout, goos.Stderr, cacheFS)
}
