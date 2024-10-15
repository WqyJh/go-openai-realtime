package openairt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	mockToken := "test"
	client := NewClient(mockToken)
	assert.Equal(t, client.config.authToken, mockToken)

	config := DefaultConfig(mockToken)
	client = NewClientWithConfig(config)
	assert.Equal(t, client.config.authToken, mockToken)
	assert.Equal(t, client.config.BaseURL, openaiAPIURLv1)
	assert.Equal(t, client.config.APIType, APITypeOpenAI)
	url := client.getUrl("test-model")
	assert.Equal(t, url, openaiAPIURLv1+"?model=test-model")
	headers := client.getHeaders()
	assert.Equal(t, headers.Get("Authorization"), "Bearer "+mockToken)
	assert.Equal(t, headers.Get("OpenAI-Beta"), "realtime=v1")

	azureUrl := "wss://my-eastus2-openai-resource.openai.azure.com/openai/realtime"
	config = DefaultAzureConfig(mockToken, azureUrl)
	client = NewClientWithConfig(config)
	assert.Equal(t, client.config.authToken, mockToken)
	assert.Equal(t, client.config.BaseURL, azureUrl)
	assert.Equal(t, client.config.APIType, APITypeAzure)
	assert.Equal(t, client.config.APIVersion, azureAPIVersion20241001Preview)
	url = client.getUrl("test-model")
	assert.Equal(t, url, azureUrl+"?api-version="+azureAPIVersion20241001Preview+"&deployment=test-model")
	headers = client.getHeaders()
	assert.Equal(t, headers.Get("api-key"), mockToken)
}
