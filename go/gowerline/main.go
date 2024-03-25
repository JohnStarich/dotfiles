package main

import (
	"fmt"
	"time"
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
	fmt.Print(FontConfig{
		Foreground: "#121212",
		Background: "default",
	})
	fmt.Print(" " + powerlineArrowPointLeftFull)
	fmt.Print(FontConfig{
		Foreground: "#797aac",
		Background: "#121212",
	})
	fmt.Print(" 🌪  57.0°F")
	fmt.Print(FontConfig{
		Foreground: "#f3e6d8",
		Background: "#121212",
	})
	fmt.Print(" " + powerlineArrowPointLeftEmpty)
	fmt.Print(FontConfig{
		Foreground: "#f3e6d8",
		Background: "#121212",
	})
	fmt.Print(" 🔥 74%")
	fmt.Print(FontConfig{
		Foreground: "#303030",
		Background: "#121212",
	})
	fmt.Print(" " + powerlineArrowPointLeftFull)
	fmt.Print(FontConfig{
		Foreground: "#9e9e9e",
		Background: "#303030",
	})
	fmt.Print(time.Now().Format(time.DateOnly))
	fmt.Print(FontConfig{
		Foreground: "#626262",
		Background: "#303030",
	})
	fmt.Print(" " + powerlineArrowPointLeftEmpty)
	fmt.Print(FontConfig{
		Foreground: "#d0d0d0",
		Background: "#303030",
		Bold:       true,
	})
	fmt.Print(time.Now().Format(time.TimeOnly))
	fmt.Println()
}
