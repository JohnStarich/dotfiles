package status

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/hack-pad/hackpadfs"
	"github.com/pkg/errors"
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
	Name      string
	Separator Separator
	Font      Font
	Content   func(Context) error
}

type Context struct {
	Context context.Context
	Writer  io.Writer
	CacheFS fs.FS
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

func (l Line) Status(ctx context.Context, w io.Writer, cacheFS fs.FS) error {
	for _, segment := range l.Segments {
		err := l.segmentStatus(ctx, w, cacheFS, segment)
		if err != nil {
			fmt.Fprint(w, " <", err.Error(), "> ")
		}
	}
	fmt.Fprintln(w)
	return nil
}

func (l Line) segmentStatus(ctx context.Context, w io.Writer, cacheFS fs.FS, segment Segment) error {
	if segment.Name == "" {
		return errors.New("segment name must be defined")
	}
	if strings.ContainsRune(segment.Name, '/') {
		return errors.Errorf("segment name must not contain a path separator: %q", segment.Name)
	}
	err := hackpadfs.MkdirAll(cacheFS, segment.Name, 0o700)
	if err != nil {
		return err
	}
	subCacheFS, err := hackpadfs.Sub(cacheFS, segment.Name)
	if err != nil {
		return err
	}
	return segment.WriteTo(Context{
		Context: ctx,
		Writer:  w,
		CacheFS: subCacheFS,
	})
}
