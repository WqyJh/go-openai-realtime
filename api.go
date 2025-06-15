package openairt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateSessionRequest struct {
	ClientSession

	// The Realtime model used for this session.
	Model string `json:"model"`
}

type ClientSecret struct {
	// Ephemeral key usable in client environments to authenticate connections to the Realtime API. Use this in client-side environments rather than a standard API token, which should only be used server-side.
	Value string `json:"value"`
	// Timestamp for when the token expires. Currently, all tokens expire after one minute.
	ExpiresAt int64 `json:"expires_at"`
}

type CreateSessionResponse struct {
	ServerSession

	// Ephemeral key returned by the API.
	ClientSecret ClientSecret `json:"client_secret"`
}

// CreateTranscriptionSessionRequest is the request for creating a transcription session.
type CreateTranscriptionSessionRequest struct {
	// The set of items to include in the transcription.
	Include []string `json:"include,omitempty"`
	// The format of input audio. Options are "pcm16", "g711_ulaw", or "g711_alaw".
	InputAudioFormat AudioFormat `json:"input_audio_format,omitempty"`
	// Configuration for input audio noise reduction.
	InputAudioNoiseReduction *InputAudioNoiseReduction `json:"input_audio_noise_reduction,omitempty"`
	// Configuration for input audio transcription.
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`

	// Attention: Keep this field empty! It's shocking that this field is documented but not supported.
	// You may get error of "Unknown parameter: 'modalities'." if this field is not empty.
	// Issue reported: https://community.openai.com/t/unknown-parameter-modalities-when-creating-transcriptionsessions/1150141/6
	// Docs: https://platform.openai.com/docs/api-reference/realtime-sessions/create-transcription#realtime-sessions-create-transcription-modalities
	// The set of modalities the model can respond with. To disable audio, set this to ["text"].
	Modalities []Modality `json:"modalities,omitempty"`

	// Configuration for turn detection.
	TurnDetection *ClientTurnDetection `json:"turn_detection,omitempty"`
}

// CreateTranscriptionSessionResponse is the response from creating a transcription session.
type CreateTranscriptionSessionResponse struct {
	// The unique ID of the session.
	ID string `json:"id"`
	// The object type, must be "realtime.transcription_session".
	Object string `json:"object"`
	// The format of input audio.
	InputAudioFormat AudioFormat `json:"input_audio_format,omitempty"`
	// Configuration of the transcription model.
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`
	// The set of modalities.
	Modalities []Modality `json:"modalities,omitempty"`
	// Configuration for turn detection.
	TurnDetection *ServerTurnDetection `json:"turn_detection,omitempty"`
	// Ephemeral key returned by the API.
	ClientSecret ClientSecret `json:"client_secret"`
}

type OpenAIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Param      string `json:"param"`
	Code       any    `json:"code"`
}

type ErrorResponse struct { //nolint:errname // this is a http error response
	StatusCode  int `json:"-"`
	OpenAIError `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.OpenAIError.Message
}

type httpOption struct {
	client  *http.Client
	headers http.Header
	method  string
}

type HTTPOption func(*httpOption)

func WithHeaders(headers http.Header) HTTPOption {
	return func(o *httpOption) {
		o.headers = headers
	}
}

func WithClient(client *http.Client) HTTPOption {
	return func(o *httpOption) {
		o.client = client
	}
}

func WithMethod(method string) HTTPOption {
	return func(o *httpOption) {
		o.method = method
	}
}

func HTTPDo[Q any, R any](ctx context.Context, url string, req *Q, opts ...HTTPOption) (*R, error) {
	opt := httpOption{
		client:  http.DefaultClient,
		headers: http.Header{},
		method:  http.MethodPost,
	}
	for _, o := range opts {
		o(&opt)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, opt.method, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header = opt.headers

	response, err := opt.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http failed: %w", err)
	}
	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		err = json.Unmarshal(data, &errResp)
		if err != nil {
			return nil, fmt.Errorf("http status code: %d, error: %s", response.StatusCode, string(data))
		}
		errResp.StatusCode = response.StatusCode
		return nil, &errResp
	}

	var resp R
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}
