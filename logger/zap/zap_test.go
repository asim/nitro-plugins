package zap

import (
	"testing"

	"github.com/micro/go-micro/v2/logger"
	"github.com/pkg/errors"
)

func TestName(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	if l.String() != "zap" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	logger.SetGlobalLogger(l)

	logger.Infof("test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	logger.SetGlobalLogger(l)

	logger.SetGlobalLevel(logger.DebugLevel)
	logger.Debugf("test show debug: %s", "debug msg")

	logger.SetGlobalLevel(logger.InfoLevel)
	logger.Debugf("test non-show debug: %s", "debug msg")
}

func TestError(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	err2 := errors.Wrap(errors.New("error message"), "from error")
	logger.SetGlobalLogger(l)
	logger.Error("test Error")
	logger.Errorw("test Errorw", err2)
}

func TestFields(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	logger.SetGlobalLogger(l)

	logger.Infow("testing: Fields", logger.Fields{
		"sumo":  "demo",
		"human": true,
		"age":   99,
	})
}

func TestSubLoggerWithFields(t *testing.T) {
	l, err := NewLogger(WithFields(logger.Fields{
		"category": "test",
		"alive":    true,
	}))
	if err != nil {
		t.Fatal(err)
	}

	logger.SetGlobalLogger(l)
	logger.Infow("testing: WithFields", logger.Fields{
		"name":  "demo",
		"human": true,
		"age":   77,
	})
	// Output:
	// {"level":"info","ts":1582075193.56922,"caller":"zap/zap.go:87","msg":"testing: WithFields","category":"test","alive":true,"name":"demo","human":true,"age":77}
}

func TestWithNamespace(t *testing.T) {
	l, err := NewLogger(WithNamespace("micro"))
	if err != nil {
		t.Fatal(err)
	}

	logger.SetGlobalLogger(l)
	logger.Infow("testing: WithFields", logger.Fields{
		"name":  "demo",
		"human": true,
		"age":   77,
	})
	// Output:
	// {"level":"info","ts":1582075193.569254,"caller":"zap/zap.go:87","msg":"testing: WithFields","micro":{"name":"demo","human":true,"age":77}}
}
