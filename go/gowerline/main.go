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
	"github.com/johnstarich/go/gowerline/internal/status"
)

const appName = "gowerline"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), goos.Interrupt)
	defer cancel()
	err := run(ctx)
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
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

	return status.Status(ctx, goos.Stdout, cacheFS, status.Line{
		Segments: []status.Segment{
			{ // weather
				Separator: status.Separator{Font: status.Font{Foreground: "#121212"}},
				Font:      status.Font{Foreground: "#797aac", Background: "#121212"},
				Content:   weatherStatus,
			},
			{ // battery
				Separator: status.Separator{Font: status.Font{Foreground: "#f3e6d8", Background: "#121212"}},
				Font:      status.Font{Foreground: "#f3e6d8", Background: "#121212"},
				Content:   batteryStatus,
			},
			{ // date
				Separator: status.Separator{Font: status.Font{Foreground: "#303030", Background: "#121212"}, FullArrow: true},
				Font:      status.Font{Foreground: "#9e9e9e", Background: "#303030"},
				Content: func(ctx status.Context) error {
					fmt.Fprint(ctx.Writer, time.Now().Format(time.DateOnly))
					return nil
				},
			},
			{ // time
				Separator: status.Separator{Font: status.Font{Foreground: "#626262", Background: "#303030"}},
				Font:      status.Font{Foreground: "#d0d0d0", Background: "#303030", Bold: true},
				Content: func(ctx status.Context) error {
					const timeFormat = "3:04 PM"
					fmt.Fprint(ctx.Writer, time.Now().Format(timeFormat))
					return nil
				},
			},
		},
	})
}
