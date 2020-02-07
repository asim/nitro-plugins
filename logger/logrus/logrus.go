package logrus

import (
	"io"
	"os"

	log "github.com/micro/go-micro/logger"
	"github.com/sirupsen/logrus"
)

var (
	formatter    logrus.Formatter = new(logrus.TextFormatter)
	lvl                           = logrus.InfoLevel
	out          io.Writer        = os.Stderr
	hooks                         = make(logrus.LevelHooks)
	reportCaller                  = false
	exit                          = os.Exit
)

type logger struct {
	*logrus.Logger
}

func (l *logger) Fields(fields ...log.Field) log.Logger {
	data := make(logrus.Fields, len(fields))
	for _, f := range fields {
		data[f.Key] = f.GetValue()
	}

	// shall we need pool here?
	// but logrus already has pool for its entry.
	return &logger{logrus.WithFields(data).Logger}
}

func (l *logger) Init(opts ...log.Option) error {
	options := &Options{}
	for _, o := range opts {
		o(&options.Options)
	}

	if options.Context != nil {
		f, ok := options.Context.Value(formatterKey{}).(logrus.Formatter)
		if ok {
			formatter = f
		}

		l, ok := options.Context.Value(levelKey{}).(log.Level)
		if ok {
			lvl = convertToLogrusLevel(l)
		}

		o, ok := options.Context.Value(outKey{}).(io.Writer)
		if ok {
			out = o
		}

		h, ok := options.Context.Value(hooksKey{}).(logrus.LevelHooks)
		if ok {
			hooks = h
		}

		r, ok := options.Context.Value(reportCallerKey{}).(bool)
		if ok {
			if r == true {
			}
			reportCaller = r
		}

		e, ok := options.Context.Value(exitKey{}).(func(int))
		if ok {
			exit = e
		}
	}

	l.Logger = &logrus.Logger{
		Out:          out,
		Formatter:    formatter,
		Hooks:        hooks,
		Level:        lvl,
		ExitFunc:     exit,
		ReportCaller: reportCaller,
	}

	return nil
}

func (l *logger) SetLevel(level log.Level) {
	l.Logger.SetLevel(convertToLogrusLevel(level))
}

func (l *logger) Level() log.Level {
	return convertLevel(l.Logger.Level)
}

func (l *logger) Log(level log.Level, args ...interface{}) {
	l.Logger.Log(convertToLogrusLevel(level), args...)
}

func (l *logger) Logf(level log.Level, format string, args ...interface{}) {
	l.Logger.Logf(convertToLogrusLevel(level), format, args...)
}

func (l *logger) String() string {
	return "logrus"
}

// New builds a new logger based on options
func New(opts ...log.Option) log.Logger {
	l := &logger{}
	_ = l.Init(opts...)
	return l
}

func convertToLogrusLevel(level log.Level) logrus.Level {
	switch level {
	case log.TraceLevel:
		return logrus.TraceLevel
	case log.DebugLevel:
		return logrus.DebugLevel
	case log.InfoLevel:
		return logrus.InfoLevel
	case log.WarnLevel:
		return logrus.WarnLevel
	case log.ErrorLevel:
		return logrus.ErrorLevel
	case log.PanicLevel:
		return logrus.PanicLevel
	case log.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func convertLevel(level logrus.Level) log.Level {
	switch level {
	case logrus.TraceLevel:
		return log.TraceLevel
	case logrus.DebugLevel:
		return log.DebugLevel
	case logrus.InfoLevel:
		return log.InfoLevel
	case logrus.WarnLevel:
		return log.WarnLevel
	case logrus.ErrorLevel:
		return log.ErrorLevel
	case logrus.PanicLevel:
		return log.PanicLevel
	case logrus.FatalLevel:
		return log.FatalLevel
	default:
		return log.InfoLevel
	}
}
