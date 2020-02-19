package zap

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zaplog struct {
	cfg zap.Config
	zap *zap.Logger
}

func (l *zaplog) Init(opts ...logger.Option) error {
	var err error

	options := &Options{logger.Options{Context: context.Background()}}
	for _, o := range opts {
		o(&options.Options)
	}

	zapConfig := zap.NewProductionConfig()
	if zconfig, ok := options.Context.Value(configKey{}).(zap.Config); ok {
		zapConfig = zconfig
	}

	if zcconfig, ok := options.Context.Value(encoderConfigKey{}).(zapcore.EncoderConfig); ok {
		zapConfig.EncoderConfig = zcconfig

	}

	zapConfig.Level = zap.NewAtomicLevel()
	if level, ok := options.Context.Value(levelKey{}).(logger.Level); ok {
		zapConfig.Level.SetLevel(loggerToZapLevel(level))
	}

	log, err := zapConfig.Build()
	if err != nil {
		return err
	}

	if fields, ok := options.Context.Value(fieldsKey{}).(logger.Fields); ok {
		data := []zap.Field{}
		for k, v := range fields {
			data = append(data, zap.Any(k, v))
		}
		log = log.With(data...)
	}

	if namespace, ok := options.Context.Value(namespaceKey{}).(string); ok {
		log = log.With(zap.Namespace(namespace))
	}

	// defer log.Sync() ??

	l.cfg = zapConfig
	l.zap = log

	return nil
}

func (l *zaplog) SetLevel(level logger.Level) {
	l.cfg.Level.SetLevel(loggerToZapLevel(level))
}

func (l *zaplog) Level() logger.Level {
	return zapToLoggerLevel(l.cfg.Level.Level())
}

func (l *zaplog) Log(level logger.Level, template string, fmtArgs []interface{}, fields logger.Fields) {
	lvl := loggerToZapLevel(level)
	if lvl < zapcore.DPanicLevel && !l.zap.Core().Enabled(lvl) {
		return
	}

	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := l.zap.Check(lvl, msg); ce != nil {
		data := []zap.Field{}
		for k, v := range fields {
			data = append(data, zap.Any(k, v))
		}
		ce.Write(data...)
	}

}

func (l *zaplog) Error(level logger.Level, template string, fmtArgs []interface{}, err error) {
	lvl := loggerToZapLevel(level)
	if lvl < zapcore.DPanicLevel && !l.zap.Core().Enabled(lvl) {
		return
	}

	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := l.zap.Check(lvl, msg); ce != nil {
		ce.Write(zap.Error(err))
	}
}

func (l *zaplog) String() string {
	return "zap"
}

// New builds a new logger based on options
func NewLogger(opts ...logger.Option) (logger.Logger, error) {
	l := &zaplog{}
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
