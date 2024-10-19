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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := openairt.NewClient(os.Getenv("OPENAI_API_KEY"))
	conn, err := client.Connect(ctx, openairt.WithLogger(openairt.StdLogger{}))
	if err != nil {
		log.Fatal(err)
	}

	// Teletype response
	responseDeltaHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseTextDelta:
			fmt.Printf(event.(openairt.ResponseTextDeltaEvent).Delta)
		}
	}

	// Full response
	responseHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseDone:
			fmt.Printf("\n\n[full] %s\n\n", event.(openairt.ResponseDoneEvent).Response.Output[0].Content[0].Text)
			fmt.Print("> ")
		}
	}

	connHandler := openairt.NewConnHandler(ctx, conn, responseDeltaHandler, responseHandler)
	connHandler.Start()

	err = conn.SendMessage(ctx, &openairt.SessionUpdateEvent{
		Session: openairt.ClientSession{
			Modalities: []openairt.Modality{openairt.ModalityText},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conversation")
	fmt.Println("---------------------")
	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "exit" || s.Text() == "quit" {
			break
		}
		conn.SendMessage(ctx, &openairt.ConversationItemCreateEvent{
			Item: openairt.MessageItem{
				ID:     openairt.GenerateID("msg_", 10),
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
		conn.SendMessage(ctx, &openairt.ResponseCreateEvent{
			Response: openairt.ResponseCreateParams{
				Modalities:      []openairt.Modality{openairt.ModalityText},
				MaxOutputTokens: 4000,
			},
		})
	}
	conn.Close()
	err = <-connHandler.Err()
	if err != nil {
		log.Printf("connection error: %v", err)
	}
}
