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

func writeTMUXConfig(w io.Writer) error {
	defaultFont := status.Font{
		Foreground: defaultPrimaryColor,
		Background: defaultSecondaryColor,
	}
	activeFont := status.Font{
		Foreground: activeColor,
		Bold:       true,
	}

	statusLeft := fmt.Sprintf(`#{?client_prefix,%s,%s} #{session_name} #{?client_prefix,%s,%s}%s `,
		defaultFont.InvertForeground().VariableSafeString(),
		activeFont.InvertForeground().VariableSafeString(),
		defaultFont.VariableSafeString(),
		activeFont.VariableSafeString(),
		status.Separator{
			FullArrow:  true,
			PointRight: true,
		},
	)

	windowFormat := fmt.Sprintf(`#{window_index}#{?window_flags,#{window_flags}, } %s #{window_name} `, status.Separator{PointRight: true})
	windowStatus := fmt.Sprintf("   %s ", windowFormat)
	currentWindowStatus := fmt.Sprintf(" %s%s %s%s%s%s",
		status.Font{Foreground: defaultSecondaryColor, Background: activeColor, Bold: true},
		status.Separator{
			FullArrow:  true,
			PointRight: true,
		},
		activeFont.InvertForeground(),
		windowFormat,
		activeFont,
		status.Separator{
			FullArrow:  true,
			PointRight: true,
		},
	)

	statusRight := `#(PATH="$HOME/go/bin:$PATH" "$HOME/.dotfiles/bin/gowerline" status-right)`
	return template.Must(template.New("").Parse(tmuxConfTemplate)).Execute(w, tmuxData{
		Options: map[string]string{
			"status":                       "on",                // Enable status line.
			"status-interval":              "2",                 // Set update interval between generating status lines.
			"status-left":                  statusLeft,          // Generate left status.
			"status-left-length":           "200",               // Set maximum width of left status.
			"status-right":                 statusRight,         // Generate right status.
			"status-right-length":          "200",               // Set maximum width of right status.
			"status-style":                 defaultFont.Style(), // Set default style like foreground and background color.
			"window-status-current-format": currentWindowStatus, // Generate status for windows on the left side.
			"window-status-format":         windowStatus,        // Generate status for windows on the left side.
		},
	})
}
