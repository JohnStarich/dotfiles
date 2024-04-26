package status

import (
	"context"
	"io"
	"io/fs"
	"time"

	"github.com/johnstarich/go/gowerline/internal/dnsresolver"
	"github.com/johnstarich/go/gowerline/internal/httpclient"
)

type Context struct {
	Cache      SegmentCache
	CacheFS    fs.FS
	Context    context.Context
	HTTPClient httpclient.Client
	Resolver   dnsresolver.Resolver
	Writer     io.Writer

	now time.Time
}

func NewContext(ctx context.Context, cache SegmentCache, fs fs.FS, writer io.Writer, now time.Time) Context {
	return Context{
		Cache:      cache,
		CacheFS:    fs,
		Context:    ctx,
		HTTPClient: httpclient.New(),
		Resolver:   dnsresolver.New(),
		Writer:     writer,
		now:        now,
	}
}

func (c Context) Now() time.Time {
	return c.now
}

func (c Context) CacheExpired() bool {
	return c.Cache.ExpiresAt.IsZero() || c.Cache.ExpiresAt.Before(c.now)
}
