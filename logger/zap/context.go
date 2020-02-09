package zap

import (
	"context"

	"github.com/micro/go-micro/v2/logger"
)

// setOption returns a function to setup a context with given value
func setOption(k, v interface{}) logger.Option {
	return func(o *logger.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
