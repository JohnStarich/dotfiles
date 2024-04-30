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
	Options map[string]any
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
	defaultFont := status.Font{Foreground: "colour231", Background: "colour233"}
	statusLeft := fmt.Sprintf(`%s #{session_name} %s `,
		status.Font{Background: "colour31", Bold: true},
		status.Separator{
			Font:       status.Font{Foreground: "colour31"},
			FullArrow:  true,
			PointRight: true,
		},
	)
	windowStatus := fmt.Sprintf(
		`#{window_index}#{?window_flags,#{window_flags}, } #{window_name} %s`,
		status.Separator{
			PointRight: true,
		},
	)
	statusRight := `#(PATH="$HOME/go/bin:$PATH" "$HOME/.dotfiles/bin/gowerline" status-right)`
	return template.Must(template.New("").Parse(tmuxConfTemplate)).Execute(w, tmuxData{
		Options: map[string]any{
			"status":                       "on",                //  Enable status line.
			"status-interval":              "2",                 //  Set update interval between generating status lines.
			"status-left":                  statusLeft,          // Generate left status.
			"status-left-length":           "200",               //  Set maximum width of left status.
			"status-right":                 statusRight,         //  Generate right status.
			"status-right-length":          "200",               //  Set maximum width of right status.
			"status-style":                 defaultFont.Style(), // Set default style like foreground and background color.
			"window-status-current-format": windowStatus,        //  Generate status for windows on the left side.
			"window-status-format":         windowStatus,        //  Generate status for windows on the left side.
		},
	})
}
