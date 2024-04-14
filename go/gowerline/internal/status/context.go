package status

import (
	"context"
	"io"
	"io/fs"
	"time"
)

type Context struct {
	Cache   SegmentCache
	CacheFS fs.FS
	Context context.Context
	Writer  io.Writer

	now time.Time
}

func (c Context) CacheExpired() bool {
	return c.Cache.ExpiresAt.IsZero() || c.Cache.ExpiresAt.Before(c.now)
}
