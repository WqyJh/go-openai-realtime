package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/WqyJh/go-openai-realtime/examples/voice/pcm"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
	// Open your PCM16 audio file
	file, err := os.Open("../record3/output.pcm")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the entire file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Define the format of your PCM data
	sampleRate := beep.SampleRate(24000) // 24kHz sample rate
	format := beep.Format{
		SampleRate:  sampleRate,
		NumChannels: 1, //  Mono audio
		Precision:   2, // 16-bit PCM
	}

	// Initialize the speaker with the format
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	// Create a new PCMStream
	stream := pcm.NewPCMStream(data, format)

	// Play the stream
	done := make(chan bool)
	speaker.Play(beep.Seq(stream, beep.Callback(func() {
		done <- true
	})))

	<-done
}
