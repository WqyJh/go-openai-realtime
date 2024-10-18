# OpenAI Realtime API SDK for Golang

[![Go Reference](https://pkg.go.dev/badge/github.com/WqyJh/go-openai-realtime.svg)](https://pkg.go.dev/github.com/WqyJh/go-openai-realtime)
[![Go Report Card](https://goreportcard.com/badge/github.com/WqyJh/go-openai-realtime)](https://goreportcard.com/report/github.com/WqyJh/go-openai-realtime)
[![codecov](https://codecov.io/gh/WqyJh/go-openai-realtime/branch/main/graph/badge.svg?token=bCbIfHLIsW)](https://codecov.io/gh/WqyJh/go-openai-realtime)

This library provides unofficial Go clients for [OpenAI Realtime API](https://platform.openai.com/docs/api-reference/realtime). We support all 9 client events and 28 server events.

[Model Support](https://platform.openai.com/docs/models/gpt-4o-realtime):
- gpt-4o-realtime-preview
- gpt-4o-realtime-preview-2024-10-01

## Installation

```bash
go get github.com/WqyJh/go-openai-realtime
```

Currently, go-openai-realtime requires Go version 1.19 or greater.

## Usage

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	openairt "github.com/WqyJh/go-openai-realtime"
)

func main() {
	client := openairt.NewClient("your-api-key")
	conn, err := client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Teletype response
	responseDeltaHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseTextDelta:
			fmt.Printf(event.(openairt.ResponseTextDeltaEvent).Delta)
		case openairt.ServerEventTypeResponseDone:
			fmt.Print("\n\n> ")
		}
	}

	connHandler := openairt.NewConnHandler(context.Background(), conn, responseDeltaHandler)
	connHandler.Start()
	defer connHandler.Stop()

	err = conn.SendMessage(context.Background(), &openairt.SessionUpdateEvent{
		Session: openairt.ClientSession{
			Modalities: []openairt.Modality{openairt.ModalityText},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		conn.SendMessage(context.Background(), &openairt.ConversationItemCreateEvent{
			Item: openairt.MessageItem{
				ID:     openairt.GenerateId("msg_", 10),
				Status: openairt.ItemStatusCompleted,
				Type:   openairt.MessageItemTypeMessage,
				Role:   openairt.MessageRoleUser,
				Content: []openairt.MessageContentPart{
					{
						Type: openairt.MessageContentTypeInputText,
						Text: s.Text(),
					},
				},
			},
		})
		conn.SendMessage(context.Background(), &openairt.ResponseCreateEvent{
			Response: openairt.ResponseCreateParams{
				Modalities:      []openairt.Modality{openairt.ModalityText},
				MaxOutputTokens: 4000,
			},
		})
	}
}
```

## More examples

- [Text-To-Text Example](./examples/text-only/README.md)
- [Text-To-Voice Example](./examples/voice/text-voice/README.md)
- [Voice-To-Voice Example](./examples/voice/voice-voice/README.md)
