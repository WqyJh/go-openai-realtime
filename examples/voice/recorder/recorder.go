package recorder

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/gordonklaus/portaudio"
)

func Record(
	inputChannels int,
	sampleRate float64,
	framesPerBuffer int,
) (data io.Reader, stop func(), err error) {
	var buf bytes.Buffer
	stream, err := portaudio.OpenDefaultStream(inputChannels, 0, sampleRate, framesPerBuffer, func(in []int16) {
		for _, sample := range in {
			binary.Write(&buf, binary.LittleEndian, sample)
		}
	})
	if err != nil {
		return nil, nil, err
	}

	err = stream.Start()
	if err != nil {
		return nil, nil, err
	}

	return &buf, func() {
		stream.Stop()
		stream.Close()
	}, nil
}
