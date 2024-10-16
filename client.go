package openairt

import (
	"context"
	"net/http"
	"net/url"

	"github.com/coder/websocket"
)

const (
	GPT4oRealtimePreview         = "gpt-4o-realtime-preview"
	GPT4oRealtimePreview20241001 = "gpt-4o-realtime-preview-2024-10-01"
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

func (c *Client) getUrl(model string) string {
	query := url.Values{}

	if c.config.APIType == APITypeAzure {
		query.Set("api-version", c.config.APIVersion)
		query.Set("deployment", model)
	} else {
		query.Set("model", model)
	}

	return c.config.BaseURL + "?" + query.Encode()
}

func (c *Client) getHeaders() http.Header {
	headers := http.Header{}

	if c.config.APIType == APITypeAzure {
		headers.Set("api-key", c.config.authToken)
	} else {
		headers.Set("Authorization", "Bearer "+c.config.authToken)
		headers.Set("OpenAI-Beta", "realtime=v1")
	}
	return headers
}

type connectOption struct {
	model       string
	dialOptions *websocket.DialOptions
	readLimit   int64
}

type ConnectOption func(*connectOption)

// WithModel sets the model to use for the connection.
func WithModel(model string) ConnectOption {
	return func(opts *connectOption) {
		opts.model = model
	}
}

// WithDialOptions sets the dial options for the connection.
func WithDialOptions(dialOptions *websocket.DialOptions) ConnectOption {
	return func(opts *connectOption) {
		opts.dialOptions = dialOptions
	}
}

// WithReadLimit sets the read limit for the connection.
func WithReadLimit(limit int64) ConnectOption {
	return func(opts *connectOption) {
		opts.readLimit = limit
	}
}

// Connect connects to the OpenAI Realtime API.
func (c *Client) Connect(ctx context.Context, opts ...ConnectOption) (*Conn, error) {
	connectOpts := connectOption{
		model: GPT4oRealtimePreview,
	}
	for _, opt := range opts {
		opt(&connectOpts)
	}

	// default headers
	headers := c.getHeaders()

	// dialOptions
	if connectOpts.dialOptions != nil {
		for k, v := range connectOpts.dialOptions.HTTPHeader {
			headers[k] = append(headers[k], v...)
		}
		connectOpts.dialOptions.HTTPHeader = headers
	} else {
		connectOpts.dialOptions = &websocket.DialOptions{
			HTTPHeader: headers,
		}
	}

	// get url by model
	url := c.getUrl(connectOpts.model)

	// dial
	conn, _, err := websocket.Dial(ctx, url, connectOpts.dialOptions)
	if err != nil {
		return nil, err
	}

	// set readLimit
	if connectOpts.readLimit > 0 {
		conn.SetReadLimit(connectOpts.readLimit)
	}
	return &Conn{conn: conn}, nil
}
