package zerolog

import (
	// "errors"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "zerolog" {
		t.Errorf("error: name expected 'zerolog' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func ExampleWithOut() {
	logger.SetGlobalLogger(NewLogger(WithOutput(os.Stdout), WithTimeFormat("ddd"), WithProductionMode()))

	logger.Info("testing: Info")
	logger.Infof("testing: %s", "Infof")
	logger.Infow("testing: Infow", logger.Fields{
		"sumo":  "demo",
		"human": true,
		"age":   99,
	})
	// Output:
	// {"level":"info","time":"ddd","message":"testing: Info"}
	// {"level":"info","time":"ddd","message":"testing: Infof"}
	// {"level":"info","age":99,"human":true,"sumo":"demo","time":"ddd","message":"testing: Infow"}
}

func TestSetLevel(t *testing.T) {
	logger.SetGlobalLogger(NewLogger())

	logger.SetGlobalLevel(logger.DebugLevel)
	logger.Debugf("test show debug: %s", "debug msg")

	logger.SetGlobalLevel(logger.InfoLevel)
	logger.Debugf("test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	logger.SetGlobalLogger(NewLogger(ReportCaller()))

	logger.Infof("testing: %s", "WithReportCaller")
}

func TestWithOutput(t *testing.T) {
	logger.SetGlobalLogger(NewLogger(WithOutput(os.Stdout)))

	logger.Infof("testing: %s", "WithOutput")
}

func TestWithDevelopmentMode(t *testing.T) {
	logger.SetGlobalLogger(NewLogger(WithDevelopmentMode(), WithTimeFormat(time.Kitchen)))

	logger.Infof("testing: %s", "DevelopmentMode")
}

func TestSubLoggerWithMoreFields(t *testing.T) {
	l := NewLogger(WithFields(logger.Fields{
		"component": "gorm",
	}))
	logger.SetGlobalLogger(l)
	logger.Infow("testing: WithFields", logger.Fields{
		"name":  "demo",
		"human": true,
		"age":   77,
	})
	// Output:
	// {"level":"info","component":"gorm","age":77,"human":true,"name":"demo","time":"2020-02-18T12:39:33-08:00","message":"testing: WithFields"}
}

func TestWithError(t *testing.T) {
	logger.SetGlobalLogger(NewLogger())
	err := errors.Wrap(errors.New("error message"), "from error")
	logger.Error("test with error")
	logger.Errorw("test with error", err)
	// Output:
	// {"level":"error","time":"2020-02-18T12:36:13-08:00","message":"test with error"}
	// {"level":"error","stack":[{"func":"TestWithError","line":"85","source":"zerolog_test.go"},{"func":"tRunner","line":"909","source":"testing.go"},{"func":"goexit","line":"1357","source":"asm_amd64.s"}],"error":"from error: error message","time":"2020-02-18T12:36:13-08:00","message":"test with error"}
}

func TestWithHooks(t *testing.T) {
	simpleHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Bool("has_level", level != zerolog.NoLevel)
		e.Str("test", "logged")
	})

	logger.SetGlobalLogger(NewLogger(WithHooks([]zerolog.Hook{simpleHook})))

	logger.Infof("testing: %s", "WithHooks")
}
