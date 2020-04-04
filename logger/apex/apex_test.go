package apex

import (
	"testing"

	log "github.com/micro/go-micro/v2/logger"
)

var (
	l = New()
)

func TestName(t *testing.T) {
	if l.String() != "apex" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l.Logf(log.InfoLevel, "test logf: %s", "name")
}

func TestJSON(t *testing.T) {
	l2 := New(WithJSONHandler())
	l2.Logf(log.InfoLevel, "test logf: %s", "name")
}

func TestText(t *testing.T) {
	l2 := New(WithTextHandler())
	l2.Logf(log.InfoLevel, "test logf: %s", "name")
}

func TestCLI(t *testing.T) {
	l2 := New(WithCLIHandler())
	l2.Logf(log.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l.SetLevel(log.DebugLevel)
	l.Logf(log.DebugLevel, "test show debug: %s", "debug msg")

	l.SetLevel(log.InfoLevel)
	l.Logf(log.DebugLevel, "test non-show debug: %s", "debug msg")
}
