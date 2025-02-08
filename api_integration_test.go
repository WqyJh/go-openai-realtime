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
	t.Logf("session: %+v", session)
}
