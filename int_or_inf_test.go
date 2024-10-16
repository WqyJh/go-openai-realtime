package openairt_test

import (
	"encoding/json"
	"math"
	"math/big"
	"net"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/assert"
)

func TestIntOrInfMarshalJSON(t *testing.T) {
	v := openairt.IntOrInf(10)
	assert.Equal(t, 10, int(v))
	assert.False(t, v.IsInf())

	b, err := json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, []byte("10"), b)

	v = openairt.Inf
	assert.Equal(t, math.MaxInt, int(v))
	assert.True(t, v.IsInf())

	b, err = json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"inf\""), b)

	type TestStruct struct {
		MaxOutputTokens openairt.IntOrInf `json:"max_output_tokens,omitempty"`
	}

	s := TestStruct{}
	b, err = json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, []byte("{}"), b)
}

func TestInetAton(t *testing.T) {
	ip := "37.57.95.27"
	ret := big.NewInt(0)
	i := net.ParseIP(ip).To4()
	assert.NotNil(t, i)
	ret.SetBytes(i)
	t.Log(ret.Uint64())
}
