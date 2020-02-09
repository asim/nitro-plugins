package logrus

import (
	"context"
	"io"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/sirupsen/logrus"
)

type formatterKey struct{}
type levelKey struct{}
type outKey struct{}
type hooksKey struct{}
type reportCallerKey struct{}
type exitKey struct{}

type Options struct {
	log.Options
}

func WithTextTextFormatter(formatter *logrus.TextFormatter) log.Option {
	return setOption(formatterKey{}, formatter)
}

func WithJSONFormatter(formatter *logrus.JSONFormatter) log.Option {
	return setOption(formatterKey{}, formatter)
}

func WithLevel(lvl log.Level) log.Option {
	return setOption(levelKey{}, lvl)
}

func WithOut(out io.Writer) log.Option {
	return setOption(outKey{}, out)
}

func WithLevelHooks(hooks logrus.LevelHooks) log.Option {
	return setOption(hooksKey{}, hooks)
}

// warning to use this option. because logrus doest not open CallerDepth option
// this will only print this package
func WithReportCaller(reportCaller bool) log.Option {
	return setOption(reportCallerKey{}, reportCaller)
}

func WithExitFunc(exit func(int)) log.Option {
	return setOption(exitKey{}, exit)
}

func setOption(k, v interface{}) log.Option {
	return func(o *log.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}