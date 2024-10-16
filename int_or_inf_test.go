package openairt_test

import (
	"encoding/json"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/assert"
)

func TestIntOrInfMarshalJSON(t *testing.T) {
	v := openairt.Int(10)
	assert.Equal(t, 10, v.Value())
	assert.False(t, v.IsInf())

	b, err := json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, []byte("10"), b)

	v = openairt.Inf()
	assert.Equal(t, 0, v.Value())
	assert.True(t, v.IsInf())

	v = openairt.Inf()
	b, err = json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"inf\""), b)
}
