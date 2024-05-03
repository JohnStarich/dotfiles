package status

import (
	"fmt"
	"strings"
)

// Reference: https://tao-of-tmux.readthedocs.io/en/latest/manuscript/09-status-bar.html

type Font struct {
	Foreground string
	Background string
	Bold       bool
	Italics    bool
	Underscore bool
}

func (f Font) InvertForeground() Font {
	originalForeground := f.Foreground
	f.Foreground = f.Background
	f.Background = originalForeground
	return f
}

func (f Font) style() []string {
	if f.Foreground == "" {
		f.Foreground = "default"
	}
	if f.Background == "" {
		f.Background = "default"
	}
	return []string{
		`fg=` + f.Foreground,
		`bg=` + f.Background,
		boolToYesNo(f.Bold) + `bold`,
		boolToYesNo(f.Italics) + `italics`,
		boolToYesNo(f.Underscore) + `underscore`,
	}
}

func (f Font) Style() string {
	return strings.Join(f.style(), ",")
}

func (f Font) VariableSafeString() string {
	return "#[" + strings.Join(f.style(), "]#[") + "]"
}

func (f Font) String() string {
	return fmt.Sprintf("#[%s]", f.Style())
}

func boolToYesNo(b bool) string {
	if b {
		return ""
	}
	return "no"
}
