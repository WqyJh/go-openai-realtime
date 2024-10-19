package openairt

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

type GorillaWebSocketOptions struct {
	ReadLimit int64
	Dialer    *websocket.Dialer
}

type GorillaWebSocketDialer struct {
	options GorillaWebSocketOptions
}

func NewGorillaWebSocketDialer(options GorillaWebSocketOptions) *GorillaWebSocketDialer {
	if options.Dialer == nil {
		options.Dialer = websocket.DefaultDialer
	}
	return &GorillaWebSocketDialer{
		options: options,
	}
}

func (d *GorillaWebSocketDialer) Dial(ctx context.Context, url string, header http.Header) (WebSocketConn, error) {
	conn, resp, err := d.options.Dialer.DialContext(ctx, url, header)
	if err != nil {
		return nil, err
	}

	conn.SetReadLimit(d.options.ReadLimit)

	return &GorillaWebSocketConn{
		conn:    conn,
		resp:    resp,
		options: d.options,
	}, nil
}

type GorillaWebSocketConn struct {
	conn    *websocket.Conn
	resp    *http.Response
	options GorillaWebSocketOptions
}

func (c *GorillaWebSocketConn) ReadMessage(ctx context.Context) (MessageType, []byte, error) {
	deadline, ok := ctx.Deadline()
	if ok {
		c.conn.SetReadDeadline(deadline)
	}

	// NextReader would block until the message is read or the connection is closed.
	// It won't be canceled by the ctx before its deadline.
	messageType, r, err := c.conn.NextReader()
	if err != nil {
		// The returned error is Permanent, the future read operations on the same connection will not succeed.
		return 0, nil, Permanent(err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return 0, nil, err
	}

	switch messageType {
	case websocket.TextMessage:
		return MessageText, data, nil
	case websocket.BinaryMessage:
		return MessageBinary, data, nil
	default:
		return 0, nil, ErrUnsupportedMessageType
	}
}

func (c *GorillaWebSocketConn) WriteMessage(ctx context.Context, messageType MessageType, data []byte) error {
	deadline, ok := ctx.Deadline()
	if ok {
		c.conn.SetWriteDeadline(deadline)
	}

	switch messageType {
	case MessageText:
		return Permanent(c.conn.WriteMessage(websocket.TextMessage, data))
	case MessageBinary:
		return Permanent(c.conn.WriteMessage(websocket.BinaryMessage, data))
	default:
		return ErrUnsupportedMessageType
	}
}

func (c *GorillaWebSocketConn) Close() error {
	return c.conn.Close()
}

func (c *GorillaWebSocketConn) Response() *http.Response {
	return c.resp
}
