package openairt_test

import (
	"encoding/json"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime/v2"
	"github.com/stretchr/testify/require"
)

func TestToolChoiceUnion(t *testing.T) {
	data := `{"type":"function","name":"get_current_weather"}`
	expectedFunction := openairt.ToolChoiceFunction{
		Name: "get_current_weather",
	}
	expected := openairt.ToolChoiceUnion{Function: &expectedFunction}
	actual := openairt.ToolChoiceUnion{}
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, expectedFunction, *actual.Function)

	data = `"auto"`
	expected = openairt.ToolChoiceUnion{Mode: openairt.ToolChoiceModeAuto}
	actual = openairt.ToolChoiceUnion{}
	err = json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, openairt.ToolChoiceModeAuto, actual.Mode)

	data = `"none"`
	expected = openairt.ToolChoiceUnion{Mode: openairt.ToolChoiceModeNone}
	actual = openairt.ToolChoiceUnion{}
	err = json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, openairt.ToolChoiceModeNone, actual.Mode)

	data = `"required"`
	expected = openairt.ToolChoiceUnion{Mode: openairt.ToolChoiceModeRequired}
	actual = openairt.ToolChoiceUnion{}
	err = json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, openairt.ToolChoiceModeRequired, actual.Mode)
}

func TestTurnDetectionUnion(t *testing.T) {
	t.Run("unmarshal server_vad", func(t *testing.T) {
		data := `{"type":"server_vad","threshold":0.5,"prefix_padding_ms":300,"silence_duration_ms":200,"idle_timeout_ms":null,"create_response":true,"interrupt_response":true}`
		expected := openairt.TurnDetectionUnion{
			ServerVad: &openairt.ServerVad{
				Threshold:         0.5,
				PrefixPaddingMs:   300,
				SilenceDurationMs: 200,
				IdleTimeoutMs:     0,
				CreateResponse:    true,
				InterruptResponse: true,
			},
		}
		actual := openairt.TurnDetectionUnion{}
		err := json.Unmarshal([]byte(data), &actual)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("unmarshal null", func(t *testing.T) {
		// Unmarshaling JSON "null" should produce an empty TurnDetectionUnion.
		data := `null`
		expected := openairt.TurnDetectionUnion{}
		actual := openairt.TurnDetectionUnion{}
		err := json.Unmarshal([]byte(data), &actual)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("marshal null", func(t *testing.T) {
		// Marshaling an empty TurnDetectionUnion should produce JSON "null"
		// to disable turn detection.
		input := openairt.TurnDetectionUnion{}
		data, err := json.Marshal(input)
		require.NoError(t, err)
		require.Equal(t, "null", string(data))
	})
}
