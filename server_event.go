package openairt

import (
	"encoding/json"
	"fmt"
)

type ServerEventType string

const (
	ServerEventTypeError                                            ServerEventType = "error"
	ServerEventTypeSessionCreated                                   ServerEventType = "session.created"
	ServerEventTypeSessionUpdated                                   ServerEventType = "session.updated"
	ServerEventTypeConversationCreated                              ServerEventType = "conversation.created"
	ServerEventTypeInputAudioBufferCommitted                        ServerEventType = "input_audio_buffer.committed"
	ServerEventTypeInputAudioBufferCleared                          ServerEventType = "input_audio_buffer.cleared"
	ServerEventTypeInputAudioBufferSpeechStarted                    ServerEventType = "input_audio_buffer.speech_started"
	ServerEventTypeInputAudioBufferSpeechStopped                    ServerEventType = "input_audio_buffer.speech_stopped"
	ServerEventTypeConversationItemCreated                          ServerEventType = "conversation.item.created"
	ServerEventTypeConversationItemInputAudioTranscriptionCompleted ServerEventType = "conversation.item.input_audio_transcription.completed"
	ServerEventTypeConversationItemInputAudioTranscriptionFailed    ServerEventType = "conversation.item.input_audio_transcription.failed"
	ServerEventTypeConversationItemTruncated                        ServerEventType = "conversation.item.truncated"
	ServerEventTypeConversationItemDeleted                          ServerEventType = "conversation.item.deleted"
	ServerEventTypeResponseCreated                                  ServerEventType = "response.created"
	ServerEventTypeResponseDone                                     ServerEventType = "response.done"
	ServerEventTypeResponseOutputItemAdded                          ServerEventType = "response.output_item.added"
	ServerEventTypeResponseOutputItemDone                           ServerEventType = "response.output_item.done"
	ServerEventTypeResponseContentPartAdded                         ServerEventType = "response.content_part.added"
	ServerEventTypeResponseContentPartDone                          ServerEventType = "response.content_part.done"
	ServerEventTypeResponseTextDelta                                ServerEventType = "response.text.delta"
	ServerEventTypeResponseTextDone                                 ServerEventType = "response.text.done"
	ServerEventTypeResponseAudioTranscriptDelta                     ServerEventType = "response.audio_transcript.delta"
	ServerEventTypeResponseAudioTranscriptDone                      ServerEventType = "response.audio_transcript.done"
	ServerEventTypeResponseAudioDelta                               ServerEventType = "response.audio.delta"
	ServerEventTypeResponseAudioDone                                ServerEventType = "response.audio.done"
	ServerEventTypeResponseFunctionCallArgumentsDelta               ServerEventType = "response.function_call_arguments.delta"
	ServerEventTypeResponseFunctionCallArgumentsDone                ServerEventType = "response.function_call_arguments.done"
	ServerEventTypeRateLimitsUpdated                                ServerEventType = "rate_limits.updated"
)

type ServerEvent interface {
	ServerEventType() ServerEventType
}

type ServerEventBase struct {
	EventID string          `json:"event_id,omitempty"`
	Type    ServerEventType `json:"type"`
}

func (m ServerEventBase) ServerEventType() ServerEventType {
	return m.Type
}

type ErrorEvent struct {
	ServerEventBase
	Error Error `json:"error"`
}

type SessionCreatedEvent struct {
	ServerEventBase
	Session ServerSession `json:"session"`
}

type SessionUpdatedEvent struct {
	ServerEventBase
	Session ServerSession `json:"session"`
}

type ConversationCreatedEvent struct {
	ServerEventBase
	Conversation Conversation `json:"conversation"`
}

type InputAudioBufferCommittedEvent struct {
	ServerEventBase
	PreviousItemID string `json:"previous_item_id,omitempty"`
	ItemID         string `json:"item_id"`
}

type InputAudioBufferClearedEvent struct {
	ServerEventBase
}

type InputAudioBufferSpeechStartedEvent struct {
	ServerEventBase
	AudioStartMs int64  `json:"audio_start_ms"` // Milliseconds since the session started when speech was detected.
	ItemID       string `json:"item_id"`
}

type InputAudioBufferSpeechStoppedEvent struct {
	ServerEventBase
	AudioEndMs int64  `json:"audio_end_ms"` // Milliseconds since the session started when speech stopped.
	ItemID     string `json:"item_id"`
}

type ConversationItemCreatedEvent struct {
	ServerEventBase
	PreviousItemID string              `json:"previous_item_id,omitempty"`
	Item           ResponseMessageItem `json:"item"`
}

type ConversationItemInputAudioTranscriptionCompletedEvent struct {
	ServerEventBase
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	Transcript   string `json:"transcript"`
}

type ConversationItemInputAudioTranscriptionFailedEvent struct {
	ServerEventBase
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	Error        Error  `json:"error"`
}

type ConversationItemTruncatedEvent struct {
	ServerEventBase
	ItemID       string `json:"item_id"`       // The ID of the assistant message item that was truncated.
	ContentIndex int    `json:"content_index"` // The index of the content part that was truncated.
	AudioEndMs   int    `json:"audio_end_ms"`  // The duration up to which the audio was truncated, in milliseconds.
}

