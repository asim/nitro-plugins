package logrus

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/micro/go-micro/v2/logger"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "logrus" {
		t.Errorf("error: name expected 'logrus' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func TestWithFields(t *testing.T) {
	logger.DefaultLogger = NewLogger(logger.WithOutput(os.Stdout))

	logger.Log(logger.InfoLevel, "testing: Info")
	logger.Logf(logger.InfoLevel, "testing: %s", "Infof")
}

func TestJSON(t *testing.T) {
	logger.DefaultLogger = NewLogger(WithJSONFormatter(&logrus.JSONFormatter{}))

	logger.Logf(logger.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	logger.DefaultLogger = NewLogger()

	logger.Init(logger.WithLevel(logger.DebugLevel))
	logger.Logf(logger.DebugLevel, "test show debug: %s", "debug msg")

	logger.Init(logger.WithLevel(logger.InfoLevel))
	logger.Logf(logger.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	logger.DefaultLogger = NewLogger(ReportCaller())

	logger.Logf(logger.InfoLevel, "testing: %s", "WithReportCaller")
}
