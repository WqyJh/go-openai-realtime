package openairt

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/coder/websocket"
)

type ServerEventHandler func(ctx context.Context, event ServerEvent)

type Conn struct {
	conn *websocket.Conn
}

func (c *Conn) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "")
}

func (c *Conn) SendMessage(ctx context.Context, msg ClientEvent) error {
	data, err := MarshalClientEvent(msg)
	if err != nil {
		return err
	}
	return c.conn.Write(ctx, websocket.MessageText, data)
}

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

type ConnHandler struct {
	ctx      context.Context
	conn     *Conn
	handlers []ServerEventHandler
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	err      error
}

func NewConnHandler(ctx context.Context, conn *Conn, handlers ...ServerEventHandler) *ConnHandler {
	ctx, cancel := context.WithCancel(ctx)
	return &ConnHandler{
		ctx:      ctx,
		conn:     conn,
		handlers: handlers,
		cancel:   cancel,
	}
}

func (c *ConnHandler) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.err = c.run()
	}()
}

func (c *ConnHandler) run() error {
	for {
		select {
		case <-c.ctx.Done():
			return nil
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

func (c *ConnHandler) Stop() {
	c.cancel()
	c.wg.Wait()
}

func (c *ConnHandler) Err() error {
	return c.err
}
