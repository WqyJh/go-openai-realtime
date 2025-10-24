# OpenAI Realtime API SDK for Golang

[![Go Reference](https://pkg.go.dev/badge/github.com/WqyJh/go-openai-realtime/v2.svg)](https://pkg.go.dev/github.com/WqyJh/go-openai-realtime/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/WqyJh/go-openai-realtime/v2)](https://goreportcard.com/report/github.com/WqyJh/go-openai-realtime/v2)
[![codecov](https://codecov.io/gh/WqyJh/go-openai-realtime/branch/main/graph/badge.svg?token=bCbIfHLIsW)](https://codecov.io/gh/WqyJh/go-openai-realtime)

This library provides unofficial Go clients for [OpenAI Realtime API](https://platform.openai.com/docs/api-reference/realtime). We support all 11 client events and 40 server events (except for 3 WebRTC-Only events).

[Model Support](https://platform.openai.com/docs/models/gpt-realtime):

- gpt-realtime
- gpt-realtime-2025-08-28
- gpt-realtime-mini
- gpt-realtime-mini-2025-10-06
- gpt-4o-realtime-preview (shutdown date 2026-02-27)
- gpt-4o-realtime-preview-2024-10-01 (shutdown date 2025-10-10)
- gpt-4o-realtime-preview-2024-12-17 (shutdown date 2026-02-27)
- gpt-4o-realtime-preview-2025-06-03 (shutdown date 2026-02-27)
- gpt-4o-mini-realtime-preview
- gpt-4o-mini-realtime-preview-2024-12-17


## Migration

OpenAI Realtime API GA was officially released on 2025-08-28, and it is not fully compatible with the Beta version.To address these breaking changes, 
we need to develop a v2 version of go-openai-realtime.For details, refer to the [Beta to GA Migration Guide](https://platform.openai.com/docs/guides/realtime#beta-to-ga-migration).

Please use `github.com/WqyJh/go-openai-realtime/v2` instead of `github.com/WqyJh/go-openai-realtime`, and update your code to adapt to the new API specifications.

Note that the Realtime API Beta and Realtime Preview Models (e.g., gpt-4o-realtime-preview) will be shut down on 2026-02-27. After this date, 
you will no longer be able to use these services or the v1 version of go-openai-realtime.

The old v1 version of go-openai-realtime is available at [v1 branch](https://github.com/WqyJh/go-openai-realtime/tree/v1).


## Installation

```bash
go get github.com/WqyJh/go-openai-realtime/v2
```

Currently, go-openai-realtime requires Go version 1.22 or greater.

## Usage

Connect to the OpenAI Realtime API. The default websocket library is [coder/websocket](https://github.com/coder/websocket).

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	openairt "github.com/WqyJh/go-openai-realtime/v2"
)

func main() {
	client := openairt.NewClient("your-api-key")
	conn, err := client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
}
```

<details>
<summary>Use another WebSocket dialer</summary>

Switch to another websocket dialer [gorilla/websocket](https://github.com/gorilla/websocket).

```go
import (
	openairt "github.com/WqyJh/go-openai-realtime/v2"
	gorilla "github.com/WqyJh/go-openai-realtime/v2/contrib/ws-gorilla"
)

func main() {
	dialer := gorilla.NewWebSocketDialer(gorilla.WebSocketOptions{})
	conn, err := client.Connect(ctx, openairt.WithDialer(dialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
```

</details>


<details>
<summary>Send message</summary>

```go

import (
	openairt "github.com/WqyJh/go-openai-realtime/v2"
)

func main() {
    err = conn.SendMessage(ctx, &openairt.SessionUpdateEvent{
        Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				OutputModalities: []openairt.Modality{openairt.ModalityText},
			},
		},
    })
}
```

</details>


<details>
<summary>Read message</summary>

`ReadMessage` is a blocking method that reads the next message from the connection.
It should be called in a standalone goroutine because it's blocking.
If the returned error is Permanent, the future read operations on the same connection will not succeed,
that means the connection is broken and should be closed or had already been closed.

```go
	for {
		msg, err := c.conn.ReadMessage(ctx)
		if err != nil {
			var permanent *PermanentError
			if errors.As(err, &permanent) {
				return permanent.Err
			}
			c.conn.logger.Warnf("read message temporary error: %+v", err)
			continue
		}
		// handle message
	}
```

</details>


<details>
<summary>Subscribe to events</summary>

`ConnHandler` is a helper that reads messages from the server in a standalone goroutine and calls the registered handlers.

Call `openairt.NewConnHandler` to create a `ConnHandler`, then call `Start` to start a new goroutine to read messages.
You can specify only one handler to handle all events or specify multiple handlers.
It's recommended to specify multiple handlers for different purposes.
The registered handlers will be called in the order of registration.

```go
	connHandler := openairt.NewConnHandler(ctx, conn, handler1, handler2, ...)
	connHandler.Start()
```

A handler is function that handle `ServerEvent`. Use `event.ServerEventType()` to determine the type of the event.
Based on the event type, you can get the event data by type assertion.


```go
	// Teletype response
	responseDeltaHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseOutputTextDelta:
			fmt.Printf(event.(openairt.ResponseOutputTextDeltaEvent).Delta)
		}
	}
```


There's no need to `Stop` the `ConnHandler`, it will exit when the connection is closed.
If you want to wait for the `ConnHandler` to exit, you can use `Err()`. This will return
a channel to receive the error.

Note that the receive of the error channel is blocking, so make sure not to call `conn.Close`
after it, which cause deadlock.

```go
    conn.Close()
	err = <-connHandler.Err()
	if err != nil {
		log.Printf("connection error: %v", err)
	}
```


</details>



## More examples

- [Text-To-Text Example](./examples/text-only/README.md)
- [Text-To-Voice Example](./examples/voice/text-voice/README.md)
- [Voice-To-Voice Example](./examples/voice/voice-voice/README.md)

## WebSocket Adapter

The default WebSocket adapter is [coder/websocket](https://github.com/coder/websocket).
There's also a [gorilla/websocket](https://github.com/gorilla/websocket) adapter.
You can easily implement your own adapter by implementing the `WebSocketConn` interface and `WebSocketDialer` interface.

Supported adapters:
- [coder/websocket](./ws_coder.go)
- [gorilla/websocket](./contrib/ws-gorilla)
