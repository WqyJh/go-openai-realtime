package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	openairt "github.com/WqyJh/go-openai-realtime"
)

func main() {
	socksProxy := os.Getenv("SOCKS_PROXY")
	if socksProxy != "" {
		proxyURL, err := url.Parse(socksProxy)
		if err != nil {
			log.Fatal(err)
		}
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := openairt.NewClient(os.Getenv("OPENAI_API_KEY"))
	conn, err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

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
	defer connHandler.Stop()

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
		conn.SendMessage(ctx, &openairt.ResponseCreateEvent{
			Response: openairt.ResponseCreateParams{
				Modalities:      []openairt.Modality{openairt.ModalityText},
				MaxOutputTokens: 4000,
			},
		})
	}
}
