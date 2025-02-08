package openairt_test

import (
	"context"
	"net/http"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/require"
)

type mockDialer struct {
	dialFunc func(ctx context.Context, url string, header http.Header) (openairt.WebSocketConn, error)
}

func (m *mockDialer) Dial(ctx context.Context, url string, header http.Header) (openairt.WebSocketConn, error) {
	return m.dialFunc(ctx, url, header)
}

type mockWebSocketConn struct {
	readMessageFunc  func(ctx context.Context) (openairt.MessageType, []byte, error)
	writeMessageFunc func(ctx context.Context, messageType openairt.MessageType, data []byte) error
	closeFunc        func() error
	responseFunc     func() *http.Response
	pingFunc         func(ctx context.Context) error
}

func (m *mockWebSocketConn) ReadMessage(ctx context.Context) (openairt.MessageType, []byte, error) {
	return m.readMessageFunc(ctx)
}

func (m *mockWebSocketConn) WriteMessage(ctx context.Context, messageType openairt.MessageType, data []byte) error {
	return m.writeMessageFunc(ctx, messageType, data)
}

func (m *mockWebSocketConn) Close() error {
	return m.closeFunc()
}

func (m *mockWebSocketConn) Response() *http.Response {
	return m.responseFunc()
}

func (m *mockWebSocketConn) Ping(ctx context.Context) error {
	return m.pingFunc(ctx)
}

func TestConnect(t *testing.T) {
	token := "mock-token"
	model := "test-model"
	client := openairt.NewClient(token)

	dialCalled := false
	readMessageCalled := false
	dialer := &mockDialer{
		dialFunc: func(_ context.Context, url string, header http.Header) (openairt.WebSocketConn, error) {
			require.Equal(t, openairt.OpenaiRealtimeAPIURLv1+"?model="+model, url)
			require.Equal(t, "Bearer "+token, header.Get("Authorization"))
			require.Equal(t, "realtime=v1", header.Get("OpenAI-Beta"))
			dialCalled = true
			return &mockWebSocketConn{
				readMessageFunc: func(_ context.Context) (openairt.MessageType, []byte, error) {
					readMessageCalled = true
					return openairt.MessageText, []byte("test-message"), nil
				},
			}, nil
		},
	}

	conn, err := client.Connect(context.Background(), openairt.WithModel(model), openairt.WithDialer(dialer), openairt.WithLogger(openairt.StdLogger{}))
	require.NoError(t, err)
	require.NotNil(t, conn)
	require.True(t, dialCalled)
	require.False(t, readMessageCalled)
}
