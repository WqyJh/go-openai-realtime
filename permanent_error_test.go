package openairt_test

import (
	"errors"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/require"
)

func TestPermanent(t *testing.T) {
	original := errors.New("test")
	err := openairt.Permanent(original)
	require.ErrorIs(t, err, original)
	var permanent *openairt.PermanentError
	require.True(t, errors.As(err, &permanent))
	require.ErrorIs(t, permanent.Err, original)
	require.EqualError(t, err, original.Error())
}
