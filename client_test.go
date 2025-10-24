package openairt //nolint:testpackage // Need to access unexported fields

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	mockToken := "test"
	client := NewClient(mockToken)
	require.Equal(t, client.config.authToken, mockToken)

	config := DefaultConfig(mockToken)
	client = NewClientWithConfig(config)
	require.Equal(t, mockToken, client.config.authToken)
	require.Equal(t, OpenaiRealtimeAPIURLv1, client.config.BaseURL)
	require.Equal(t, APITypeOpenAI, client.config.APIType)
	url := client.getURL("test-model")
	require.Equal(t, OpenaiRealtimeAPIURLv1+"?model=test-model", url)
	headers := client.getHeaders()
	require.Equal(t, "Bearer "+mockToken, headers.Get("Authorization"))

	azureURL := "wss://my-eastus2-openai-resource.openai.azure.com/openai/realtime"
	config = DefaultAzureConfig(mockToken, azureURL)
	client = NewClientWithConfig(config)
	require.Equal(t, mockToken, client.config.authToken)
	require.Equal(t, azureURL, client.config.BaseURL)
	require.Equal(t, APITypeAzure, client.config.APIType)
	require.Equal(t, azureAPIVersion20241001Preview, client.config.APIVersion)
	url = client.getURL("test-model")
	require.Equal(t, azureURL+"?api-version="+azureAPIVersion20241001Preview+"&deployment=test-model", url)
	headers = client.getHeaders()
	require.Equal(t, mockToken, headers.Get("api-key"))
}

func TestUnmarshalCreateClientSecretResponseRealtime(t *testing.T) {
	data := `{
  "value": "ek_68f9c9627aa081919efa4dda52ac3c08",
  "expires_at": 1761201082,
  "session": {
    "type": "realtime",
    "object": "realtime.session",
    "id": "sess_CTizKRJmcL3l4dbpQ8qbI",
    "model": "gpt-realtime-2025-08-28",
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
          "type": "server_vad",
          "threshold": 0.5,
          "prefix_padding_ms": 300,
          "silence_duration_ms": 200,
          "idle_timeout_ms": null,
          "create_response": true,
          "interrupt_response": true
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
	var resp CreateClientSecretResponse
	err := json.Unmarshal([]byte(data), &resp)
	require.NoError(t, err)
	require.Equal(t, "ek_68f9c9627aa081919efa4dda52ac3c08", resp.Value)
	require.Equal(t, int64(1761201082), resp.ExpiresAt)
	session := resp.Session.Realtime
	require.Equal(t, "sess_CTizKRJmcL3l4dbpQ8qbI", session.ID)
	require.Equal(t, "gpt-realtime-2025-08-28", session.Model)
	require.Equal(t, "You are a friendly assistant.", session.Instructions)
	require.Equal(t, []Modality{ModalityAudio}, session.OutputModalities)
	require.Equal(t, "auto", session.ToolChoice.Mode.ToolChoiceType())
	require.Equal(t, Inf, session.MaxOutputTokens)
	require.Equal(t, "auto", session.Truncation.Strategy.TruncationStrategy())
	vad := session.Audio.Input.TurnDetection.ServerVad
	require.Equal(t, "server_vad", string(vad.VadType()))
	require.InDelta(t, 0.5, vad.Threshold, 0.0001)
	require.Equal(t, int64(300), vad.PrefixPaddingMs)
	require.Equal(t, int64(200), vad.SilenceDurationMs)
	require.True(t, vad.CreateResponse)
	require.True(t, vad.InterruptResponse)
	require.Equal(t, int(24000), session.Audio.Input.Format.PCM.Rate)
	require.Equal(t, VoiceAlloy, session.Audio.Output.Voice)
	require.InDelta(t, float32(1.0), session.Audio.Output.Speed, 0.0001)
	require.Equal(t, int(24000), session.Audio.Output.Format.PCM.Rate)
}

func TestUnmarshalCreateClientSecretResponseTranscription(t *testing.T) {
	data := `{
  "value": "ek_68f9cf1281188191903f8dde88960d41",
  "expires_at": 1761202538,
  "session": {
    "type": "transcription",
    "object": "realtime.transcription_session",
    "id": "sess_CTjMo5AmgZuoEWQTvyBqd",
    "expires_at": 0,
    "audio": {
      "input": {
        "format": {
          "type": "audio/pcm",
          "rate": 24000
        },
        "transcription": {
          "model": "gpt-4o-transcribe",
          "language": "en",
          "prompt": null
        },
        "noise_reduction": {
          "type": "near_field"
        },
        "turn_detection": {
          "type": "server_vad",
          "threshold": 0.6,
          "prefix_padding_ms": 300,
          "silence_duration_ms": 500,
          "idle_timeout_ms": null
        }
      }
    },
    "include": null
  }
}`
	var resp CreateClientSecretResponse
	err := json.Unmarshal([]byte(data), &resp)
	require.NoError(t, err)
	require.Equal(t, "ek_68f9cf1281188191903f8dde88960d41", resp.Value)
	require.Equal(t, int64(1761202538), resp.ExpiresAt)
	session := resp.Session.Transcription
	require.Equal(t, "sess_CTjMo5AmgZuoEWQTvyBqd", session.ID)
	require.Equal(t, SessionTypeTranscription, session.Type())
	require.Equal(t, int64(0), session.ExpiresAt)
	require.Equal(t, "gpt-4o-transcribe", session.Audio.Input.Transcription.Model)
	require.Equal(t, "en", session.Audio.Input.Transcription.Language)
	require.Equal(t, int64(300), session.Audio.Input.TurnDetection.ServerVad.PrefixPaddingMs)
	require.Equal(t, int64(500), session.Audio.Input.TurnDetection.ServerVad.SilenceDurationMs)
	require.Equal(t, NoiseReductionNearField, session.Audio.Input.NoiseReduction.Type)
	require.Equal(t, TurnDetectionTypeServerVad, session.Audio.Input.TurnDetection.ServerVad.VadType())
	require.InDelta(t, 0.6, session.Audio.Input.TurnDetection.ServerVad.Threshold, 0.0001)
	require.Equal(t, int(24000), session.Audio.Input.Format.PCM.Rate)
}
