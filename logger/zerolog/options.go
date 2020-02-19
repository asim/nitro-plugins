package zerolog

import (
	"context"
	"io"

	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
)

type Options struct {
	logger.Options
}

type reportCallerKey struct{}

func ReportCaller() logger.Option {
	return setOption(reportCallerKey{}, true)
}

type useAsDefaultKey struct{}

func UseAsDefault() logger.Option {
	return setOption(useAsDefaultKey{}, true)
}

type developmentModeKey struct{}

func WithDevelopmentMode() logger.Option {
	return setOption(developmentModeKey{}, true)
}

type productionModeKey struct{}

func WithProductionMode() logger.Option {
	return setOption(productionModeKey{}, true)
}

type outputKey struct{}

func WithOutput(out io.Writer) logger.Option {
	return setOption(outputKey{}, out)
}

type fieldsKey struct{}

func WithFields(fields logger.Fields) logger.Option {
	return setOption(fieldsKey{}, fields)
}

type levelKey struct{}

func WithLevel(lvl logger.Level) logger.Option {
	return setOption(levelKey{}, lvl)
}

type timeFormatKey struct{}

func WithTimeFormat(timeFormat string) logger.Option {
	return setOption(timeFormatKey{}, timeFormat)
}

type hooksKey struct{}

func WithHooks(hooks []zerolog.Hook) logger.Option {
	return setOption(hooksKey{}, hooks)
}

type exitKey struct{}

func WithExitFunc(exit func(int)) logger.Option {
	return setOption(exitKey{}, exit)
}

func setOption(k, v interface{}) logger.Option {
	return func(o *logger.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
