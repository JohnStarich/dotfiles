package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	err = writeTMUXConfig(w, debug)
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
	bellColor                = "#bb0000"
)

func writeTMUXConfig(w io.Writer, debug bool) error {
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

	statusLeft := status.Join(
		status.Ternary{
			Variable: "client_prefix",
			True:     status.String(activeWindowFont.String()),
			False:    status.String(inactiveLeftFont.InvertForeground().String()),
		},
		status.String(" "),
		status.Variable("session_name"),
		status.String(" "),
		status.Ternary{
			Variable: "client_prefix",
			True:     status.String(activeWindowSeparatorFont.InvertForeground().String()),
			False:    status.String(inactiveLeftSeparatorFont.String()),
		},
		status.Separator{FullArrow: true, PointRight: true},
	)

	windowFormat := status.Join(
		status.String(" "),
		status.Variable("window_index"),
		status.Ternary{
			Variable: "window_flags",
			True:     status.Variable("window_flags"),
			False:    status.String(" "),
		},
		status.String(" "),
		status.Separator{PointRight: true},
		status.String(" "),
		status.Variable("window_name"),
		status.String(" "),
	)
	windowStatus := fmt.Sprintf(" %s ", windowFormat)

	bellWindowStyle := status.Font{
		Foreground: bellColor,
		Bold:       true,
	}.Style()

	currentWindowStatus := status.Join(
		activeWindowSeparatorFont,
		status.Separator{FullArrow: true, PointRight: true},
		activeWindowFont,
		status.String(windowFormat),
		activeWindowSeparatorFont.InvertForeground(),
		status.Separator{FullArrow: true, PointRight: true},
	)

	statusRight := status.Command{
		Name: "$HOME/.dotfiles/bin/gowerline",
		Args: []string{"status-right"},
		Environment: map[string]string{
			"PATH": "$HOME/go/bin:$PATH",
		},
	}
	if debug {
		statusRight.Args = append(statusRight.Args, "--debug")
	}

	statusInterval := "2"
	if debug {
		statusInterval = "10" // reduce polling rate when caching is disabled
	}
	return template.Must(template.New("").Parse(tmuxConfTemplate)).Execute(w, tmuxData{
		Options: map[string]string{
			"status":                       "on",                       // Enable status line.
			"status-interval":              statusInterval,             // Set update interval between generating status lines.
			"status-left":                  statusLeft,                 // Generate left status.
			"status-left-length":           "200",                      // Set maximum width of left status.
			"status-right":                 statusRight.String(),       // Generate right status.
			"status-right-length":          "200",                      // Set maximum width of right status.
			"status-style":                 inactiveWindowFont.Style(), // Set default style like foreground and background color.
			"window-status-bell-style":     bellWindowStyle,            // Generate status for a window that has triggered a bell (BEL).
			"window-status-current-format": currentWindowStatus,        // Generate status for the current window on the left side.
			"window-status-format":         windowStatus,               // Generate status for windows on the left side.
		},
	})
}
