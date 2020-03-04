package logrus

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/micro/go-micro/v2/logger"
)

type logrusLogger struct {
	*logrus.Logger
	opts Options
}

func (l *logrusLogger) Init(opts ...logger.Option) error {

	for _, o := range opts {
		o(&l.opts.Options)
	}

	if formatter, ok := l.opts.Context.Value(formatterKey{}).(logrus.Formatter); ok {
		l.opts.Formatter = formatter
	}
	if hs, ok := l.opts.Context.Value(hooksKey{}).(logrus.LevelHooks); ok {
		l.opts.Hooks = hs
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}

	log := logrus.New() // defaults
	if ll, ok := l.opts.Context.Value(logrusLoggerKey{}).(*logrus.Logger); ok {
		log = ll

		// overwrite default options
		l.opts.Level = logrusToLoggerLevel(log.GetLevel())
		l.opts.Out = log.Out
		l.opts.Formatter = log.Formatter
		l.opts.Hooks = log.Hooks
		l.opts.ReportCaller = log.ReportCaller
		l.opts.ExitFunc = log.ExitFunc
	} else {
		log.SetLevel(loggerToLogrusLevel(l.opts.Level))
		log.SetOutput(l.opts.Out)
		log.SetFormatter(l.opts.Formatter)
		log.ReplaceHooks(l.opts.Hooks)
		log.SetReportCaller(l.opts.ReportCaller)
		log.ExitFunc = l.opts.ExitFunc
	}

	l.Logger = log

	return nil
}

func (l *logrusLogger) String() string {
	return "logrus"
}

func (l *logrusLogger) Fields(fields map[string]interface{}) logger.Logger {
	// shall we need pool here?
	// but logrus already has pool for its entry.
	return &logrusLogger{logrus.WithFields(fields).Logger, l.opts}
}

func (l *logrusLogger) Error(err error) logger.Logger {
	return &logrusLogger{logrus.WithError(err).Logger, l.opts}
}

func (l *logrusLogger) Log(level logger.Level, args ...interface{}) {
	l.Logger.Log(loggerToLogrusLevel(level), args...)
}

func (l *logrusLogger) Logf(level logger.Level, format string, args ...interface{}) {
	l.Logger.Logf(loggerToLogrusLevel(level), format, args...)
}

func (l *logrusLogger) Options() logger.Options {
	// FIXME: How to return full opts?
	return l.opts.Options
}

// New builds a new logger based on options
func NewLogger(opts ...logger.Option) logger.Logger {
	// Default options
	options := Options{
		Options: logger.Options{
			Level:   logger.InfoLevel,
			Fields:  make(map[string]interface{}),
			Out:     os.Stderr,
			Context: context.Background(),
		},
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
	}
	l := &logrusLogger{opts: options}
	_ = l.Init(opts...)
	return l
}

func loggerToLogrusLevel(level logger.Level) logrus.Level {
	switch level {
	case logger.TraceLevel:
		return logrus.TraceLevel
	case logger.DebugLevel:
		return logrus.DebugLevel
	case logger.InfoLevel:
		return logrus.InfoLevel
	case logger.WarnLevel:
		return logrus.WarnLevel
	case logger.ErrorLevel:
		return logrus.ErrorLevel
	case logger.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func logrusToLoggerLevel(level logrus.Level) logger.Level {
	switch level {
	case logrus.TraceLevel:
		return logger.TraceLevel
	case logrus.DebugLevel:
		return logger.DebugLevel
	case logrus.InfoLevel:
		return logger.InfoLevel
	case logrus.WarnLevel:
		return logger.WarnLevel
	case logrus.ErrorLevel:
		return logger.ErrorLevel
	case logrus.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}
