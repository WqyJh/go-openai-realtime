package openairt

// APIType is the type of API.
type APIType string

const (
	// APITypeOpenAI is the type of API for OpenAI.
	APITypeOpenAI APIType = "OPEN_AI"
	// APITypeAzure is the type of API for Azure.
	APITypeAzure APIType = "AZURE"
)

const (
	// OpenaiRealtimeAPIURLv1 is the base URL for the OpenAI Realtime API.
	OpenaiRealtimeAPIURLv1 = "wss://api.openai.com/v1/realtime"
)

const (
	// azureAPIVersion20241001Preview is the API version for Azure.
	azureAPIVersion20241001Preview = "2024-10-01-preview"
)

// ClientConfig is the configuration for the client.
type ClientConfig struct {
	authToken string

	BaseURL    string  // Base URL for the API. Defaults to "wss://api.openai.com/v1/realtime"
	APIType    APIType // API type. Defaults to APITypeOpenAI
	APIVersion string  // required when APIType is APITypeAzure
}

// DefaultConfig creates a new ClientConfig with the given auth token.
// Defaults to using the OpenAI Realtime API.
func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,
		BaseURL:   OpenaiRealtimeAPIURLv1,
		APIType:   APITypeOpenAI,
	}
}

// DefaultAzureConfig creates a new ClientConfig with the given auth token and base URL.
// Defaults to using the Azure Realtime API.
func DefaultAzureConfig(apiKey, baseURL string) ClientConfig {
	return ClientConfig{
		authToken:  apiKey,
		BaseURL:    baseURL,
		APIType:    APITypeAzure,
		APIVersion: azureAPIVersion20241001Preview,
	}
}

// String returns a string representation of the ClientConfig.
func (c ClientConfig) String() string {
	return "<OpenAI Realtime API ClientConfig>"
}
