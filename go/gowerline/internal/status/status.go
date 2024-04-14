package status

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"

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
	Font            Font
	GenerateContent func(Context) (time.Duration, error)
	Name            string // required
	Separator       Separator
}

func (s Segment) Status(ctx Context) (SegmentCache, error) {
	fmt.Fprint(ctx.Writer, s.Separator.Font)
	fmt.Fprint(ctx.Writer, " ")
	separator := powerlineArrowPointLeftEmpty
	if s.Separator.FullArrow {
		separator = powerlineArrowPointLeftFull
	}
	fmt.Fprint(ctx.Writer, separator)
	fmt.Fprint(ctx.Writer, s.Font)
	fmt.Fprint(ctx.Writer, " ")

	generatorCtx := ctx
	var newStatus bytes.Buffer
	generatorCtx.Writer = io.MultiWriter(ctx.Writer, &newStatus)

	cacheDuration, err := s.GenerateContent(generatorCtx)
	return SegmentCache{
		Content:   newStatus.String(),
		ExpiresAt: ctx.now.Add(cacheDuration),
	}, err
}

type Line struct {
	Segments []Segment
}

func (l Line) Status(ctx context.Context, w io.Writer, cacheFS fs.FS) error {
	lineCacheData, err := readLineCache(cacheFS)
	if err != nil {
		return err
	}
	now := time.Now()

	for _, segment := range l.Segments {
		segmentCache := lineCacheData.Segments[segment.Name]
		newSegmentCache, err := l.segmentStatus(ctx, segment, w, cacheFS, segmentCache, now)
		if err != nil {
			if !segmentCache.ExpiresAt.IsZero() {
				fmt.Fprint(w, segmentCache.Content)
			}
			fmt.Fprint(w, " <", err.Error(), "> ")
		}
		lineCacheData.Segments[segment.Name] = newSegmentCache
	}
	fmt.Fprintln(w)
	return writeLineCache(cacheFS, lineCacheData)
}

func (l Line) segmentStatus(ctx context.Context, segment Segment, w io.Writer, cacheFS fs.FS, segmentCache SegmentCache, now time.Time) (SegmentCache, error) {
	if segment.Name == "" {
		return SegmentCache{}, errors.New("segment name must be defined")
	}
	if strings.ContainsRune(segment.Name, '/') {
		return SegmentCache{}, errors.Errorf("segment name must not contain a path separator: %q", segment.Name)
	}

	err := hackpadfs.MkdirAll(cacheFS, segment.Name, 0o700)
	if err != nil {
		return SegmentCache{}, err
	}
	subCacheFS, err := hackpadfs.Sub(cacheFS, segment.Name)
	if err != nil {
		return SegmentCache{}, err
	}

	return segment.Status(Context{
		Cache:   segmentCache,
		CacheFS: subCacheFS,
		Context: ctx,
		Writer:  w,
		now:     now,
	})
}
