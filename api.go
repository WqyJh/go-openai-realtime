package openairt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClientSecret struct {
	// Ephemeral key usable in client environments to authenticate connections to the Realtime API. Use this in client-side environments rather than a standard API token, which should only be used server-side.
	Value string `json:"value"`
	// Timestamp for when the token expires. Currently, all tokens expire after one minute.
	ExpiresAt int64 `json:"expires_at"`
}

type ExpiresAfter struct {
	// The anchor point for the client secret expiration, meaning that seconds will be added to the created_at time of the client secret to produce an expiration timestamp. Only created_at is currently supported.
	Anchor string `json:"anchor,omitzero"`

	// The number of seconds from the anchor point to the expiration. Select a value between 10 and 7200 (2 hours). This default to 600 seconds (10 minutes) if not specified.
	Seconds int `json:"seconds,omitzero"`
}

type CreateClientSecretRequest struct {
	// Configuration for the client secret expiration. Expiration refers to the time after which a client secret will no longer be valid for creating sessions. The session itself may continue after that time once started. A secret can be used to create multiple sessions until it expires.
	ExpiresAfter ExpiresAfter `json:"expires_after,omitzero"`

	// Session configuration to use for the client secret. Choose either a realtime session or a transcription session.
	Session SessionUnion `json:"session,omitzero"`
}

type CreateClientSecretResponse struct {
	// Ephemeral key returned by the API.
	ClientSecret

	// Session configuration to use for the client secret. Choose either a realtime session or a transcription session.
	Session SessionUnion `json:"session,omitzero"`
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
