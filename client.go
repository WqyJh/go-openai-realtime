package openairt

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

const (
	GPTRealtime             = "gpt-realtime"
	GPTRealtime20250828     = "gpt-realtime-2025-08-28"
	GPTRealtimeMini         = "gpt-realtime-mini"
	GPTRealtimeMini20251006 = "gpt-realtime-mini-2025-10-06"

	// Deprecated: shutdown date 2026-02-27.
	GPT4oRealtimePreview = "gpt-4o-realtime-preview"
	// Deprecated: shutdown date 2025-10-10.
	GPT4oRealtimePreview20241001 = "gpt-4o-realtime-preview-2024-10-01"
	// Deprecated: shutdown date 2026-02-27.
	GPT4oRealtimePreview20241217 = "gpt-4o-realtime-preview-2024-12-17"
	// Deprecated: shutdown date 2026-02-27.
	GPT4oRealtimePreview20250603 = "gpt-4o-realtime-preview-2025-06-03"

	GPT4oMiniRealtimePreview         = "gpt-4o-mini-realtime-preview"
	GPT4oMiniRealtimePreview20241217 = "gpt-4o-mini-realtime-preview-2024-12-17"
)

// Transcription models.
const (
	GPT4oTranscribe     = "gpt-4o-transcribe"
	GPT4oMiniTranscribe = "gpt-4o-mini-transcribe"
	Whisper1            = "whisper-1"
)

// Client is OpenAI Realtime API client.
type Client struct {
	config ClientConfig
}

// NewClient creates new OpenAI Realtime API client.
func NewClient(authToken string) *Client {
	config := DefaultConfig(authToken)
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates new OpenAI Realtime API client for specified config.
func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) getURL(model string) string {
	query := url.Values{}
	query.Set("model", model)

	return c.config.BaseURL + "?" + query.Encode()
}

func (c *Client) getHeaders() http.Header {
	headers := http.Header{}

	if c.config.APIType == APITypeAzure {
		headers.Set("api-key", c.config.authToken)
	} else {
		headers.Set("Authorization", "Bearer "+c.config.authToken)
	}
	return headers
}

type connectOption struct {
	model  string
	intent string
	dialer WebSocketDialer
	logger Logger
}

type ConnectOption func(*connectOption)

// WithModel sets the model to use for the connection.
func WithModel(model string) ConnectOption {
	return func(opts *connectOption) {
		opts.model = model
	}
}

// Set transcription intent instead of model.
func WithIntent() ConnectOption {
	return func(opts *connectOption) {
		opts.intent = "transcription"
	}
}

// WithDialer sets the dialer for the connection.
func WithDialer(dialer WebSocketDialer) ConnectOption {
	return func(opts *connectOption) {
		opts.dialer = dialer
	}
}

// WithLogger sets the logger for the connection.
func WithLogger(logger Logger) ConnectOption {
	return func(opts *connectOption) {
		opts.logger = logger
	}
}

// Connect connects to the OpenAI Realtime API.
func (c *Client) Connect(ctx context.Context, opts ...ConnectOption) (*Conn, error) {
	connectOpts := connectOption{
		model:  GPT4oRealtimePreview,
		logger: NopLogger{},
	}
	for _, opt := range opts {
		opt(&connectOpts)
	}
	if connectOpts.dialer == nil {
		connectOpts.dialer = DefaultDialer()
	}

	// default headers
	headers := c.getHeaders()

	// get url by model
	var url string
	if connectOpts.intent == "" { //nolint:gocritic // if conditions would be determined in order
		url = c.getURL(connectOpts.model)
	} else if c.config.APIType != APITypeOpenAI {
		return nil, errors.New("intent not supported for Azure API type")
	} else {
		url = c.config.BaseURL + "?" + "intent=" + connectOpts.intent
	}

	// dial
	conn, err := connectOpts.dialer.Dial(ctx, url, headers)
	if err != nil {
		return nil, err
	}

	return &Conn{conn: conn, logger: connectOpts.logger}, nil
}

func (c *Client) getAPIHeaders() http.Header {
	headers := http.Header{}

	if c.config.APIType == APITypeAzure {
		headers.Set("api-key", c.config.authToken)
	} else {
		headers.Set("Authorization", "Bearer "+c.config.authToken)
	}
	headers.Set("Content-Type", "application/json")
	return headers
}

func (c *Client) CreateClientSecret(
	ctx context.Context,
	req *CreateClientSecretRequest,
) (*CreateClientSecretResponse, error) {
	return HTTPDo[CreateClientSecretRequest, CreateClientSecretResponse](
		ctx,
		c.config.APIBaseURL+"/realtime/client_secrets",
		req,
		WithClient(c.config.HTTPClient),
		WithMethod(http.MethodPost),
		WithHeaders(c.getAPIHeaders()),
	)
}
