package nats_test

import (
	"testing"
	"time"

	"github.com/micro/go-micro/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSocketClose(t *testing.T) {
	t.SkipNow()

	tp, err := setUpTestTransport()
	require.NoError(t, err)

	listener, err := tp.Listen("")
	require.NoError(t, err)
	defer listener.Close()

	result := make(chan error)
	go func() {
		listener.Accept(func(socket transport.Socket) {
			result <- socket.Recv(&transport.Message{})
		})
	}()

	client, err := tp.Dial(listener.Addr())
	require.NoError(t, err)

	require.NoError(t, client.Close())
	select {
	case <-result:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "server socket doesn't get closed")
	}
}
