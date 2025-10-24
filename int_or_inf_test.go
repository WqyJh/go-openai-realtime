package openairt_test

import (
	"encoding/json"
	"math"
	"math/big"
	"net"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime/v2"
	"github.com/stretchr/testify/require"
)

func TestIntOrInfMarshalJSON(t *testing.T) {
	v := openairt.IntOrInf(10)
	require.Equal(t, 10, int(v))
	require.False(t, v.IsInf())

	b, err := json.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, []byte("10"), b)

	v = openairt.Inf
	require.Equal(t, math.MaxInt, int(v))
	require.True(t, v.IsInf())

	b, err = json.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, []byte("\"inf\""), b)

	type TestStruct struct {
		MaxOutputTokens openairt.IntOrInf `json:"max_output_tokens,omitempty"`
	}

	s := TestStruct{}
	b, err = json.Marshal(s)
	require.NoError(t, err)
	require.Equal(t, []byte("{}"), b)
}

func TestInetAton(t *testing.T) {
	ip := "37.57.95.27"
	ret := big.NewInt(0)
	i := net.ParseIP(ip).To4()
	require.NotNil(t, i)
	ret.SetBytes(i)
	t.Log(ret.Uint64())
}
