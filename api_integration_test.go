package openairt_test

import (
	"context"
	"os"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime/v2"
	"github.com/stretchr/testify/require"
)

func TestCreateRealtimeSession(t *testing.T) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		t.Skip("OPENAI_API_KEY is not set")
	}
	config := openairt.DefaultConfig(key)
	if baseUrl := os.Getenv("OPENAI_BASE_URL"); baseUrl != "" {
		config.BaseURL = baseUrl
	}
	client := openairt.NewClientWithConfig(config)
	session, err := client.CreateClientSecret(context.Background(), &openairt.CreateClientSecretRequest{
		ExpiresAfter: openairt.ExpiresAfter{
			Anchor:  "created_at",
			Seconds: 600,
		},
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				Model: openairt.GPTRealtime20250828,
				// If you specify ["audio", "text"], you'll get error `Invalid modalities: ['audio', 'text']. Supported combinations are: ['text'] and ['audio'].`
				// That's because Realtime API GA no longer accepts both text and audio for the parameter.
				// Just passing ["audio"] and you can receive transcription of the input/output audio
				// See https://github.com/openai/openai-agents-python/issues/1771#issuecomment-3317018366.
				// OutputModalities: []openairt.Modality{
				// 	openairt.ModalityAudio,
				// },
				Instructions: "You are a friendly assistant.",
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, session.ClientSecret.Value)
	require.NotZero(t, session.ClientSecret.ExpiresAt)
	require.Equal(t, openairt.GPTRealtime20250828, session.Session.Realtime.Model)
	require.Equal(t, "You are a friendly assistant.", session.Session.Realtime.Instructions)
	t.Logf("session: %+v", session)
}

func TestCreateTranscriptionSession(t *testing.T) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		t.Skip("OPENAI_API_KEY is not set")
	}
	config := openairt.DefaultConfig(key)
	if baseUrl := os.Getenv("OPENAI_BASE_URL"); baseUrl != "" {
		config.BaseURL = baseUrl
	}
	client := openairt.NewClientWithConfig(config)
	session, err := client.CreateClientSecret(context.Background(), &openairt.CreateClientSecretRequest{
		ExpiresAfter: openairt.ExpiresAfter{
			Anchor:  "created_at",
			Seconds: 600,
		},
		Session: openairt.SessionUnion{
			Transcription: &openairt.TranscriptionSession{
				Audio: openairt.TranscriptionSessionAudio{
					Input: &openairt.SessionAudioInput{
						Format: openairt.AudioFormatUnion{
							PCM: &openairt.AudioFormatPCM{
								Rate: 24000,
							},
						},
						Transcription: openairt.AudioTranscription{
							Model:    openairt.GPT4oTranscribe,
							Language: "en",
						},
						NoiseReduction: openairt.AudioNoiseReduction{
							Type: openairt.NoiseReductionNearField,
						},
						TurnDetection: openairt.TurnDetectionUnion{
							ServerVad: &openairt.ServerVad{
								Threshold:         0.6,
								PrefixPaddingMs:   300,
								SilenceDurationMs: 500,
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, session.ClientSecret.Value)
	require.NotZero(t, session.ClientSecret.ExpiresAt)
	require.Equal(t, "realtime.transcription_session", session.Session.Transcription.Object)
	require.Equal(t, int(24000), session.Session.Transcription.Audio.Input.Format.PCM.Rate)
	require.Equal(t, openairt.GPT4oTranscribe, session.Session.Transcription.Audio.Input.Transcription.Model)
	require.Equal(t, "en", session.Session.Transcription.Audio.Input.Transcription.Language)
	require.NotNil(t, session.Session.Transcription.Audio.Input.TurnDetection.ServerVad)
	require.Nil(t, session.Session.Transcription.Audio.Input.TurnDetection.SemanticVad)
	require.InEpsilon(t, 0.6, session.Session.Transcription.Audio.Input.TurnDetection.ServerVad.Threshold, 0.0001)
	require.Equal(t, int64(300), session.Session.Transcription.Audio.Input.TurnDetection.ServerVad.PrefixPaddingMs)
	require.Equal(t, int64(500), session.Session.Transcription.Audio.Input.TurnDetection.ServerVad.SilenceDurationMs)
	t.Logf("transcription session: %+v", session)
}
