package openairt

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai/jsonschema"
)

type Voice string

const (
	VoiceAlloy   Voice = "alloy"
	VoiceShimmer Voice = "shimmer"
	VoiceEcho    Voice = "echo"
)

type AudioFormat string

const (
	AudioFormatPcm16    AudioFormat = "pcm16"
	AudioFormatG711Ulaw AudioFormat = "g711-ulaw"
	AudioFormatG711Alaw AudioFormat = "g711-alaw"
)

type Modality string

const (
	ModalityText  Modality = "text"
	ModalityAudio Modality = "audio"
)

type TurnDetectionType string

const (
	TurnDetectionTypeNone      TurnDetectionType = "none"
	TurnDetectionTypeServerVad TurnDetectionType = "server_vad"
)

type TurnDetection struct {
	Type              TurnDetectionType `json:"type"`
	Threshold         float64           `json:"threshold,omitempty"`
	PrefixPaddingMs   int               `json:"prefix_padding_ms,omitempty"`
	SilenceDurationMs int               `json:"silence_duration_ms,omitempty"`
}

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type ToolChoiceInterface interface {
	ToolChoice()
}

type ToolChoiceString string

func (ToolChoiceString) ToolChoice() {}

const (
	ToolChoiceAuto     ToolChoiceString = "auto"
	ToolChoiceNone     ToolChoiceString = "none"
	ToolChoiceRequired ToolChoiceString = "required"
)

type ToolChoice struct {
	Type     ToolType     `json:"type"`
	Function ToolFunction `json:"function,omitempty"`
}

func (t ToolChoice) ToolChoice() {}

type ToolFunction struct {
	Name string `json:"name"`
}

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleUser      MessageRole = "user"
)

type InputAudioTranscription struct {
	Enabled bool   `json:"enabled"`
	Model   string `json:"model"`
}

type Tool struct {
	Type        ToolType              `json:"type"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  jsonschema.Definition `json:"parameters"`
}

type MessageItemType string

const (
	MessageItemTypeMessage            MessageItemType = "message"
	MessageItemTypeFunctionCall       MessageItemType = "function_call"
	MessageItemTypeFunctionCallOutput MessageItemType = "function_call_output"
)

type MessageContentType string

const (
	MessageContentTypeText       MessageContentType = "text"
	MessageContentTypeAudio      MessageContentType = "audio"
	MessageContentTypeTranscript MessageContentType = "transcript"
	MessageContentTypeInputText  MessageContentType = "input_text"
	MessageContentTypeInputAudio MessageContentType = "input_audio"
)

type MessageContentPart struct {
	Type       MessageContentType `json:"type"`
	Text       string             `json:"text,omitempty"`
	Audio      string             `json:"audio,omitempty"`
	Transcript string             `json:"transcript,omitempty"`
}

type MessageItem struct {
	ID      string               `json:"id"`
	Type    MessageItemType      `json:"type"`
	Status  ItemStatus           `json:"status"`
	Role    MessageRole          `json:"role"`
	Content []MessageContentPart `json:"content"`
}

type ResponseMessageItem struct {
	MessageItem
	Object string `json:"object,omitempty"`
}

type Error struct {
	Message string `json:"message,omitempty"`
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
	Param   string `json:"param,omitempty"`
	EventID string `json:"event_id,omitempty"`
}

// ServerToolChoice is a type that can be used to choose a tool response from the server.
type ServerToolChoice struct {
	String   ToolChoiceString
	Function ToolChoice
}

// UnmarshalJSON is a custom unmarshaler for ServerToolChoice.
func (m *ServerToolChoice) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &m.Function)
	if err != nil {
		if data[0] == '"' {
			data = data[1:]
		}
		if data[len(data)-1] == '"' {
			data = data[:len(data)-1]
		}
		m.String = ToolChoiceString(data)
		m.Function = ToolChoice{}
		return nil
	}
	return nil
}

// IsFunction returns true if the tool choice is a function call.
func (m *ServerToolChoice) IsFunction() bool {
	return m.Function.Type == ToolTypeFunction
}

// Get returns the ToolChoiceInterface based on the type of tool choice.
func (m ServerToolChoice) Get() ToolChoiceInterface {
	if m.IsFunction() {
		return m.Function
	}
	return m.String
}

type ServerSession struct {
	ID                      string                   `json:"id"`
	Object                  string                   `json:"object"`
	Model                   string                   `json:"model"`
	Modalities              []Modality               `json:"modalities,omitempty"`
	Instructions            string                   `json:"instructions,omitempty"`
	Voice                   Voice                    `json:"voice,omitempty"`
	InputAudioFormat        AudioFormat              `json:"input_audio_format,omitempty"`
	OutputAudioFormat       AudioFormat              `json:"output_audio_format,omitempty"`
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`
	TurnDetection           *TurnDetection           `json:"turn_detection,omitempty"`
	Tools                   []Tool                   `json:"tools,omitempty"`
	ToolChoice              ServerToolChoice         `json:"tool_choice,omitempty"`
	Temperature             *float32                 `json:"temperature,omitempty"`
	MaxOutputTokens         int                      `json:"max_output_tokens,omitempty"`
}

type ItemStatus string

const (
	ItemStatusInProgress ItemStatus = "in_progress"
	ItemStatusCompleted  ItemStatus = "completed"
	ItemStatusIncomplete ItemStatus = "incomplete"
)

type Conversation struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

type ResponseStatus string

const (
	ResponseStatusInProgress ResponseStatus = "in_progress"
	ResponseStatusCompleted  ResponseStatus = "completed"
	ResponseStatusCancelled  ResponseStatus = "cancelled"
	ResponseStatusIncomplete ResponseStatus = "incomplete"
	ResponseStatusFailed     ResponseStatus = "failed"
)

type Usage struct {
	TotalTokens  int `json:"total_tokens"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type Response struct {
	ID            string                `json:"id"`
	Object        string                `json:"object"`
	Status        ResponseStatus        `json:"status"`
	StatusDetails any                   `json:"status_details,omitempty"`
	Output        []ResponseMessageItem `json:"output"`
	Usage         *Usage                `json:"usage,omitempty"`
}

type RateLimit struct {
	Name         string  `json:"name"`
	Limit        int     `json:"limit"`
	Remaining    int     `json:"remaining"`
	ResetSeconds float64 `json:"reset_seconds"`
}
