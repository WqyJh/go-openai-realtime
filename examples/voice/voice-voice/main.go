package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/WqyJh/go-openai-realtime/examples/voice/pcm"
	"github.com/WqyJh/go-openai-realtime/examples/voice/recorder"
	openairt "github.com/WqyJh/go-openai-realtime/v2"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gordonklaus/portaudio"
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
			// ignore
		case openairt.ServerEventTypeResponseOutputAudioTranscriptDone:
			fmt.Printf("[response] %s\n", event.(openairt.ResponseOutputAudioTranscriptDoneEvent).Transcript)
		}
	}

	// Full response
	responseHandler := func(ctx context.Context, event openairt.ServerEvent) {
		switch event.ServerEventType() {
		case openairt.ServerEventTypeConversationItemInputAudioTranscriptionCompleted:
			question := event.(openairt.ConversationItemInputAudioTranscriptionCompletedEvent).Transcript
			fmt.Printf("\n[question] %s\n", question)
		case openairt.ServerEventTypeResponseDone:
			fmt.Print("\n> ")
		}
	}

	// Log handler
	logHandler := func(ctx context.Context, event openairt.ServerEvent) {
		// fmt.Printf("[%s]\n", event.ServerEventType())

		switch event.ServerEventType() {
		case openairt.ServerEventTypeError,
			openairt.ServerEventTypeSessionUpdated:
			// data, err := json.Marshal(event)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Printf("[%s] %s\n", event.ServerEventType(), string(data))
		}
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

	connHandler := openairt.NewConnHandler(ctx, conn, logHandler, responseHandler, responseDeltaHandler, audioResponseHandler)
	connHandler.Start()

	err = conn.SendMessage(ctx, openairt.SessionUpdateEvent{
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				OutputModalities: []openairt.Modality{openairt.ModalityAudio},
				Audio: openairt.RealtimeSessionAudio{
					Input: &openairt.SessionAudioInput{
						Transcription: openairt.AudioTranscription{
							Model: openairt.GPTRealtime,
						},
					},
					Output: &openairt.SessionAudioOutput{
						Voice: openairt.VoiceMarin,
						Format: openairt.AudioFormatUnion{
							PCM: &openairt.AudioFormatPCM{
								Rate: 24000,
							},
						},
					},
				},
				MaxOutputTokens: 4000,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Conversation (input 'start' to start recording)")
	fmt.Println("---------------------")
	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)

	// Audio init
	portaudio.Initialize()
	defer portaudio.Terminate()

	var data io.Reader
	var stop func()

Loop:
	for s.Scan() {

		switch s.Text() {
		case "exit", "quit":
			break Loop
		case "start":
			data, stop, err = recorder.Record(1, 24000, 64)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Recording...(input 'stop' to stop recording)")
			fmt.Print("> ")
		case "stop":
			if stop == nil {
				log.Print("not recording")
				continue
			}
			stop()
			audioData, err := io.ReadAll(data)
			if err != nil {
				log.Fatal(err)
			}
			base64Audio := base64.StdEncoding.EncodeToString(audioData)
			conn.SendMessage(ctx, openairt.InputAudioBufferAppendEvent{
				Audio: base64Audio,
			})
			conn.SendMessage(ctx, openairt.InputAudioBufferCommitEvent{})
			conn.SendMessage(ctx, openairt.ResponseCreateEvent{
				Response: openairt.ResponseCreateParams{
					OutputModalities: []openairt.Modality{openairt.ModalityAudio},
				},
			})
		}

	}
}
