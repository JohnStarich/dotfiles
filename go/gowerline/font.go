package main

import "fmt"

// Reference: https://tao-of-tmux.readthedocs.io/en/latest/manuscript/09-status-bar.html

type Font struct {
	Foreground string
	Background string
	Bold       bool
	Italics    bool
	Underscore bool
}

func (f Font) String() string {
	if f.Foreground == "" {
		f.Foreground = "default"
	}
	if f.Background == "" {
		f.Background = "default"
	}
	return fmt.Sprintf(`#[fg=%s,bg=%s,%sbold,%sitalics,%sunderscore]`, f.Foreground, f.Background, boolToYesNo(f.Bold), boolToYesNo(f.Italics), boolToYesNo(f.Underscore))
}

func boolToYesNo(b bool) string {
	if b {
		return ""
	}
	return "no"
}
