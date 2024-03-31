package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	powerlineArrowPointLeftFull  = ""
	powerlineArrowPointLeftEmpty = ""
)

func status(w io.Writer) error {
	// hide first arrow?
	// weather
	fmt.Fprint(w, FontConfig{
		Foreground: "#121212",
		Background: "default",
	})
	fmt.Fprint(w, " ")
	fmt.Fprint(w, powerlineArrowPointLeftFull)
	fmt.Fprint(w, FontConfig{
		Foreground: "#797aac",
		Background: "#121212",
	})
	fmt.Fprint(w, "🌪  57.0°F")

	// battery
	fmt.Fprint(w, FontConfig{
		Foreground: "#f3e6d8",
		Background: "#121212",
	})
	fmt.Fprint(w, " ")
	fmt.Fprint(w, powerlineArrowPointLeftEmpty)
	fmt.Fprint(w, FontConfig{
		Foreground: "#f3e6d8",
		Background: "#121212",
	})
	fmt.Fprint(w, " ")
	chargeNowBytes, err := os.ReadFile("/sys/class/power_supply/BAT1/charge_now")
	if err != nil {
		return err
	}
	totalChargeBytes, err := os.ReadFile("/sys/class/power_supply/BAT1/charge_full")
	if err != nil {
		return err
	}
	chargeNow, err := strconv.ParseFloat(strings.TrimSpace(string(chargeNowBytes)), 64)
	if err != nil {
		return err
	}
	totalCharge, err := strconv.ParseFloat(strings.TrimSpace(string(totalChargeBytes)), 64)
	if err != nil {
		return err
	}
	statusBytes, err := os.ReadFile("/sys/class/power_supply/BAT1/status")
	if err != nil {
		return err
	}
	status := "⚠️"
	if strings.TrimSpace(string(statusBytes)) == "Discharging" {
		status = "🔥"
	}
	fmt.Printf("%s %.0f%%", status, chargeNow/totalCharge*100)

	// date
	fmt.Fprint(w, FontConfig{
		Foreground: "#303030",
		Background: "#121212",
	})
	fmt.Fprint(w, " ")
	fmt.Fprint(w, powerlineArrowPointLeftFull)
	fmt.Fprint(w, FontConfig{
		Foreground: "#9e9e9e",
		Background: "#303030",
	})
	fmt.Fprint(w, " ")
	fmt.Fprint(w, time.Now().Format(time.DateOnly))

	// time
	fmt.Fprint(w, FontConfig{
		Foreground: "#626262",
		Background: "#303030",
	})
	fmt.Fprint(w, " ")
	fmt.Fprint(w, powerlineArrowPointLeftEmpty)
	fmt.Fprint(w, FontConfig{
		Foreground: "#d0d0d0",
		Background: "#303030",
		Bold:       true,
	})
	fmt.Fprint(w, " ")
	const timeFormat = "3:04 PM"
	fmt.Fprint(w, time.Now().Format(timeFormat))
	fmt.Fprintln(w)
	return nil
}

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
