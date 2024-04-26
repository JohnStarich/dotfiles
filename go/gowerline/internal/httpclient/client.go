package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

var _ Client = &http.Client{}

func New() *http.Client {
	return &http.Client{}
}

type Test struct {
	tb      testing.TB
	handler http.Handler
}

var _ Client = &Test{}

func NewTest(tb testing.TB, handler http.Handler) *Test {
	return &Test{
		tb:      tb,
		handler: handler,
	}
}

func (t *Test) Do(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	t.handler.ServeHTTP(resp, req)
	return resp.Result(), nil
}
