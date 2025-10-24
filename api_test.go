package openairt_test

import (
	"encoding/json"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/WqyJh/jsontools"
	"github.com/stretchr/testify/require"
)

func TestCreateClientSecretRequest(t *testing.T) {
	data := `{
  "expires_after": {
    "anchor": "created_at",
    "seconds": 600
  },
  "session": {
    "type": "realtime",
    "model": "gpt-realtime",
    "instructions": "You are a friendly assistant."
  }
}`
	expected := openairt.CreateClientSecretRequest{
		ExpiresAfter: openairt.ExpiresAfter{
			Anchor:  "created_at",
			Seconds: 600,
		},
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				Model:        openairt.GPTRealtime,
				Instructions: "You are a friendly assistant.",
			},
		},
	}

	var actual openairt.CreateClientSecretRequest
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	actualBytes, err := json.Marshal(actual)
	require.NoError(t, err)
	jsontools.RequireJSONEq(t, data, string(actualBytes))
}

func TestCreateClientSecretResponse(t *testing.T) {
	data := `{
  "value": "ek_68af296e8e408191a1120ab6383263c2",
  "expires_at": 1756310470,
  "session": {
    "type": "realtime",
    "object": "realtime.session",
    "id": "sess_C9CiUVUzUzYIssh3ELY1d",
    "model": "gpt-realtime",
    "output_modalities": ["audio"],
    "instructions": "You are a friendly assistant.",
    "tools": [],
    "tool_choice": "auto",
    "max_output_tokens": "inf",
    "tracing": null,
    "truncation": "auto",
    "prompt": null,
    "expires_at": 0,
    "audio": {
      "input": {
        "format": {
          "type": "audio/pcm",
          "rate": 24000
        },
        "transcription": null,
        "noise_reduction": null,
        "turn_detection": {
          "type": "server_vad"
        }
      },
      "output": {
        "format": {
          "type": "audio/pcm",
          "rate": 24000
        },
        "voice": "alloy",
        "speed": 1.0
      }
    },
    "include": null
  }
}`
	expected := openairt.CreateClientSecretResponse{
		ClientSecret: openairt.ClientSecret{
			Value:     "ek_68af296e8e408191a1120ab6383263c2",
			ExpiresAt: 1756310470,
		},
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				ID:     "sess_C9CiUVUzUzYIssh3ELY1d",
				Object: "realtime.session",
				Model:  openairt.GPTRealtime,
				OutputModalities: []openairt.Modality{
					openairt.ModalityAudio,
				},
				Instructions: "You are a friendly assistant.",
				Truncation: openairt.TruncationUnion{
					Strategy: openairt.TruncationStrategyAuto,
				},
				Audio: openairt.RealtimeSessionAudio{
					Input: &openairt.SessionAudioInput{
						Format: openairt.AudioFormatUnion{
							PCM: &openairt.AudioFormatPCM{
								Rate: 24000,
							},
						},
						TurnDetection: openairt.TurnDetectionUnion{
							ServerVad: &openairt.ServerVad{},
						},
					},
					Output: &openairt.SessionAudioOutput{
						Format: openairt.AudioFormatUnion{
							PCM: &openairt.AudioFormatPCM{
								Rate: 24000,
							},
						},
						Voice: "alloy",
						Speed: 1.0,
					},
				},
				Tools:           []openairt.ToolUnion{},
				ToolChoice:      openairt.ToolChoiceUnion{Mode: openairt.ToolChoiceModeAuto},
				MaxOutputTokens: openairt.Inf,
			},
		},
	}

	var actual openairt.CreateClientSecretResponse
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
