package test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type echoServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
	// interval is the interval between each message.
	interval time.Duration
}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		if strings.HasPrefix(k, "X-") {
			for _, v := range v {
				w.Header().Add(k, v)
			}
		}
	}

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer func() {
		err2 := c.CloseNow()
		if err2 != nil {
			s.logf("failed to close connection: %v", err2)
		}
	}()

	for {
		err = echo(r.Context(), c, s.interval)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.logf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

const (
	readTimeout = time.Second * 10
)

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(ctx context.Context, c *websocket.Conn, interval time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	time.Sleep(interval)

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}

type Server struct {
	URL    string
	Server *httptest.Server
}

func NewServer(_ *testing.T, interval time.Duration) Server {
	var testServer Server
	testServer.Server = httptest.NewServer(echoServer{logf: func(f string, v ...interface{}) {
		log.Printf("[test server] "+f, v...)
	}, interval: interval})
	testServer.URL = makeWsProto(testServer.Server.URL)
	return testServer
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}
