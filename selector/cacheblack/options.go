package cacheblack

import (
	"time"

	"github.com/micro/go-micro/selector"
	"golang.org/x/net/context"
)

var (
	DefaultCacheTTL      = time.Minute // the cache ttl
	DefaultBlackTTL      = 30          // the ttl to blacklist for
	DefaultBlackErrCount = 3           // number of times we see an error before blacklisting
)

type cacheTTLKey struct{}

// Set the cache ttl
func CacheTTL(t time.Duration) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, cacheTTLKey{}, t)
	}
}

type blackTTLKey struct{}

// Set the ttl to blacklist for
func BlackTTL(n int) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, blackTTLKey{}, n)
	}
}

type blackErrCountKey struct{}

// Set number of times we see an error before blacklisting
func BlackErrCount(n int) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, blackErrCountKey{}, n)
	}
}
