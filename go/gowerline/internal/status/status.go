package status

import (
	"context"
	"fmt"
	"io"
	"io/fs"
)

const (
	powerlineArrowPointLeftFull  = ""
	powerlineArrowPointLeftEmpty = ""
)

type Separator struct {
	Font      Font
	FullArrow bool
}

type Segment struct {
	// TODO request minimum delay between updates
	Separator Separator
	Font      Font
	Content   func(Context) error
}

type Context struct {
	Context context.Context
	Writer  io.Writer
}

func (s Segment) WriteTo(ctx Context) error {
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

type Line struct {
	Segments []Segment
}

func (l Line) WriteTo(ctx Context) error {
	for _, segment := range l.Segments {
		err := segment.WriteTo(ctx)
		if err != nil {
			fmt.Fprint(ctx.Writer, err.Error())
		}
	}
	fmt.Fprintln(ctx.Writer)
	return nil
}

func Status(ctx context.Context, w io.Writer, cacheFS fs.FS, statusLine Line) error {
	return statusLine.WriteTo(Context{
		Context: ctx,
		Writer:  w,
	})
}
