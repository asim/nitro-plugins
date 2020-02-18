package zerolog

import (
	// "errors"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/logger/log"
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
	log.SetGlobalLogger(NewLogger(WithOutput(os.Stdout), WithTimeFormat("ddd"), WithProductionMode()))

	log.Info("testing: Info")
	log.Infof("testing: %s", "Infof")
	log.Infow("testing: Infow", map[string]interface{}{
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
	log.SetGlobalLogger(NewLogger())

	log.SetGlobalLevel(logger.DebugLevel)
	log.Debugf("test show debug: %s", "debug msg")

	log.SetGlobalLevel(logger.InfoLevel)
	log.Debugf("test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	log.SetGlobalLogger(NewLogger(ReportCaller()))

	log.Infof("testing: %s", "WithReportCaller")
}

func TestWithOutput(t *testing.T) {
	log.SetGlobalLogger(NewLogger(WithOutput(os.Stdout)))

	log.Infof("testing: %s", "WithOutput")
}

func TestWithDevelopmentMode(t *testing.T) {
	log.SetGlobalLogger(NewLogger(WithDevelopmentMode(), WithTimeFormat(time.Kitchen)))

	log.Infof("testing: %s", "DevelopmentMode")
}

func TestWithFields(t *testing.T) {
	log.SetGlobalLogger(NewLogger())

	log.Infow("testing: WithFields", map[string]interface{}{
		"sumo":  "demo",
		"human": true,
		"age":   99,
	})
}

func TestWithError(t *testing.T) {
	l := NewLogger(WithFields(map[string]interface{}{
		"name":  "sumo",
		"age":   99,
		"alive": true,
	}))
	err := errors.Wrap(errors.New("error message"), "from error")
	log.SetGlobalLogger(l)
	log.Error("test with error")
	log.Errorw("test with error", err)
	// Output:
	// {"level":"error","age":99,"alive":true,"name":"sumo","time":"2020-02-18T03:11:42-08:00","message":"test with error"}
	// {"level":"error","age":99,"alive":true,"name":"sumo","stack":[{"func":"TestWithError","line":"86","source":"zerolog_test.go"},{"func":"tRunner","line":"909","source":"testing.go"},{"func":"goexit","line":"1357","source":"asm_amd64.s"}],"error":"from error: error message","time":"2020-02-18T03:11:42-08:00","message":"test with error"}
}

func TestWithHooks(t *testing.T) {
	simpleHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Bool("has_level", level != zerolog.NoLevel)
		e.Str("test", "logged")
	})

	log.SetGlobalLogger(NewLogger(WithHooks([]zerolog.Hook{simpleHook})))

	log.Infof("testing: %s", "WithHooks")
}
