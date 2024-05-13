package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/johnstarich/go/gowerline/internal/status"
)

//go:embed tmux.conf.tmpl
var tmuxConfTemplate string

type tmuxData struct {
	Options map[string]string
}

func setUpTmux(ctx context.Context, debug bool) error {
	args := []string{"source-file"}
	if debug {
		args = append(args, "-v")
	}
	args = append(args, "/dev/stdin")
	cmd := exec.CommandContext(ctx, "tmux", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	w, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = writeTMUXConfig(w)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

const (
	defaultPrimaryColor      = "#ffffff"
	defaultSecondaryColor    = "#111111"
	inactivePrimaryColor     = "#bbbbbb"
	inactiveLeftPrimaryColor = "#eeeeee"
	activeColor              = "#6699cc"
)

func writeTMUXConfig(w io.Writer) error {
	inactiveLeftFont := status.Font{
		Foreground: inactiveLeftPrimaryColor,
		Background: defaultSecondaryColor,
		Bold:       true,
	}
	inactiveLeftSeparatorFont := status.Font{
		Foreground: inactiveLeftPrimaryColor,
		Background: defaultSecondaryColor,
		Bold:       true,
	}
	inactiveWindowFont := status.Font{
		Foreground: inactivePrimaryColor,
		Background: defaultSecondaryColor,
	}
	activeWindowFont := status.Font{
		Foreground: defaultPrimaryColor,
		Background: activeColor,
		Bold:       true,
	}
	activeWindowSeparatorFont := status.Font{
		Foreground: defaultSecondaryColor,
		Background: activeColor,
		Bold:       true,
	}

	statusLeft := join(
		`#{?client_prefix,`, activeWindowFont.VariableSafeString(), `,`, inactiveLeftFont.InvertForeground().VariableSafeString(), `}`,
		` #{session_name} `,
		`#{?client_prefix,`, activeWindowSeparatorFont.InvertForeground().VariableSafeString(), `,`, inactiveLeftSeparatorFont.VariableSafeString(), `}`,
		status.Separator{FullArrow: true, PointRight: true}.String(),
	)

	windowFormat := join(
		`#{window_index}`,
		`#{?window_flags,#{window_flags}, } `,
		status.Separator{PointRight: true}.String(),
		` #{window_name} `,
	)
	windowStatus := fmt.Sprintf("   %s ", windowFormat)
	currentWindowStatus := join(
		` `,
		activeWindowSeparatorFont.String(),
		status.Separator{FullArrow: true, PointRight: true}.String(),
		` `,
		activeWindowFont.String(),
		windowFormat,
		activeWindowSeparatorFont.InvertForeground().String(),
		status.Separator{FullArrow: true, PointRight: true}.String(),
	)

	statusRight := `#(PATH="$HOME/go/bin:$PATH" "$HOME/.dotfiles/bin/gowerline" status-right)`
	return template.Must(template.New("").Parse(tmuxConfTemplate)).Execute(w, tmuxData{
		Options: map[string]string{
			"status":                       "on",                       // Enable status line.
			"status-interval":              "2",                        // Set update interval between generating status lines.
			"status-left":                  statusLeft,                 // Generate left status.
			"status-left-length":           "200",                      // Set maximum width of left status.
			"status-right":                 statusRight,                // Generate right status.
			"status-right-length":          "200",                      // Set maximum width of right status.
			"status-style":                 inactiveWindowFont.Style(), // Set default style like foreground and background color.
			"window-status-current-format": currentWindowStatus,        // Generate status for windows on the left side.
			"window-status-format":         windowStatus,               // Generate status for windows on the left side.
		},
	})
}

func join(s ...string) string {
	return strings.Join(s, "")
}
