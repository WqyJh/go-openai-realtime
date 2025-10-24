package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	openairt "github.com/WqyJh/go-openai-realtime/v2"
	"github.com/WqyJh/go-openai-realtime/v2/examples/voice/pcm"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
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
		case openairt.ServerEventTypeResponseOutputAudioTranscriptDelta:
			fmt.Printf(event.(openairt.ResponseOutputAudioTranscriptDeltaEvent).Delta)
		}
	}

	// Full response
	responseHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeResponseOutputAudioTranscriptDone:
			fmt.Printf("\n\n[full] %s\n\n", event.(openairt.ResponseOutputAudioTranscriptDoneEvent).Transcript)
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
		case openairt.ServerEventTypeResponseOutputAudioDelta:
			msg := event.(openairt.ResponseOutputAudioDeltaEvent)
			// log.Printf("audioResponseHandler: %v", delta)
			data, err := base64.StdEncoding.DecodeString(msg.Delta)
			if err != nil {
				log.Fatal(err)
			}
			// os.WriteFile(fmt.Sprintf("%s.pcm", msg.EventID), data, 0o644)
			// streamer.Append(data)
			datas <- data
		case openairt.ServerEventTypeResponseOutputAudioDone:
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

	err = conn.SendMessage(ctx, &openairt.SessionUpdateEvent{
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				OutputModalities: []openairt.Modality{openairt.ModalityAudio},
				Audio: &openairt.RealtimeSessionAudio{
					Output: &openairt.SessionAudioOutput{
						Voice: openairt.VoiceMarin,
						Format: &openairt.AudioFormatUnion{
							PCM: &openairt.AudioFormatPCM{
								Rate: 24000,
							},
						},
					},
				},
			},
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
				OutputModalities: []openairt.Modality{openairt.ModalityAudio},
				MaxOutputTokens:  4000,
			},
		})
	}
}
