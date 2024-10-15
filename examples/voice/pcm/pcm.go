package pcm

import (
	"encoding/binary"

	"github.com/faiface/beep"
)

type PCMStream struct {
	data   []byte
	pos    int
	format beep.Format
}

func (s *PCMStream) Stream(samples [][2]float64) (n int, ok bool) {
	bytesPerSample := s.format.Precision
	for i := range samples {
		if s.pos+bytesPerSample > len(s.data) {
			return i, false
		}

		// Read PCM16 sample
		sample := int16(binary.LittleEndian.Uint16(s.data[s.pos:]))
		s.pos += bytesPerSample

		// Convert to float64 sample
		floatSample := float64(sample) / (1 << 15)

		// Since it's mono, we copy the same sample to both left and right channels
		samples[i][0] = floatSample
		samples[i][1] = floatSample
	}
	return len(samples), true
}

func (s *PCMStream) Err() error {
	return nil
}

func NewPCMStream(data []byte, format beep.Format) *PCMStream {
	return &PCMStream{
		data:   data,
		format: format,
	}
}
