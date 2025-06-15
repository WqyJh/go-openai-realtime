package openairt_test

import (
	"context"
	"os"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		t.Skip("OPENAI_API_KEY is not set")
	}
	client := openairt.NewClient(key)
	session, err := client.CreateSession(context.Background(), &openairt.CreateSessionRequest{
		Model: openairt.GPT4oRealtimePreview20241217,
		ClientSession: openairt.ClientSession{
			Modalities: []openairt.Modality{
				openairt.ModalityAudio,
				openairt.ModalityText,
			},
			Instructions: "You are a friendly assistant.",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, session.ClientSecret.Value)
	require.NotZero(t, session.ClientSecret.ExpiresAt)
	require.Equal(t, openairt.GPT4oRealtimePreview20241217, session.Model)
	require.Equal(t, "You are a friendly assistant.", session.Instructions)
	t.Logf("session: %+v", session)
}

func TestCreateTranscriptionSession(t *testing.T) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		t.Skip("OPENAI_API_KEY is not set")
	}
	client := openairt.NewClient(key)
	session, err := client.CreateTranscriptionSession(context.Background(), &openairt.CreateTranscriptionSessionRequest{
		InputAudioFormat: openairt.AudioFormatPcm16,
		InputAudioTranscription: &openairt.InputAudioTranscription{
			Model:    openairt.GPT4oTranscribe,
			Language: "en",
		},
		InputAudioNoiseReduction: &openairt.InputAudioNoiseReduction{
			Type: openairt.NearFieldNoiseReduction,
		},
		// Attention: Keep this field empty! It's shocking that this field is documented but not supported.
		// Modalities: []openairt.Modality{
		// 	openairt.ModalityText,
		// },
		TurnDetection: &openairt.ClientTurnDetection{
			Type: openairt.ClientTurnDetectionTypeServerVad,
			TurnDetectionParams: openairt.TurnDetectionParams{
				Threshold:         0.6,
				PrefixPaddingMs:   300,
				SilenceDurationMs: 500,
			},
		},
		Include: []string{},
	})
	require.NoError(t, err)
	require.NotEmpty(t, session.ClientSecret.Value)
	require.NotZero(t, session.ClientSecret.ExpiresAt)
	require.Equal(t, "realtime.transcription_session", session.Object)
	require.Equal(t, openairt.AudioFormatPcm16, session.InputAudioFormat)
	require.Equal(t, openairt.GPT4oTranscribe, session.InputAudioTranscription.Model)
	require.Equal(t, "en", session.InputAudioTranscription.Language)
	require.Equal(t, openairt.ServerTurnDetectionTypeServerVad, session.TurnDetection.Type)
	require.InEpsilon(t, 0.6, session.TurnDetection.Threshold, 0.0001)
	require.Equal(t, 300, session.TurnDetection.PrefixPaddingMs)
	require.Equal(t, 500, session.TurnDetection.SilenceDurationMs)
	require.Empty(t, session.Modalities)
	t.Logf("transcription session: %+v", session)
}
