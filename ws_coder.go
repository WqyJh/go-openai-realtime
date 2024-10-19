package openairt

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/coder/websocket"
)

var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
)

type CoderWebSocketOptions struct {
	// ReadLimit is the maximum size of a message in bytes. -1 means no limit. Default is -1.
	ReadLimit int64
	// DialOptions is the options to pass to the websocket.Dial function.
	DialOptions *websocket.DialOptions
}

type CoderWebSocketDialer struct {
	options CoderWebSocketOptions
}

func NewCoderWebSocketDialer(
	options CoderWebSocketOptions,
) *CoderWebSocketDialer {
	// set default read limit
	if options.ReadLimit <= 0 {
		options.ReadLimit = -1
	}
	return &CoderWebSocketDialer{
		options: options,
	}
}

func (d *CoderWebSocketDialer) Dial(ctx context.Context, url string, header http.Header) (WebSocketConn, error) {
	mergedHeader := http.Header{}
	for k, v := range header {
		mergedHeader[k] = append(mergedHeader[k], v...)
	}
	if d.options.DialOptions == nil {
		d.options.DialOptions = &websocket.DialOptions{
			HTTPHeader: mergedHeader,
		}
	} else {
		for k, v := range d.options.DialOptions.HTTPHeader {
			mergedHeader[k] = append(mergedHeader[k], v...)
		}
		d.options.DialOptions.HTTPHeader = mergedHeader
	}

	conn, resp, err := websocket.Dial(ctx, url, d.options.DialOptions)
	if err != nil {
		return nil, err
	}

	conn.SetReadLimit(d.options.ReadLimit)

	return &CoderWebSocketConn{conn: conn, options: d.options, resp: resp}, nil
}

type CoderWebSocketConn struct {
	conn    *websocket.Conn
	resp    *http.Response
	options CoderWebSocketOptions
}

func (c *CoderWebSocketConn) ReadMessage(ctx context.Context) (MessageType, []byte, error) {
	messageType, r, err := c.conn.Reader(ctx)
	if err != nil {
		return 0, nil, Permanent(err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return 0, nil, Permanent(err)
	}

	switch messageType {
	case websocket.MessageText:
		return MessageText, data, nil
	case websocket.MessageBinary:
		return MessageBinary, data, nil
	default:
		return 0, nil, ErrUnsupportedMessageType
	}
}

func (c *CoderWebSocketConn) WriteMessage(ctx context.Context, messageType MessageType, data []byte) error {
	switch messageType {
	case MessageText:
		return Permanent(c.conn.Write(ctx, websocket.MessageText, data))
	case MessageBinary:
		return Permanent(c.conn.Write(ctx, websocket.MessageBinary, data))
	default:
		return ErrUnsupportedMessageType
	}
}

func (c *CoderWebSocketConn) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "")
}

func (c *CoderWebSocketConn) Response() *http.Response {
	return c.resp
}
