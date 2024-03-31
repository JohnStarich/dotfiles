package main

import (
	"context"
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

type Separator struct {
	Font      Font
	FullArrow bool
}

type StatusSegment struct {
	Separator Separator
	Font      Font
	Content   func(context.Context, io.Writer) error
}

func (s StatusSegment) WriteTo(ctx context.Context, w io.Writer) error {
	fmt.Fprint(w, s.Separator.Font)
	fmt.Fprint(w, " ")
	separator := powerlineArrowPointLeftEmpty
	if s.Separator.FullArrow {
		separator = powerlineArrowPointLeftFull
	}
	fmt.Fprint(w, separator)
	fmt.Fprint(w, s.Font)
	fmt.Fprint(w, " ")
	return s.Content(ctx, w)
}

type StatusLine struct {
	Segments []StatusSegment
}

func (l StatusLine) WriteTo(ctx context.Context, w io.Writer) error {
	for _, segment := range l.Segments {
		err := segment.WriteTo(ctx, w)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(w)
	return nil
}

func status(ctx context.Context, w io.Writer) error {
	segments := StatusLine{
		Segments: []StatusSegment{
			{ // weather
				Separator: Separator{Font: Font{Foreground: "#121212"}},
				Font:      Font{Foreground: "#797aac", Background: "#121212"},
				Content: func(_ context.Context, w io.Writer) error {
					fmt.Fprint(w, "🌪  57.0°F")
					return nil
				},
			},
			{ // battery
				Separator: Separator{Font: Font{Foreground: "#f3e6d8", Background: "#121212"}},
				Font:      Font{Foreground: "#f3e6d8", Background: "#121212"},
				Content: func(_ context.Context, w io.Writer) error {
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
					fmt.Fprintf(w, "%s %.0f%%", status, chargeNow/totalCharge*100)
					return nil
				},
			},
			{ // date
				Separator: Separator{Font: Font{Foreground: "#303030", Background: "#121212"}, FullArrow: true},
				Font:      Font{Foreground: "#9e9e9e", Background: "#303030"},
				Content: func(_ context.Context, w io.Writer) error {
					fmt.Fprint(w, time.Now().Format(time.DateOnly))
					return nil
				},
			},
			{ // time
				Separator: Separator{Font: Font{Foreground: "#626262", Background: "#303030"}},
				Font:      Font{Foreground: "#d0d0d0", Background: "#303030", Bold: true},
				Content: func(_ context.Context, w io.Writer) error {
					const timeFormat = "3:04 PM"
					fmt.Fprint(w, time.Now().Format(timeFormat))
					return nil
				},
			},
		},
	}

	return segments.WriteTo(ctx, w)
}
