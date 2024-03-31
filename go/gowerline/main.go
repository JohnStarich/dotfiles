package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// Reference: https://tao-of-tmux.readthedocs.io/en/latest/manuscript/09-status-bar.html

/*
#[fg=#121212,bg=default,nobold,noitalics,nounderscore] #[fg=#797aac,bg=#121212,nobold,noitalics,nounderscore] 🌪  57.0°F#[fg=#f3e6d8,bg=#121212,nobold,noitalics,nounderscore] #[fg=#f3e6d8,bg=#121212,nobold,noitalics,nounderscore] 🔥 74%#[fg=#303030,bg=#121212,nobold,noitalics,nounderscore] #[fg=#9e9e9e,bg=#303030,nobold,noitalics,nounderscore] Mon Mar 25#[fg=#626262,bg=#303030,nobold,noitalics,nounderscore] #[fg=#d0d0d0,bg=#303030,bold,noitalics,nounderscore] 05:12 PM
*/

type FontConfig struct {
	Foreground string
	Background string
	Bold       bool
	Italics    bool
	Underscore bool
}

func (f FontConfig) String() string {
	return fmt.Sprintf(`#[fg=%s,bg=%s,%sbold,%sitalics,%sunderscore]`, f.Foreground, f.Background, boolToYesNo(f.Bold), boolToYesNo(f.Italics), boolToYesNo(f.Underscore))
}

func boolToYesNo(b bool) string {
	if b {
		return ""
	}
	return "no"
}

const (
	powerlineArrowPointLeftFull  = ""
	powerlineArrowPointLeftEmpty = ""
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := run(ctx)
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	return status(os.Stdout)
}
