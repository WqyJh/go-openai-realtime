package openairt //nolint:testpackage // Need to access unexported fields

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	mockToken := "test"
	client := NewClient(mockToken)
	require.Equal(t, client.config.authToken, mockToken)

	config := DefaultConfig(mockToken)
	client = NewClientWithConfig(config)
	require.Equal(t, mockToken, client.config.authToken)
	require.Equal(t, openaiAPIURLv1, client.config.BaseURL)
	require.Equal(t, APITypeOpenAI, client.config.APIType)
	url := client.getURL("test-model")
	require.Equal(t, openaiAPIURLv1+"?model=test-model", url)
	headers := client.getHeaders()
	require.Equal(t, "Bearer "+mockToken, headers.Get("Authorization"))
	require.Equal(t, "realtime=v1", headers.Get("OpenAI-Beta"))

	azureURL := "wss://my-eastus2-openai-resource.openai.azure.com/openai/realtime"
	config = DefaultAzureConfig(mockToken, azureURL)
	client = NewClientWithConfig(config)
	require.Equal(t, mockToken, client.config.authToken)
	require.Equal(t, azureURL, client.config.BaseURL)
	require.Equal(t, APITypeAzure, client.config.APIType)
	require.Equal(t, azureAPIVersion20241001Preview, client.config.APIVersion)
	url = client.getURL("test-model")
	require.Equal(t, azureURL+"?api-version="+azureAPIVersion20241001Preview+"&deployment=test-model", url)
	headers = client.getHeaders()
	require.Equal(t, mockToken, headers.Get("api-key"))
}
