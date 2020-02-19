package logrus

import (
	"os"
	"testing"

	"github.com/micro/go-micro/v2/logger"
	"github.com/sirupsen/logrus"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "logrus" {
		t.Errorf("error: name expected 'logrus' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func TestWithFields(t *testing.T) {
	logger.SetGlobalLogger(NewLogger(WithOutput(os.Stdout)))

	logger.Info("testing: Info")
	logger.Infof("testing: %s", "Infof")
	logger.Infow("testing: Infow", logger.Fields{
		"sumo":  "demo",
		"human": true,
		"age":   99,
	})
}

func TestJSON(t *testing.T) {
	logger.SetGlobalLogger(NewLogger(WithJSONFormatter(&logrus.JSONFormatter{})))

	logger.Infof("test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	logger.SetGlobalLogger(NewLogger())

	logger.SetGlobalLevel(logger.DebugLevel)
	logger.Debugf("test show debug: %s", "debug msg")

	logger.SetGlobalLevel(logger.InfoLevel)
	logger.Debugf("test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	logger.SetGlobalLogger(NewLogger(WithReportCaller(true)))

	logger.Infof("testing: %s", "WithReportCaller")
}
