package multi_test

import (
	"errors"
	"testing"

	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/mock"
	"github.com/micro/go-plugins/transport/multi"
)

type deadTransport struct{}

func (f *deadTransport) Dial(addr string, opts ...transport.DialOption) (client transport.Client, err error) {
	return nil, errors.New("Connection refused")
}

func (f *deadTransport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	return nil, errors.New("No port found")
}

func (f *deadTransport) String() string {
	return "dead"
}

func emptyListenerCheck(t *testing.T, transport transport.Transport) {
	t.Helper()

	listener, err := transport.Listen("127.0.0.1")
	if expect := "Not supported"; err.Error() != expect {
		t.Errorf("Got %v expected %v", err, expect)
	}

	if listener != nil {
		t.Error("Expected nil listener on error")
	}
}

func TestNoTransports(t *testing.T) {
	transport := multi.NewTransport()

	emptyListenerCheck(t, transport)

	conn, err := transport.Dial("127.0.0.1")
	if expect := "No transports provided"; err.Error() != expect {
		t.Errorf("Got %v expected %v", err, expect)
	}

	if conn != nil {
		t.Error("expected nil connection on error")
	}
}

func TestNoSuccessfulTransports(t *testing.T) {
	transport := multi.NewTransport(
		multi.WithTransports(
			&deadTransport{},
			&deadTransport{},
		),
	)

	emptyListenerCheck(t, transport)

	conn, err := transport.Dial("127.0.0.1")
	if expect := "Connection refused"; err.Error() != expect {
		t.Errorf("Got %v expected %v", err, expect)
	}

	if conn != nil {
		t.Error("expected nil connection on error")
	}
}

func TestSuccessfulTransports(t *testing.T) {
	mt := mock.NewTransport()

	l, _ := mt.Listen("127.0.0.1:9999")
	go l.Accept(func(sock transport.Socket) {
		sock.Close()
	})
	defer l.Close()

	transport := multi.NewTransport(
		multi.WithTransports(
			&deadTransport{},
			mt,
		),
	)

	emptyListenerCheck(t, transport)

	conn, err := transport.Dial("127.0.0.1:9999")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if conn == nil {
		t.Error("Expected connection to not be nil")
	}

	conn.Close()
}

func TestListnerTransport(t *testing.T) {
	transport := multi.NewTransport(
		multi.WithListenTransport(mock.NewTransport()),
	)

	listener, err := transport.Listen("127.0.0.1")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if listener == nil {
		t.Error("Expected listener to not be nil")
	}

	listener.Close()
}

func TestString(t *testing.T) {
	transport := multi.NewTransport()

	if expect, got := "multi", transport.String(); expect != got {
		t.Errorf("Expected %q got %q", expect, got)
	}
}
