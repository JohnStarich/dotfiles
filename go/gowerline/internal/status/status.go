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
	powerlineArrowPointLeftFull   = ""
	powerlineArrowPointLeftEmpty  = ""
	powerlineArrowPointRightFull  = ""
	powerlineArrowPointRightEmpty = ""
)

type Separator struct {
	FullArrow  bool
	PointRight bool
}

func (s Separator) String() string {
	switch {
	case !s.FullArrow && !s.PointRight:
		return powerlineArrowPointLeftEmpty
	case !s.FullArrow && s.PointRight:
		return powerlineArrowPointRightEmpty
	case s.FullArrow && !s.PointRight:
		return powerlineArrowPointLeftFull
	case s.FullArrow && s.PointRight:
		return powerlineArrowPointRightFull
	default:
		return powerlineArrowPointLeftEmpty
	}
}

type Segment struct {
	Font            Font
	GenerateContent func(Context) (time.Duration, error)
	Name            string // required
	SeparatorFont   Font
	Separator       Separator
}

func (s Segment) Status(ctx Context) (SegmentCache, error) {
	fmt.Fprint(ctx.Writer, " ")
	fmt.Fprint(ctx.Writer, s.SeparatorFont.String())
	fmt.Fprint(ctx.Writer, s.Separator.String())
	fmt.Fprint(ctx.Writer, s.Font)
	fmt.Fprint(ctx.Writer, " ")

	generatorCtx := ctx
	var newStatus bytes.Buffer
	generatorCtx.Writer = io.MultiWriter(ctx.Writer, &newStatus)

	cacheDuration, err := generateContentSafely(generatorCtx, s)
	return SegmentCache{
		Content:   newStatus.String(),
		ExpiresAt: ctx.now.Add(cacheDuration),
	}, err
}

func generateContentSafely(ctx Context, segment Segment) (_ time.Duration, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered from panic: %v", r)
		}
	}()
	return segment.GenerateContent(ctx)
}

type Line struct {
	Segments []Segment
}

func (l Line) Status(ctx context.Context, writer, errWriter io.Writer, cacheFS fs.FS) error {
	lineCacheData, _ := readLineCache(cacheFS)
	// failed to read cache, proceed with a blank cache
	now := time.Now()

	segmentCacheUpdates := make(map[string]SegmentCache)
	for _, segment := range l.Segments {
		segmentCache := lineCacheData.Segments[segment.Name]
		newSegmentCache, err := l.segmentStatus(ctx, segment, writer, cacheFS, segmentCache, now)
		if err != nil {
			fmt.Fprintln(errWriter, err.Error())
		} else if newSegmentCache.ExpiresAt != now {
			segmentCacheUpdates[segment.Name] = newSegmentCache
		}
	}
	if len(segmentCacheUpdates) > 0 {
		for name, segment := range segmentCacheUpdates {
			lineCacheData.Segments[name] = segment
		}
		return writeLineCache(cacheFS, lineCacheData)
	}
	fmt.Fprintln(writer)
	return nil
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

	return segment.Status(NewContext(ctx, segmentCache, subCacheFS, w, now))
}
