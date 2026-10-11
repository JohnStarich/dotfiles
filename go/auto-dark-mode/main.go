package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	gsettings, err := exec.LookPath("gsettings")
	if err != nil {
		fmt.Println("gsettings not found, assuming not running on GNOME, so skipping auto-dark-mode")
		return nil
	}
	currentTheme, err := runCommand(ctx, gsettings, "get", "org.gnome.desktop.interface", "color-scheme")
	if err != nil {
		return err
	}
	currentTheme = strings.TrimSpace(currentTheme)
	desiredTheme := desiredTheme(time.Now())
	if currentTheme != desiredTheme {
		if _, err := runCommand(ctx, gsettings, "set", "org.gnome.desktop.interface", "color-scheme", desiredTheme); err != nil {
			return fmt.Errorf("failed to update theme to %s: %w", &desiredTheme, err)
		}
		fmt.Println("Updated theme", currentTheme, "to", desiredTheme)
	}
	return nil
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	output := buf.String()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, output)
	}
	return output, nil
}

func desiredTheme(t time.Time) string {
	// Use 'gsettings describe org.gnome.desktop.interface color-scheme' to get valid values.
	hour := t.Local().Hour()
	const (
		sunrise = 6
		sunset  = 20
	)
	if sunrise < hour && hour < sunset {
		return "'prefer-light'"
	}
	return "'prefer-dark'"
}
