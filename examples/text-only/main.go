package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	openairt "github.com/WqyJh/go-openai-realtime/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := openairt.NewClient(os.Getenv("OPENAI_API_KEY"))
	conn, err := client.Connect(ctx, openairt.WithLogger(openairt.StdLogger{}))
	if err != nil {
		log.Fatal(err)
	}

	var sendLock sync.Mutex

	// Teletype response
	responseDeltaHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseOutputTextDelta:
			fmt.Printf(event.(openairt.ResponseOutputTextDeltaEvent).Delta)
		}
	}

	// Full response
	responseHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseDone:
			fmt.Printf("\n\n[full] %s\n\n", event.(openairt.ResponseDoneEvent).Response.Output[0].Assistant.Content[0].Text)
			fmt.Print("> ")
		}
	}

	connHandler := openairt.NewConnHandler(ctx, conn, responseDeltaHandler, responseHandler)
	connHandler.Start()

	err = conn.SendMessage(ctx, &openairt.SessionUpdateEvent{
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				OutputModalities: []openairt.Modality{openairt.ModalityText},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(30 * time.Second):
				sendLock.Lock()
				err := conn.Ping(ctx)
				sendLock.Unlock()
				if err != nil {
					log.Printf("ping error: %v", err)
				} else {
					log.Printf("ping success")
				}
			}
		}
	}()

	fmt.Println("Conversation")
	fmt.Println("---------------------")
	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "exit" || s.Text() == "quit" {
			break
		}
		sendLock.Lock()
		conn.SendMessage(ctx, &openairt.ConversationItemCreateEvent{
			Item: openairt.MessageItemUnion{
				User: &openairt.MessageItemUser{
					ID:     openairt.GenerateID("msg_", 10),
					Status: openairt.ItemStatusCompleted,
					Content: []openairt.MessageContentInput{
						{
							Type: openairt.MessageContentTypeInputText,
							Text: s.Text(),
						},
					},
				},
			},
		})
		conn.SendMessage(ctx, &openairt.ResponseCreateEvent{
			Response: openairt.ResponseCreateParams{
				OutputModalities: []openairt.Modality{openairt.ModalityText},
				MaxOutputTokens:  4000,
			},
		})
		sendLock.Unlock()
	}
	conn.Close()
	err = <-connHandler.Err()
	if err != nil {
		log.Printf("connection error: %v", err)
	}
}
