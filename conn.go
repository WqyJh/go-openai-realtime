package openairt

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/coder/websocket"
)

type ServerEventHandler func(ctx context.Context, event ServerEvent)

// Conn is a connection to the OpenAI Realtime API.
type Conn struct {
	conn *websocket.Conn
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "")
}

// SendMessage sends a client event to the server.
func (c *Conn) SendMessage(ctx context.Context, msg ClientEvent) error {
	data, err := MarshalClientEvent(msg)
	if err != nil {
		return err
	}
	return c.conn.Write(ctx, websocket.MessageText, data)
}

// ReadMessage reads a server event from the server.
func (c *Conn) ReadMessage(ctx context.Context) (ServerEvent, error) {
	messageType, data, err := c.conn.Read(ctx)
	if err != nil {
		return nil, err
	}
	if messageType != websocket.MessageText {
		return nil, fmt.Errorf("expected text message, got %d", messageType)
	}
	event, err := UnmarshalServerEvent(data)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// ConnHandler is a handler for a connection to the OpenAI Realtime API.
// It reads messages from the server in a standalone goroutine and calls the registered handlers.
// It is the responsibility of the caller to call Start and Stop.
// The handlers are called in the order they are registered.
// Users should not call ReadMessage directly when using ConnHandler.
type ConnHandler struct {
	ctx      context.Context
	conn     *Conn
	handlers []ServerEventHandler
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewConnHandler creates a new ConnHandler with the given context and connection.
func NewConnHandler(ctx context.Context, conn *Conn, handlers ...ServerEventHandler) *ConnHandler {
	ctx, cancel := context.WithCancel(ctx)
	return &ConnHandler{
		ctx:      ctx,
		conn:     conn,
		handlers: handlers,
		cancel:   cancel,
	}
}

// Start starts the ConnHandler.
func (c *ConnHandler) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.run()
	}()
}

func (c *ConnHandler) run() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		msg, err := c.conn.ReadMessage(c.ctx)
		if err != nil {
			log.Printf("ReadMessage: %v", err)
			continue
		}
		for _, handler := range c.handlers {
			handler(c.ctx, msg)
		}
	}
}

// Stop stops the ConnHandler.
func (c *ConnHandler) Stop() {
	c.cancel()
	c.wg.Wait()
}
