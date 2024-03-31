package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
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
