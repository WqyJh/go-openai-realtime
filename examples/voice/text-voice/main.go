package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/WqyJh/go-openai-realtime/examples/voice/pcm"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
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
	conn, err := client.Connect(ctx, openairt.WithReadLimit(32768*2))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Teletype response
	responseDeltaHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseAudioTranscriptDelta:
			fmt.Printf(event.(openairt.ResponseAudioTranscriptDeltaEvent).Delta)
		}
	}

	// Full response
	responseHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseAudioTranscriptDone:
			fmt.Printf("\n\n[full] %s\n\n", event.(openairt.ResponseAudioTranscriptDoneEvent).Transcript)
			fmt.Print("> ")
		}
	}

	// All event handler
	allHandler := func(ctx context.Context, event openairt.ServerEvent) {
		// data, err := json.Marshal(event)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Printf("[%s] %s\n\n", event.ServerEventType(), string(data))

		// fmt.Printf("[%s]\n\n", event.ServerEventType())
	}

	// Init speaker
	speaker.Init(24000, 24000/10)
	defer speaker.Close()

	streamers := make(chan *pcm.PCMStream, 100)
	datas := make(chan []byte, 100)
	beep.NewBuffer(beep.Format{
		SampleRate:  beep.SampleRate(24000),
		NumChannels: 1,
		Precision:   2, // 16-bit PCM
	})

	// Audio response
	audioResponseHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseAudioDelta:
			msg := event.(openairt.ResponseAudioDeltaEvent)
			// log.Printf("audioResponseHandler: %v", delta)
			data, err := base64.StdEncoding.DecodeString(msg.Delta)
			if err != nil {
				log.Fatal(err)
			}
			// os.WriteFile(fmt.Sprintf("%s.pcm", msg.EventID), data, 0o644)
			// streamer.Append(data)
			datas <- data
		case openairt.ServerEventTypeResponseAudioDone:
			fulldata := []byte{}
			close(datas)
			for data := range datas {
				fulldata = append(fulldata, data...)
			}
			// os.WriteFile(fmt.Sprintf("%s.pcm", event.(openairt.ResponseAudioDoneEvent).EventID), fulldata, 0o644)
			streamer := pcm.NewPCMStream(
				fulldata,
				beep.Format{
					SampleRate:  beep.SampleRate(24000),
					NumChannels: 1,
					Precision:   2, // 16-bit PCM
				},
			)
			datas = make(chan []byte, 100)
			streamers <- streamer
		}
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case streamer := <-streamers:
				done := make(chan bool, 1)
				speaker.Play(beep.Seq(streamer, beep.Callback(func() {
					done <- true
				})))
				<-done // wait for done
			}
		}
	}()

	connHandler := openairt.NewConnHandler(ctx, conn, allHandler, responseDeltaHandler, responseHandler, audioResponseHandler)
	connHandler.Start()
	defer connHandler.Stop()

	err = conn.SendMessage(ctx, &openairt.SessionUpdateEvent{
		Session: openairt.ClientSession{
			Modalities:        []openairt.Modality{openairt.ModalityText, openairt.ModalityAudio},
			Voice:             openairt.VoiceShimmer,
			OutputAudioFormat: openairt.AudioFormatPcm16,
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
		fmt.Println("---------------------")
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
				Modalities:      []openairt.Modality{openairt.ModalityText, openairt.ModalityAudio},
				MaxOutputTokens: 4000,
			},
		})
	}
}
