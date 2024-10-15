package openairt

type APIType string

const (
	APITypeOpenAI APIType = "OPEN_AI"
	APITypeAzure  APIType = "AZURE"
)

const (
	openaiAPIURLv1 = "wss://api.openai.com/v1/realtime"
)

const (
	azureAPIVersion20241001Preview = "2024-10-01-preview"
)

type ClientConfig struct {
	authToken string

	BaseURL    string
	APIType    APIType
	APIVersion string // required when APIType is APITypeAzure
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,
		BaseURL:   openaiAPIURLv1,
		APIType:   APITypeOpenAI,
	}
}

func DefaultAzureConfig(apiKey, baseURL string) ClientConfig {
	return ClientConfig{
		authToken:  apiKey,
		BaseURL:    baseURL,
		APIType:    APITypeAzure,
		APIVersion: azureAPIVersion20241001Preview,
	}
}

func (ClientConfig) String() string {
	return "<OpenAI Realtime API ClientConfig>"
}
