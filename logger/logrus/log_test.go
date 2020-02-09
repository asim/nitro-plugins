package logrus

import (
	"testing"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/sirupsen/logrus"
)

var (
	l = New()
)

func TestName(t *testing.T) {
	if l.String() != "logrus" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l.Logf(log.InfoLevel, "test logf: %s", "name")
}

func TestJSON(t *testing.T) {
	l2 := New(WithJSONFormatter(&logrus.JSONFormatter{}))
	l2.Logf(log.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l.SetLevel(log.DebugLevel)
	l.Logf(log.DebugLevel, "test show debug: %s", "debug msg")

	l.SetLevel(log.InfoLevel)
	l.Logf(log.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestReportCaller(t *testing.T) {
	l2 := New(WithReportCaller(true))
	l2.Logf(log.InfoLevel, "test logf: %s", "name")
}