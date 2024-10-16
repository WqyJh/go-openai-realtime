package openairt_test

import (
	"strings"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	generated := openairt.GenerateId("test_", 10)
	t.Logf("generated: %s", generated)
	assert.True(t, strings.HasPrefix(generated, "test_"))
	assert.Equal(t, 10, len(generated))

	generated = openairt.GenerateId("test_", 20)
	t.Logf("generated: %s", generated)
	assert.True(t, strings.HasPrefix(generated, "test_"))
	assert.Equal(t, 20, len(generated))

	generated = openairt.GenerateId("test_", 5)
	t.Logf("generated: %s", generated)
	assert.Equal(t, "test_", generated)
}
