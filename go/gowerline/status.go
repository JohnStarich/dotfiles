package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
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
	// TODO request minimum delay between updates
	Separator Separator
	Font      Font
	Content   func(StatusContext) error
}

type StatusContext struct {
	Context context.Context
	Writer  io.Writer
}

func (s StatusSegment) WriteTo(ctx StatusContext) error {
	fmt.Fprint(ctx.Writer, s.Separator.Font)
	fmt.Fprint(ctx.Writer, " ")
	separator := powerlineArrowPointLeftEmpty
	if s.Separator.FullArrow {
		separator = powerlineArrowPointLeftFull
	}
	fmt.Fprint(ctx.Writer, separator)
	fmt.Fprint(ctx.Writer, s.Font)
	fmt.Fprint(ctx.Writer, " ")
	return s.Content(ctx)
}

type StatusLine struct {
	Segments []StatusSegment
}

func (l StatusLine) WriteTo(ctx StatusContext) error {
	for _, segment := range l.Segments {
		err := segment.WriteTo(ctx)
		if err != nil {
			fmt.Fprint(ctx.Writer, err.Error())
		}
	}
	fmt.Fprintln(ctx.Writer)
	return nil
}

func status(ctx context.Context, w io.Writer, cacheFS fs.FS) error {
	statusCtx := StatusContext{
		Context: ctx,
		Writer:  w,
	}
	segments := StatusLine{
		Segments: []StatusSegment{
			{ // weather
				Separator: Separator{Font: Font{Foreground: "#121212"}},
				Font:      Font{Foreground: "#797aac", Background: "#121212"},
				Content:   weatherStatus,
			},
			{ // battery
				Separator: Separator{Font: Font{Foreground: "#f3e6d8", Background: "#121212"}},
				Font:      Font{Foreground: "#f3e6d8", Background: "#121212"},
				Content:   batteryStatus,
			},
			{ // date
				Separator: Separator{Font: Font{Foreground: "#303030", Background: "#121212"}, FullArrow: true},
				Font:      Font{Foreground: "#9e9e9e", Background: "#303030"},
				Content: func(ctx StatusContext) error {
					fmt.Fprint(ctx.Writer, time.Now().Format(time.DateOnly))
					return nil
				},
			},
			{ // time
				Separator: Separator{Font: Font{Foreground: "#626262", Background: "#303030"}},
				Font:      Font{Foreground: "#d0d0d0", Background: "#303030", Bold: true},
				Content: func(ctx StatusContext) error {
					const timeFormat = "3:04 PM"
					fmt.Fprint(ctx.Writer, time.Now().Format(timeFormat))
					return nil
				},
			},
		},
	}

	return segments.WriteTo(statusCtx)
}
