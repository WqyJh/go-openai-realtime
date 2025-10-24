package openairt_test

import (
	"strings"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime/v2"
	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	generated := openairt.GenerateID("test_", 10)
	t.Logf("generated: %s", generated)
	require.True(t, strings.HasPrefix(generated, "test_"))
	require.Len(t, generated, 10)

	generated = openairt.GenerateID("test_", 20)
	t.Logf("generated: %s", generated)
	require.True(t, strings.HasPrefix(generated, "test_"))
	require.Len(t, generated, 20)

	generated = openairt.GenerateID("test_", 5)
	t.Logf("generated: %s", generated)
	require.Equal(t, "test_", generated)
}
