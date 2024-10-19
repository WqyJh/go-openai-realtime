package openairt_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	openairt "github.com/WqyJh/go-openai-realtime"
	test "github.com/WqyJh/go-openai-realtime/test"
	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

func TestCoderWebSocket(t *testing.T) {
	s := test.NewServer(t, time.Millisecond)
	defer s.Server.Close()

	dialer := openairt.NewCoderWebSocketDialer(openairt.CoderWebSocketOptions{})

	conn, err := dialer.Dial(context.Background(), s.URL, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	err = conn.WriteMessage(context.Background(), openairt.MessageBinary+1, []byte("hello"))
	require.ErrorIs(t, err, openairt.ErrUnsupportedMessageType)

	err = conn.WriteMessage(context.Background(), openairt.MessageBinary, []byte("hello"))
	require.NoError(t, err)

	msgType, data, err := conn.ReadMessage(context.Background())
	require.NoError(t, err)
	require.Equal(t, openairt.MessageBinary, msgType)
	require.Equal(t, []byte("hello"), data)

	err = conn.WriteMessage(context.Background(), openairt.MessageText, []byte("world"))
	require.NoError(t, err)

	msgType, data, err = conn.ReadMessage(context.Background())
	require.NoError(t, err)
	require.Equal(t, openairt.MessageText, msgType)
	require.Equal(t, []byte("world"), data)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*5)
	defer cancel()

	_, _, err = conn.ReadMessage(ctx)
	var permanent *openairt.PermanentError
	require.ErrorAs(t, err, &permanent)
	require.True(t, errors.Is(permanent.Err, context.DeadlineExceeded) || errors.Is(permanent.Err, net.ErrClosed))
	t.Logf("permanent error: %+v", permanent.Err)
}

func TestCoderWebSocketReadLimitError(t *testing.T) {
	s := test.NewServer(t, time.Millisecond)
	defer s.Server.Close()

	dialer := openairt.NewCoderWebSocketDialer(openairt.CoderWebSocketOptions{
		ReadLimit: 5,
	})

	conn, err := dialer.Dial(context.Background(), s.URL, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.WriteMessage(context.Background(), openairt.MessageBinary, []byte("hello world"))
	require.NoError(t, err)

	_, _, err = conn.ReadMessage(context.Background())
	require.Error(t, err)
	// 6 = 5 + 1
	// Because SetReadLimit set n+1 as the limit, for 1 byte fin frame.
	require.Contains(t, err.Error(), "read limited at 6 bytes")

	err = conn.Close()
	status := websocket.CloseStatus(err)
	require.Equal(t, websocket.StatusMessageTooBig, status)
}

func TestCoderWebSocketReadLimitOK(t *testing.T) {
	s := test.NewServer(t, time.Millisecond)
	defer s.Server.Close()

	dialer := openairt.NewCoderWebSocketDialer(openairt.CoderWebSocketOptions{
		ReadLimit: 11,
	})

	conn, err := dialer.Dial(context.Background(), s.URL, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	err = conn.WriteMessage(context.Background(), openairt.MessageBinary, []byte("hello world"))
	require.NoError(t, err)

	msgType, data, err := conn.ReadMessage(context.Background())
	require.NoError(t, err)
	require.Equal(t, openairt.MessageBinary, msgType)
	require.Equal(t, []byte("hello world"), data)
}

func TestCoderWebSocketDialOptions(t *testing.T) {
	s := test.NewServer(t, time.Millisecond)
	defer s.Server.Close()

	dialer := openairt.NewCoderWebSocketDialer(openairt.CoderWebSocketOptions{
		DialOptions: &websocket.DialOptions{
			HTTPHeader: http.Header{
				"X-Test":  {"test"},
				"X-Test2": {"test2", "test3"},
			},
		},
	})

	conn, err := dialer.Dial(context.Background(), s.URL, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)

	header := conn.Response().Header
	require.Equal(t, "test", header.Get("X-Test"))
	require.Equal(t, []string{"test2", "test3"}, header["X-Test2"])
}
