package nats_test

import (
	"testing"
	"time"

	"github.com/micro/go-micro/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleMessageEcho(t *testing.T) {
	tp, err := setUpTestTransport(t)
	require.NoError(t, err)

	listener, err := tp.Listen("")
	require.NoError(t, err)
	defer listener.Close()

	go func() {
		listener.Accept(func(socket transport.Socket) {
			m := &transport.Message{}
			require.NoError(t, socket.Recv(m))
			require.NoError(t, socket.Send(m))
		})
	}()

	client, err := tp.Dial(listener.Addr())
	require.NoError(t, err)

	sendMessage := &transport.Message{
		Header: map[string]string{"key": "value"},
		Body:   []byte("test"),
	}
	recvMessage := &transport.Message{}
	require.NoError(t, client.Send(sendMessage))
	require.NoError(t, client.Recv(recvMessage))
	assert.Equal(t, sendMessage, recvMessage)

	listener.Close()
	time.Sleep(1 * time.Second)
}

func BenchmarkSimpleMessageEcho(b *testing.B) {
	tp, err := setUpTestTransport(b)
	require.NoError(b, err)

	listener, err := tp.Listen("")
	require.NoError(b, err)
	defer listener.Close()

	go func() {
		listener.Accept(func(socket transport.Socket) {
			m := &transport.Message{}
			require.NoError(b, socket.Recv(m))
			require.NoError(b, socket.Send(m))
		})
	}()

	sendMessage := &transport.Message{
		Header: map[string]string{"key": "value"},
		Body:   []byte("test"),
	}
	recvMessage := &transport.Message{}

	b.ResetTimer()
	for i := 0; i < b.N/10; i++ {
		client, err := tp.Dial(listener.Addr())
		require.NoError(b, err)

		require.NoError(b, client.Send(sendMessage))
		require.NoError(b, client.Recv(recvMessage))
		assert.Equal(b, sendMessage, recvMessage)
	}
}
