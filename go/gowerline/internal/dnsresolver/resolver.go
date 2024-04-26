package dnsresolver

import (
	"context"
	"net"
	"testing"

	"github.com/pkg/errors"
)

type Resolver interface {
	LookupIPWithResolverHost(ctx context.Context, resolverHostPort, hostname string) (net.IP, error)
}

type NetResolver struct{}

func New() *NetResolver {
	return &NetResolver{}
}

func (n *NetResolver) LookupIPWithResolverHost(ctx context.Context, resolverHostPort, hostname string) (net.IP, error) {
	// Equivalent of running: nslookup resolverHostPort hostname
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", resolverHostPort)
		},
	}
	ipAddresses, err := resolver.LookupIPAddr(ctx, hostname)
	if err != nil {
		return nil, err
	}
	if len(ipAddresses) == 0 {
		return nil, errors.Errorf("could not resolve %s @ %s", hostname, resolverHostPort)
	}
	return ipAddresses[0].IP, nil
}

type Test struct {
	tb                  testing.TB
	resolvedIPAddresses []TestIP
}

type TestIP struct {
	ResolverHostPort string
	Hostname         string
	IP               net.IP
}

func NewTest(tb testing.TB, resolvedIPAddresses []TestIP) *Test {
	tb.Helper()
	return &Test{
		tb:                  tb,
		resolvedIPAddresses: resolvedIPAddresses,
	}
}

func (t *Test) LookupIPWithResolverHost(ctx context.Context, resolverHostPort, hostname string) (net.IP, error) {
	for _, ip := range t.resolvedIPAddresses {
		if ip.ResolverHostPort == resolverHostPort && ip.Hostname == hostname {
			return ip.IP, nil
		}
	}
	return nil, errors.Errorf("could not resolve %s @ %s", hostname, resolverHostPort)
}
