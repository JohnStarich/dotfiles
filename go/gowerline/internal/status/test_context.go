package status

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/hack-pad/hackpadfs/mem"
	"github.com/johnstarich/go/gowerline/internal/dnsresolver"
	"github.com/johnstarich/go/gowerline/internal/httpclient"
)

type TestConfig struct {
	Handler             http.Handler
	Now                 time.Time
	ResolvedIPAddresses []dnsresolver.TestIP
}

type TestContext struct {
	Context
}

func NewTestContext(tb testing.TB, config TestConfig) TestContext {
	tb.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	tb.Cleanup(cancel)
	memFS, err := mem.NewFS()
	if err != nil {
		tb.Fatal(err)
	}
	return TestContext{
		Context{
			Cache:      SegmentCache{},
			CacheFS:    memFS,
			Context:    ctx,
			HTTPClient: httpclient.NewTest(tb, config.Handler),
			Resolver:   dnsresolver.NewTest(tb, config.ResolvedIPAddresses),
			Writer:     bytes.NewBuffer(nil),
			now:        config.Now,
		},
	}
}

func (t TestContext) Output() string {
	return t.Context.Writer.(*bytes.Buffer).String()
}
