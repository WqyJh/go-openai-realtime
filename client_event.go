package openairt

import "encoding/json"

type ClientEventType string

const (
	ClientEventTypeSessionUpdate            ClientEventType = "session.update"
	ClientEventTypeInputAudioBufferAppend   ClientEventType = "input_audio_buffer.append"
	ClientEventTypeInputAudioBufferCommit   ClientEventType = "input_audio_buffer.commit"
	ClientEventTypeInputAudioBufferClear    ClientEventType = "input_audio_buffer.clear"
	ClientEventTypeConversationItemCreate   ClientEventType = "conversation.item.create"
	ClientEventTypeConversationItemTruncate ClientEventType = "conversation.item.truncate"
	ClientEventTypeConversationItemDelete   ClientEventType = "conversation.item.delete"
	ClientEventTypeResponseCreate           ClientEventType = "response.create"
	ClientEventTypeResponseCancel           ClientEventType = "response.cancel"
)

type ClientEvent interface {
	ClientEventType() ClientEventType
}

type EventBase struct {
	EventID string `json:"event_id,omitempty"`
}

type ClientSession struct {
	Modalities              []Modality               `json:"modalities,omitempty"`
	Instructions            string                   `json:"instructions,omitempty"`
	Voice                   Voice                    `json:"voice,omitempty"`
	InputAudioFormat        AudioFormat              `json:"input_audio_format,omitempty"`
	OutputAudioFormat       AudioFormat              `json:"output_audio_format,omitempty"`
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`
	TurnDetection           *TurnDetection           `json:"turn_detection,omitempty"`
	Tools                   []Tool                   `json:"tools,omitempty"`
	ToolChoice              ToolChoiceInterface      `json:"tool_choice,omitempty"`
	Temperature             *float32                 `json:"temperature,omitempty"`
	MaxOutputTokens         int                      `json:"max_output_tokens,omitempty"`
}

type SessionUpdateEvent struct {
	EventBase
	Session ClientSession `json:"session"`
}

func (m SessionUpdateEvent) ClientEventType() ClientEventType {
	return ClientEventTypeSessionUpdate
}

func (m SessionUpdateEvent) MarshalJSON() ([]byte, error) {
	type sessionUpdateEvent SessionUpdateEvent
	v := struct {
		*sessionUpdateEvent
		Type ClientEventType `json:"type"`
	}{
		sessionUpdateEvent: (*sessionUpdateEvent)(&m),
		Type:               m.ClientEventType(),
	}
	return json.Marshal(v)
}

type InputAudioBufferAppendEvent struct {
	EventBase
	Audio string `json:"audio"`
}

func (m InputAudioBufferAppendEvent) ClientEventType() ClientEventType {
	return ClientEventTypeInputAudioBufferAppend
}

func (m InputAudioBufferAppendEvent) MarshalJSON() ([]byte, error) {
	type inputAudioBufferAppendEvent InputAudioBufferAppendEvent
	v := struct {
		*inputAudioBufferAppendEvent
		Type ClientEventType `json:"type"`
	}{
		inputAudioBufferAppendEvent: (*inputAudioBufferAppendEvent)(&m),
		Type:                        m.ClientEventType(),
	}
	return json.Marshal(v)
}

type InputAudioBufferCommitEvent struct {
	EventBase
}

func (m InputAudioBufferCommitEvent) ClientEventType() ClientEventType {
	return ClientEventTypeInputAudioBufferCommit
}

func (m InputAudioBufferCommitEvent) MarshalJSON() ([]byte, error) {
	type inputAudioBufferCommitEvent InputAudioBufferCommitEvent
	v := struct {
		*inputAudioBufferCommitEvent
		Type ClientEventType `json:"type"`
	}{
		inputAudioBufferCommitEvent: (*inputAudioBufferCommitEvent)(&m),
		Type:                        m.ClientEventType(),
	}
	return json.Marshal(v)
}

type InputAudioBufferClearEvent struct {
	EventBase
}

func (m InputAudioBufferClearEvent) ClientEventType() ClientEventType {
	return ClientEventTypeInputAudioBufferClear
}

func (m InputAudioBufferClearEvent) MarshalJSON() ([]byte, error) {
	type inputAudioBufferClearEvent InputAudioBufferClearEvent
	v := struct {
		*inputAudioBufferClearEvent
		Type ClientEventType `json:"type"`
	}{
		inputAudioBufferClearEvent: (*inputAudioBufferClearEvent)(&m),
		Type:                       m.ClientEventType(),
	}
	return json.Marshal(v)
}

type ConversationItemCreateEvent struct {
	EventBase
	PreviousItemID string      `json:"previous_item_id,omitempty"`
	Item           MessageItem `json:"item"`
}

func (m ConversationItemCreateEvent) ClientEventType() ClientEventType {
	return ClientEventTypeConversationItemCreate
}

func (m ConversationItemCreateEvent) MarshalJSON() ([]byte, error) {
	type conversationItemCreateEvent ConversationItemCreateEvent
	v := struct {
		*conversationItemCreateEvent
		Type ClientEventType `json:"type"`
	}{
		conversationItemCreateEvent: (*conversationItemCreateEvent)(&m),
		Type:                        m.ClientEventType(),
	}
	return json.Marshal(v)
}

type ConversationItemTruncateEvent struct {
	EventBase
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	AudioEndMs   int    `json:"audio_end_ms"`
}

func (m ConversationItemTruncateEvent) ClientEventType() ClientEventType {
	return ClientEventTypeConversationItemTruncate
}

func (m ConversationItemTruncateEvent) MarshalJSON() ([]byte, error) {
	type conversationItemTruncateEvent ConversationItemTruncateEvent
	v := struct {
		*conversationItemTruncateEvent
		Type ClientEventType `json:"type"`
	}{
		conversationItemTruncateEvent: (*conversationItemTruncateEvent)(&m),
		Type:                          m.ClientEventType(),
	}
	return json.Marshal(v)
}

type ConversationItemDeleteEvent struct {
	EventBase
	ItemID string `json:"item_id"`
}

func (m ConversationItemDeleteEvent) ClientEventType() ClientEventType {
	return ClientEventTypeConversationItemDelete
}

func (m ConversationItemDeleteEvent) MarshalJSON() ([]byte, error) {
	type conversationItemDeleteEvent ConversationItemDeleteEvent
	v := struct {
		*conversationItemDeleteEvent
		Type ClientEventType `json:"type"`
	}{
		conversationItemDeleteEvent: (*conversationItemDeleteEvent)(&m),
		Type:                        m.ClientEventType(),
	}
	return json.Marshal(v)
}

type ResponseCreateParams struct {
	Modalities        []Modality          `json:"modalities,omitempty"`
	Instructions      string              `json:"instructions,omitempty"`
	Voice             Voice               `json:"voice,omitempty"`
	OutputAudioFormat AudioFormat         `json:"output_audio_format,omitempty"`
	Tools             []Tool              `json:"tools,omitempty"`
	ToolChoice        ToolChoiceInterface `json:"tool_choice,omitempty"`
	Temperature       *float32            `json:"temperature,omitempty"`
	MaxOutputTokens   int                 `json:"max_output_tokens,omitempty"`
}

type ResponseCreateEvent struct {
	EventBase
	Response ResponseCreateParams `json:"response"`
}

func (m ResponseCreateEvent) ClientEventType() ClientEventType {
	return ClientEventTypeResponseCreate
}

func (m ResponseCreateEvent) MarshalJSON() ([]byte, error) {
	type responseCreateEvent ResponseCreateEvent
	v := struct {
		*responseCreateEvent
		Type ClientEventType `json:"type"`
	}{
		responseCreateEvent: (*responseCreateEvent)(&m),
		Type:                m.ClientEventType(),
	}
	return json.Marshal(v)
}

type ResponseCancelEvent struct {
	EventBase
}

func (m ResponseCancelEvent) ClientEventType() ClientEventType {
	return ClientEventTypeResponseCancel
}

func (m ResponseCancelEvent) MarshalJSON() ([]byte, error) {
	type responseCancelEvent ResponseCancelEvent
	v := struct {
		*responseCancelEvent
		Type ClientEventType `json:"type"`
	}{
		responseCancelEvent: (*responseCancelEvent)(&m),
		Type:                m.ClientEventType(),
	}
	return json.Marshal(v)
}

func MarshalClientEvent(event ClientEvent) ([]byte, error) {
	return json.Marshal(event)
}