type ConversationItemDeletedEvent struct {
	ServerEventBase
	ItemID string `json:"item_id"` // The ID of the item that was deleted.
}

type ResponseCreatedEvent struct {
	ServerEventBase
	Response Response `json:"response"`
}

type ResponseDoneEvent struct {
	ServerEventBase
	Response Response `json:"response"`
}

type ResponseOutputItemAddedEvent struct {
	ServerEventBase
	ResponseID  string              `json:"response_id"`
	OutputIndex int                 `json:"output_index"`
	Item        ResponseMessageItem `json:"item"`
}

type ResponseOutputItemDoneEvent struct {
	ServerEventBase
	ResponseID  string              `json:"response_id"`
	OutputIndex int                 `json:"output_index"`
	Item        ResponseMessageItem `json:"item"`
}

type ResponseContentPartAddedEvent struct {
	ServerEventBase
	ResponseID   string             `json:"response_id"`
	ItemID       string             `json:"item_id"`
	OutputIndex  int                `json:"output_index"`
	ContentIndex int                `json:"content_index"`
	Part         MessageContentPart `json:"part"`
}

type ResponseContentPartDoneEvent struct {
	ServerEventBase
	ResponseID   string             `json:"response_id"`
	ItemID       string             `json:"item_id"`
	OutputIndex  int                `json:"output_index"`
	ContentIndex int                `json:"content_index"`
	Part         MessageContentPart `json:"part"`
}

type ResponseTextDeltaEvent struct {
	ServerEventBase
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

type ResponseTextDoneEvent struct {
	ServerEventBase
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Text         string `json:"text"`
}

type ResponseAudioTranscriptDeltaEvent struct {
	ServerEventBase
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

type ResponseAudioTranscriptDoneEvent struct {
	ServerEventBase
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Transcript   string `json:"transcript"`
}

type ResponseAudioDeltaEvent struct {
	ServerEventBase
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

type ResponseAudioDoneEvent struct {
	ServerEventBase
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
}

type ResponseFunctionCallArgumentsDeltaEvent struct {
	ServerEventBase
	ResponseID  string `json:"response_id"`
	ItemID      string `json:"item_id"`
	OutputIndex int    `json:"output_index"`
	CallID      string `json:"call_id"`
	Delta       string `json:"delta"`
}

type ResponseFunctionCallArgumentsDoneEvent struct {
	ServerEventBase
	ResponseID  string `json:"response_id"`
	ItemID      string `json:"item_id"`
	OutputIndex int    `json:"output_index"`
	CallID      string `json:"call_id"`
	Name        string `json:"name"`
	Arguments   string `json:"arguments"`
}

type RateLimitsUpdatedEvent struct {
	ServerEventBase
	RateLimits []RateLimit `json:"rate_limits"`
}

func UnmarshalServerEvent(data []byte) (ServerEvent, error) {
	var eventType struct {
		Type ServerEventType `json:"type"`
	}
	err := json.Unmarshal(data, &eventType)
	if err != nil {
		return nil, err
	}
	switch eventType.Type {
	case ServerEventTypeError:
		var event ErrorEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeSessionCreated:
		var event SessionCreatedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeSessionUpdated:
		var event SessionUpdatedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeConversationCreated:
		var event ConversationCreatedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeInputAudioBufferCommitted:
		var event InputAudioBufferCommittedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeInputAudioBufferCleared:
		var event InputAudioBufferClearedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeInputAudioBufferSpeechStarted:
		var event InputAudioBufferSpeechStartedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeInputAudioBufferSpeechStopped:
		var event InputAudioBufferSpeechStoppedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeConversationItemCreated:
		var event ConversationItemCreatedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeConversationItemInputAudioTranscriptionCompleted:
		var event ConversationItemInputAudioTranscriptionCompletedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeConversationItemInputAudioTranscriptionFailed:
		var event ConversationItemInputAudioTranscriptionFailedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeConversationItemTruncated:
		var event ConversationItemTruncatedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeConversationItemDeleted:
		var event ConversationItemDeletedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseCreated:
		var event ResponseCreatedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseDone:
		var event ResponseDoneEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseOutputItemAdded:
		var event ResponseOutputItemAddedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseOutputItemDone:
		var event ResponseOutputItemDoneEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseContentPartAdded:
		var event ResponseContentPartAddedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseContentPartDone:
		var event ResponseContentPartDoneEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseTextDelta:
		var event ResponseTextDeltaEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseTextDone:
		var event ResponseTextDoneEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseAudioTranscriptDelta:
		var event ResponseAudioTranscriptDeltaEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseAudioTranscriptDone:
		var event ResponseAudioTranscriptDoneEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseAudioDelta:
		var event ResponseAudioDeltaEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseAudioDone:
		var event ResponseAudioDoneEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseFunctionCallArgumentsDelta:
		var event ResponseFunctionCallArgumentsDeltaEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeResponseFunctionCallArgumentsDone:
		var event ResponseFunctionCallArgumentsDoneEvent
		err = json.Unmarshal(data, &event)
		return event, err
	case ServerEventTypeRateLimitsUpdated:
		var event RateLimitsUpdatedEvent
		err = json.Unmarshal(data, &event)
		return event, err
	default:
		// This should never happen.
		return nil, fmt.Errorf("unknown server event type: %s", eventType.Type)
	}
}
