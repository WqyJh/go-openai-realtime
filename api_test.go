package openairt_test

import (
	"encoding/json"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/WqyJh/jsontools"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/require"
)

func TestCreateSessionRequest(t *testing.T) {
	data := `{
    "model": "gpt-4o-realtime-preview-2024-12-17",
    "modalities": ["audio", "text"],
    "instructions": "You are a friendly assistant."
}`
	expected := openairt.CreateSessionRequest{
		Model: openairt.GPT4oRealtimePreview20241217,
		ClientSession: openairt.ClientSession{
			Modalities: []openairt.Modality{
				openairt.ModalityAudio,
				openairt.ModalityText,
			},
			Instructions: "You are a friendly assistant.",
		},
	}

	var actual openairt.CreateSessionRequest
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	actualBytes, err := json.Marshal(actual)
	require.NoError(t, err)
	jsontools.RequireJSONEq(t, data, string(actualBytes))
}

func TestCreateSessionResponse(t *testing.T) {
	data := `{
  "id": "sess_001",
  "object": "realtime.session",
  "model": "gpt-4o-realtime-preview-2024-12-17",
  "modalities": ["audio", "text"],
  "instructions": "You are a friendly assistant.",
  "voice": "alloy",
  "input_audio_format": "pcm16",
  "output_audio_format": "pcm16",
  "input_audio_transcription": {
      "model": "whisper-1"
  },
  "turn_detection": null,
  "tools": [],
  "tool_choice": "none",
  "temperature": 0.7,
  "max_response_output_tokens": 200,
  "client_secret": {
    "value": "ek_abc123", 
    "expires_at": 1234567890
  }
}
`
	temperature := float32(0.7)
	expected := openairt.CreateSessionResponse{
		ClientSecret: openairt.ClientSecret{
			Value:     "ek_abc123",
			ExpiresAt: 1234567890,
		},
		ServerSession: openairt.ServerSession{
			ID:     "sess_001",
			Object: "realtime.session",
			Model:  openairt.GPT4oRealtimePreview20241217,
			Modalities: []openairt.Modality{
				openairt.ModalityAudio,
				openairt.ModalityText,
			},
			Instructions:      "You are a friendly assistant.",
			Voice:             openairt.VoiceAlloy,
			InputAudioFormat:  openairt.AudioFormatPcm16,
			OutputAudioFormat: openairt.AudioFormatPcm16,
			InputAudioTranscription: &openairt.InputAudioTranscription{
				Model: openai.Whisper1,
			},
			TurnDetection:   nil,
			Tools:           []openairt.Tool{},
			ToolChoice:      openairt.ServerToolChoice{String: openairt.ToolChoiceNone},
			Temperature:     &temperature,
			MaxOutputTokens: openairt.IntOrInf(200),
		},
	}

	var actual openairt.CreateSessionResponse
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestCreateTranscriptionSessionRequest(t *testing.T) {
	data := `{
    "input_audio_format": "pcm16",
    "input_audio_transcription": {
        "model": "gpt-4o-transcribe",
        "language": "en",
        "prompt": "This is a test transcription"
    },
    "input_audio_noise_reduction": {
        "type": "near_field"
    },
    "modalities": ["text"],
    "turn_detection": {
        "type": "server_vad",
        "threshold": 0.6,
        "prefix_padding_ms": 300,
        "silence_duration_ms": 500
    },
    "include": ["item.input_audio_transcription.logprobs"]
}`
	expected := openairt.CreateTranscriptionSessionRequest{
		InputAudioFormat: openairt.AudioFormatPcm16,
		InputAudioTranscription: &openairt.InputAudioTranscription{
			Model:    "gpt-4o-transcribe",
			Language: "en",
			Prompt:   "This is a test transcription",
		},
		InputAudioNoiseReduction: &openairt.InputAudioNoiseReduction{
			Type: "near_field",
		},
		Modalities: []openairt.Modality{
			openairt.ModalityText,
		},
		TurnDetection: &openairt.ClientTurnDetection{
			Type: openairt.ClientTurnDetectionTypeServerVad,
			TurnDetectionParams: openairt.TurnDetectionParams{
				Threshold:         0.6,
				PrefixPaddingMs:   300,
				SilenceDurationMs: 500,
			},
		},
		Include: []string{"item.input_audio_transcription.logprobs"},
	}

	var actual openairt.CreateTranscriptionSessionRequest
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	actualBytes, err := json.Marshal(actual)
	require.NoError(t, err)
	jsontools.RequireJSONEq(t, data, string(actualBytes))
}

func TestCreateTranscriptionSessionResponse(t *testing.T) {
	data := `{
    "id": "trans_123456",
    "object": "realtime.transcription_session",
    "input_audio_format": "pcm16",
    "input_audio_transcription": {
        "model": "gpt-4o-transcribe",
        "language": "en"
    },
    "modalities": ["text"],
    "turn_detection": {
        "type": "server_vad",
        "threshold": 0.6,
        "prefix_padding_ms": 300,
        "silence_duration_ms": 500
    },
    "client_secret": {
        "value": "ek_trans_abc123",
        "expires_at": 1234567890
    }
}`
	expected := openairt.CreateTranscriptionSessionResponse{
		ID:               "trans_123456",
		Object:           "realtime.transcription_session",
		InputAudioFormat: openairt.AudioFormatPcm16,
		InputAudioTranscription: &openairt.InputAudioTranscription{
			Model:    "gpt-4o-transcribe",
			Language: "en",
		},
		Modalities: []openairt.Modality{
			openairt.ModalityText,
		},
		TurnDetection: &openairt.ServerTurnDetection{
			Type: openairt.ServerTurnDetectionTypeServerVad,
			TurnDetectionParams: openairt.TurnDetectionParams{
				Threshold:         0.6,
				PrefixPaddingMs:   300,
				SilenceDurationMs: 500,
			},
		},
		ClientSecret: openairt.ClientSecret{
			Value:     "ek_trans_abc123",
			ExpiresAt: 1234567890,
		},
	}

	var actual openairt.CreateTranscriptionSessionResponse
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
