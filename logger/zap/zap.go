package zap

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/micro/go-micro/v2/logger"
)

type zaplog struct {
	cfg  zap.Config
	zap  *zap.Logger
	opts logger.Options
}

func (l *zaplog) Init(opts ...logger.Option) error {
	var err error

	for _, o := range opts {
		o(&l.opts)
	}

	zapConfig := zap.NewProductionConfig()
	if zconfig, ok := l.opts.Context.Value(configKey{}).(zap.Config); ok {
		zapConfig = zconfig
	}

	if zcconfig, ok := l.opts.Context.Value(encoderConfigKey{}).(zapcore.EncoderConfig); ok {
		zapConfig.EncoderConfig = zcconfig

	}

	// Set log Level if not default
	zapConfig.Level = zap.NewAtomicLevel()
	if l.opts.Level != logger.InfoLevel {
		zapConfig.Level.SetLevel(loggerToZapLevel(l.opts.Level))
	}

	log, err := zapConfig.Build()
	if err != nil {
		return err
	}

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		data := []zap.Field{}
		for k, v := range l.opts.Fields {
			data = append(data, zap.Any(k, v))
		}
		log = log.With(data...)
	}

	// Adding namespace
	if namespace, ok := l.opts.Context.Value(namespaceKey{}).(string); ok {
		log = log.With(zap.Namespace(namespace))
	}

	// defer log.Sync() ??

	l.cfg = zapConfig
	l.zap = log

	return nil
}

func (l *zaplog) Fields(fields map[string]interface{}) logger.Logger {
	data := make([]zap.Field, len(fields))
	for k, v := range fields {
		data = append(data, zap.Any(k, v))
	}
	l.zap = l.zap.With(data...)
	return l
}

func (l *zaplog) Error(err error) logger.Logger {
	l.zap = l.zap.With(zap.Error(err))
	return l
}

func (l *zaplog) Log(level logger.Level, args ...interface{}) {
	lvl := loggerToZapLevel(level)
	msg := fmt.Sprint(args...)
	switch lvl {
	case zap.DebugLevel:
		l.zap.Debug(msg)
	case zap.InfoLevel:
		l.zap.Info(msg)
	case zap.WarnLevel:
		l.zap.Warn(msg)
	case zap.ErrorLevel:
		l.zap.Error(msg)
	case zap.PanicLevel:
		l.zap.Panic(msg)
	case zap.FatalLevel:
		l.zap.Fatal(msg)
	}
}

func (l *zaplog) Logf(level logger.Level, format string, args ...interface{}) {
	lvl := loggerToZapLevel(level)
	msg := fmt.Sprintf(format, args...)
	switch lvl {
	case zap.DebugLevel:
		l.zap.Debug(msg)
	case zap.InfoLevel:
		l.zap.Info(msg)
	case zap.WarnLevel:
		l.zap.Warn(msg)
	case zap.ErrorLevel:
		l.zap.Error(msg)
	case zap.PanicLevel:
		l.zap.Panic(msg)
	case zap.FatalLevel:
		l.zap.Fatal(msg)
	}
}

func (l *zaplog) String() string {
	return "zap"
}

func (l *zaplog) Options() logger.Options {
	return l.opts
}

// New builds a new logger based on options
func NewLogger(opts ...logger.Option) (logger.Logger, error) {
	// Default options
	options := logger.Options{
		Level:   logger.InfoLevel,
		Fields:  make(map[string]interface{}),
		Out:     os.Stderr,
		Context: context.Background(),
	}

	l := &zaplog{opts: options}
	if err := l.Init(opts...); err != nil {
		return nil, err
	}

	return l, nil
}

func loggerToZapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.TraceLevel, logger.DebugLevel:
		return zap.DebugLevel
	case logger.InfoLevel:
		return zap.InfoLevel
	case logger.WarnLevel:
		return zap.WarnLevel
	case logger.ErrorLevel:
		return zap.ErrorLevel
	case logger.PanicLevel:
		return zap.PanicLevel
	case logger.FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func zapToLoggerLevel(level zapcore.Level) logger.Level {
	switch level {
	case zap.DebugLevel:
		return logger.DebugLevel
	case zap.InfoLevel:
		return logger.InfoLevel
	case zap.WarnLevel:
		return logger.WarnLevel
	case zap.ErrorLevel:
		return logger.ErrorLevel
	case zap.PanicLevel:
		return logger.PanicLevel
	case zap.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}
