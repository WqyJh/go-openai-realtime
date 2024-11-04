package openairt_test

import (
	"encoding/json"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/require"
)

func TestErrorEvent(t *testing.T) {
	data := `{
    "event_id": "event_890",
    "type": "error",
    "error": {
        "type": "invalid_request_error",
        "code": "invalid_event",
        "message": "The 'type' field is missing.",
        "param": null,
        "event_id": "event_567"
    }
}`
	expected := openairt.ErrorEvent{
		ServerEventBase: openairt.ServerEventBase{
			Type:    openairt.ServerEventTypeError,
			EventID: "event_890",
		},
		Error: openairt.Error{
			Type:    "invalid_request_error",
			Code:    "invalid_event",
			Message: "The 'type' field is missing.",
			Param:   "",
			EventID: "event_567",
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeError, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ErrorEvent))
}

func TestServerToolChoice(t *testing.T) {
	data := `{"type": "function", "function": {"name": "get_current_weather"}}`
	expectedFunction := openairt.ToolChoice{
		Type: openairt.ToolTypeFunction,
		Function: openairt.ToolFunction{
			Name: "get_current_weather",
		},
	}
	expected := openairt.ServerToolChoice{Function: expectedFunction}
	actual := openairt.ServerToolChoice{}
	err := json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, expectedFunction, actual.Get())

	data = `"auto"`
	expected = openairt.ServerToolChoice{String: openairt.ToolChoiceAuto}
	actual = openairt.ServerToolChoice{}
	err = json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, openairt.ToolChoiceAuto, actual.Get())

	data = `"none"`
	expected = openairt.ServerToolChoice{String: openairt.ToolChoiceNone}
	actual = openairt.ServerToolChoice{}
	err = json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, openairt.ToolChoiceNone, actual.Get())

	data = `"required"`
	expected = openairt.ServerToolChoice{String: openairt.ToolChoiceRequired}
	actual = openairt.ServerToolChoice{}
	err = json.Unmarshal([]byte(data), &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, openairt.ToolChoiceRequired, actual.Get())
}

func TestSessionCreatedEvent(t *testing.T) {
	data := `{
    "event_id": "event_1234",
    "type": "session.created",
    "session": {
        "id": "sess_001",
        "object": "realtime.session",
        "model": "gpt-4o-realtime-preview-2024-10-01",
        "modalities": ["text", "audio"],
        "instructions": "",
        "voice": "alloy",
        "input_audio_format": "pcm16",
        "output_audio_format": "pcm16",
        "input_audio_transcription": null,
        "turn_detection": {
            "type": "server_vad",
            "threshold": 0.5,
            "prefix_padding_ms": 300,
            "silence_duration_ms": 200
        },
        "tools": [],
        "tool_choice": "auto",
        "temperature": 0.8,
        "max_response_output_tokens": null
    }
}`
	temperature := float32(0.8)
	expected := openairt.SessionCreatedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_1234",
			Type:    openairt.ServerEventTypeSessionCreated,
		},
		Session: openairt.ServerSession{
			ID:                      "sess_001",
			Object:                  "realtime.session",
			Model:                   openairt.GPT4oRealtimePreview20241001,
			Modalities:              []openairt.Modality{openairt.ModalityText, openairt.ModalityAudio},
			Instructions:            "",
			Voice:                   openairt.VoiceAlloy,
			InputAudioFormat:        openairt.AudioFormatPcm16,
			OutputAudioFormat:       openairt.AudioFormatPcm16,
			InputAudioTranscription: nil,
			TurnDetection: &openairt.ServerTurnDetection{
				Type: openairt.ServerTurnDetectionTypeServerVad,
				TurnDetectionParams: openairt.TurnDetectionParams{
					Threshold:         0.5,
					PrefixPaddingMs:   300,
					SilenceDurationMs: 200,
				},
			},
			Tools:           []openairt.Tool{},
			ToolChoice:      openairt.ServerToolChoice{String: openairt.ToolChoiceAuto},
			Temperature:     &temperature,
			MaxOutputTokens: 0,
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeSessionCreated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.SessionCreatedEvent))
}

func TestSessionUpdatedEvent(t *testing.T) {
	data := `{
    "event_id": "event_5678",
    "type": "session.updated",
    "session": {
        "id": "sess_001",
        "object": "realtime.session",
        "model": "gpt-4o-realtime-preview-2024-10-01",
        "modalities": ["text"],
        "instructions": "New instructions",
        "voice": "alloy",
        "input_audio_format": "pcm16",
        "output_audio_format": "pcm16",
        "input_audio_transcription": {
            "model": "whisper-1"
        },
        "turn_detection": {
            "type": "none"
        },
        "tools": [],
        "tool_choice": "none",
        "temperature": 0.7,
        "max_response_output_tokens": 200
    }
}`
	temperature := float32(0.7)
	expected := openairt.SessionUpdatedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_5678",
			Type:    openairt.ServerEventTypeSessionUpdated,
		},
		Session: openairt.ServerSession{
			ID:                "sess_001",
			Object:            "realtime.session",
			Model:             openairt.GPT4oRealtimePreview20241001,
			Modalities:        []openairt.Modality{openairt.ModalityText},
			Instructions:      "New instructions",
			Voice:             openairt.VoiceAlloy,
			InputAudioFormat:  openairt.AudioFormatPcm16,
			OutputAudioFormat: openairt.AudioFormatPcm16,
			InputAudioTranscription: &openairt.InputAudioTranscription{
				Model: "whisper-1",
			},
			TurnDetection: &openairt.ServerTurnDetection{
				Type: openairt.ServerTurnDetectionTypeNone,
			},
			Tools:           []openairt.Tool{},
			ToolChoice:      openairt.ServerToolChoice{String: openairt.ToolChoiceNone},
			Temperature:     &temperature,
			MaxOutputTokens: 200,
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeSessionUpdated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.SessionUpdatedEvent))

	data = `{
		"event_id": "event_5678",
		"type": "session.updated",
		"session": {
			"id": "sess_001",
			"object": "realtime.session",
			"model": "gpt-4o-realtime-preview-2024-10-01",
			"modalities": ["text"],
			"instructions": "New instructions",
			"voice": "alloy",
			"input_audio_format": "pcm16",
			"output_audio_format": "pcm16",
			"input_audio_transcription": {
				"model": "whisper-1"
			},
			"turn_detection": {
				"type": "none"
			},
			"tools": [],
			"tool_choice": "none",
			"temperature": 0.7,
			"max_response_output_tokens": "inf"
		}
	}`
	expected.Session.MaxOutputTokens = openairt.Inf
	actual, err = openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeSessionUpdated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.SessionUpdatedEvent))
}

func TestConversationCreatedEvent(t *testing.T) {
	data := `{
    "event_id": "event_9101",
    "type": "conversation.created",
    "conversation": {
        "id": "conv_001",
        "object": "realtime.conversation"
    }
}`
	expected := openairt.ConversationCreatedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_9101",
			Type:    openairt.ServerEventTypeConversationCreated,
		},
		Conversation: openairt.Conversation{
			ID:     "conv_001",
			Object: "realtime.conversation",
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeConversationCreated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ConversationCreatedEvent))
}

func TestConversationItemCreatedEvent(t *testing.T) {
	data := `{
    "event_id": "event_1920",
    "type": "conversation.item.created",
    "previous_item_id": "msg_002",
    "item": {
        "id": "msg_003",
        "object": "realtime.item",
        "type": "message",
        "status": "completed",
        "role": "user",
        "content": [
            {
                "type": "input_audio",
                "transcript": null
            }
        ]
    }
}`
	expected := openairt.ConversationItemCreatedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_1920",
			Type:    openairt.ServerEventTypeConversationItemCreated,
		},
		PreviousItemID: "msg_002",
		Item: openairt.ResponseMessageItem{
			Object: "realtime.item",
			MessageItem: openairt.MessageItem{
				ID:     "msg_003",
				Type:   openairt.MessageItemTypeMessage,
				Status: openairt.ItemStatusCompleted,
				Role:   openairt.MessageRoleUser,
				Content: []openairt.MessageContentPart{
					{
						Type:       openairt.MessageContentTypeInputAudio,
						Transcript: "",
					},
				},
			},
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeConversationItemCreated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ConversationItemCreatedEvent))
}

func TestConversationItemInputAudioTranscriptionCompletedEvent(t *testing.T) {
	data := `{
    "event_id": "event_2122",
    "type": "conversation.item.input_audio_transcription.completed",
    "item_id": "msg_003",
    "content_index": 1,
    "transcript": "Hello, how are you?"
}`
	expected := openairt.ConversationItemInputAudioTranscriptionCompletedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_2122",
			Type:    openairt.ServerEventTypeConversationItemInputAudioTranscriptionCompleted,
		},
		ItemID:       "msg_003",
		ContentIndex: 1,
		Transcript:   "Hello, how are you?",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeConversationItemInputAudioTranscriptionCompleted, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ConversationItemInputAudioTranscriptionCompletedEvent))
}

func TestConversationItemInputAudioTranscriptionFailedEvent(t *testing.T) {
	data := `{
    "event_id": "event_2324",
    "type": "conversation.item.input_audio_transcription.failed",
    "item_id": "msg_003",
    "content_index": 0,
    "error": {
        "type": "transcription_error",
        "code": "audio_unintelligible",
        "message": "The audio could not be transcribed.",
        "param": null
    }
}`
	expected := openairt.ConversationItemInputAudioTranscriptionFailedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_2324",
			Type:    openairt.ServerEventTypeConversationItemInputAudioTranscriptionFailed,
		},
		ItemID:       "msg_003",
		ContentIndex: 0,
		Error: openairt.Error{
			Type:    "transcription_error",
			Code:    "audio_unintelligible",
			Message: "The audio could not be transcribed.",
			Param:   "",
			EventID: "",
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeConversationItemInputAudioTranscriptionFailed, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ConversationItemInputAudioTranscriptionFailedEvent))
}

func TestConversationItemTruncatedEvent(t *testing.T) {
	data := `{
    "event_id": "event_2526",
    "type": "conversation.item.truncated",
    "item_id": "msg_004",
    "content_index": 2,
    "audio_end_ms": 1500
}`
	expected := openairt.ConversationItemTruncatedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_2526",
			Type:    openairt.ServerEventTypeConversationItemTruncated,
		},
		ItemID:       "msg_004",
		ContentIndex: 2,
		AudioEndMs:   1500,
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeConversationItemTruncated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ConversationItemTruncatedEvent))
}

func TestConversationItemDeletedEvent(t *testing.T) {
	data := `{
    "event_id": "event_2728",
    "type": "conversation.item.deleted",
    "item_id": "msg_005"
}`
	expected := openairt.ConversationItemDeletedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_2728",
			Type:    openairt.ServerEventTypeConversationItemDeleted,
		},
		ItemID: "msg_005",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeConversationItemDeleted, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ConversationItemDeletedEvent))
}

func TestInputAudioBufferCommittedEvent(t *testing.T) {
	data := `{
    "event_id": "event_1121",
    "type": "input_audio_buffer.committed",
    "previous_item_id": "msg_001",
    "item_id": "msg_002"
}`
	expected := openairt.InputAudioBufferCommittedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_1121",
			Type:    openairt.ServerEventTypeInputAudioBufferCommitted,
		},
		PreviousItemID: "msg_001",
		ItemID:         "msg_002",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeInputAudioBufferCommitted, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.InputAudioBufferCommittedEvent))
}

func TestInputAudioBufferClearedEvent(t *testing.T) {
	data := `{
    "event_id": "event_1314",
    "type": "input_audio_buffer.cleared"
}`
	expected := openairt.InputAudioBufferClearedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_1314",
			Type:    openairt.ServerEventTypeInputAudioBufferCleared,
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeInputAudioBufferCleared, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.InputAudioBufferClearedEvent))
}

func TestInputAudioBufferSpeechStartedEvent(t *testing.T) {
	data := `{
    "event_id": "event_1516",
    "type": "input_audio_buffer.speech_started",
    "audio_start_ms": 1000,
    "item_id": "msg_003"
}`
	expected := openairt.InputAudioBufferSpeechStartedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_1516",
			Type:    openairt.ServerEventTypeInputAudioBufferSpeechStarted,
		},
		AudioStartMs: 1000,
		ItemID:       "msg_003",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeInputAudioBufferSpeechStarted, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.InputAudioBufferSpeechStartedEvent))
}

func TestInputAudioBufferSpeechStoppedEvent(t *testing.T) {
	data := `{
    "event_id": "event_1718",
    "type": "input_audio_buffer.speech_stopped",
    "audio_end_ms": 2000,
    "item_id": "msg_003"
}`
	expected := openairt.InputAudioBufferSpeechStoppedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_1718",
			Type:    openairt.ServerEventTypeInputAudioBufferSpeechStopped,
		},
		AudioEndMs: 2000,
		ItemID:     "msg_003",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeInputAudioBufferSpeechStopped, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.InputAudioBufferSpeechStoppedEvent))
}

func TestResponseCreatedEvent(t *testing.T) {
	data := `{
    "event_id": "event_2930",
    "type": "response.created",
    "response": {
        "id": "resp_001",
        "object": "realtime.response",
        "status": "in_progress",
        "status_details": null,
        "output": [],
        "usage": null
    }
}`
	expected := openairt.ResponseCreatedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_2930",
			Type:    openairt.ServerEventTypeResponseCreated,
		},
		Response: openairt.Response{
			ID:            "resp_001",
			Object:        "realtime.response",
			Status:        "in_progress",
			StatusDetails: nil,
			Output:        []openairt.ResponseMessageItem{},
			Usage:         nil,
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseCreated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseCreatedEvent))
}

func TestResponseDoneEvent(t *testing.T) {
	data := `{
    "event_id": "event_3132",
    "type": "response.done",
    "response": {
        "id": "resp_001",
        "object": "realtime.response",
        "status": "completed",
        "status_details": null,
        "output": [
            {
                "id": "msg_006",
                "object": "realtime.item",
                "type": "message",
                "status": "completed",
                "role": "assistant",
                "content": [
                    {
                        "type": "text",
                        "text": "Sure, how can I assist you today?"
                    }
                ]
            }
        ],
        "usage": {
            "total_tokens": 50,
            "input_tokens": 20,
            "output_tokens": 30
        }
    }
}`
	expected := openairt.ResponseDoneEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_3132",
			Type:    openairt.ServerEventTypeResponseDone,
		},
		Response: openairt.Response{
			ID:            "resp_001",
			Object:        "realtime.response",
			Status:        openairt.ResponseStatusCompleted,
			StatusDetails: nil,
			Output: []openairt.ResponseMessageItem{
				{
					Object: "realtime.item",
					MessageItem: openairt.MessageItem{
						ID:     "msg_006",
						Type:   openairt.MessageItemTypeMessage,
						Status: openairt.ItemStatusCompleted,
						Role:   openairt.MessageRoleAssistant,
						Content: []openairt.MessageContentPart{
							{
								Type: openairt.MessageContentTypeText,
								Text: "Sure, how can I assist you today?",
							},
						},
					},
				},
			},
			Usage: &openairt.Usage{
				TotalTokens:  50,
				InputTokens:  20,
				OutputTokens: 30,
			},
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseDone, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseDoneEvent))
}

func TestResponseOutputItemAddedEvent(t *testing.T) {
	data := `{
    "event_id": "event_3334",
    "type": "response.output_item.added",
    "response_id": "resp_001",
    "output_index": 10,
    "item": {
        "id": "msg_007",
        "object": "realtime.item",
        "type": "message",
        "status": "in_progress",
        "role": "assistant",
        "content": []
    }
}`
	expected := openairt.ResponseOutputItemAddedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_3334",
			Type:    openairt.ServerEventTypeResponseOutputItemAdded,
		},
		ResponseID:  "resp_001",
		OutputIndex: 10,
		Item: openairt.ResponseMessageItem{
			Object: "realtime.item",
			MessageItem: openairt.MessageItem{
				ID:      "msg_007",
				Type:    openairt.MessageItemTypeMessage,
				Status:  openairt.ItemStatusInProgress,
				Role:    openairt.MessageRoleAssistant,
				Content: []openairt.MessageContentPart{},
			},
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseOutputItemAdded, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseOutputItemAddedEvent))
}

func TestResponseOutputItemDoneEvent(t *testing.T) {
	data := `{
    "event_id": "event_3536",
    "type": "response.output_item.done",
    "response_id": "resp_001",
    "output_index": 0,
    "item": {
        "id": "msg_007",
        "object": "realtime.item",
        "type": "message",
        "status": "completed",
        "role": "assistant",
        "content": [
            {
                "type": "text",
                "text": "Sure, I can help with that."
            }
        ]
    }
}`
	expected := openairt.ResponseOutputItemDoneEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_3536",
			Type:    openairt.ServerEventTypeResponseOutputItemDone,
		},
		ResponseID:  "resp_001",
		OutputIndex: 0,
		Item: openairt.ResponseMessageItem{
			Object: "realtime.item",
			MessageItem: openairt.MessageItem{
				ID:     "msg_007",
				Type:   openairt.MessageItemTypeMessage,
				Status: openairt.ItemStatusCompleted,
				Role:   openairt.MessageRoleAssistant,
				Content: []openairt.MessageContentPart{
					{
						Type: openairt.MessageContentTypeText,
						Text: "Sure, I can help with that.",
					},
				},
			},
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseOutputItemDone, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseOutputItemDoneEvent))
}

func TestResponseContentPartAddedEvent(t *testing.T) {
	data := `{
    "event_id": "event_3738",
    "type": "response.content_part.added",
    "response_id": "resp_001",
    "item_id": "msg_007",
    "output_index": 1,
    "content_index": 2,
    "part": {
        "type": "text",
        "text": ""
    }
}`
	expected := openairt.ResponseContentPartAddedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_3738",
			Type:    openairt.ServerEventTypeResponseContentPartAdded,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_007",
		OutputIndex:  1,
		ContentIndex: 2,
		Part: openairt.MessageContentPart{
			Type: openairt.MessageContentTypeText,
			Text: "",
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseContentPartAdded, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseContentPartAddedEvent))
}

func TestResponseContentPartDoneEvent(t *testing.T) {
	data := `{
    "event_id": "event_3940",
    "type": "response.content_part.done",
    "response_id": "resp_001",
    "item_id": "msg_007",
    "output_index": 0,
    "content_index": 0,
    "part": {
        "type": "text",
        "text": "Sure, I can help with that."
    }
}`
	expected := openairt.ResponseContentPartDoneEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_3940",
			Type:    openairt.ServerEventTypeResponseContentPartDone,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_007",
		OutputIndex:  0,
		ContentIndex: 0,
		Part: openairt.MessageContentPart{
			Type: openairt.MessageContentTypeText,
			Text: "Sure, I can help with that.",
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseContentPartDone, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseContentPartDoneEvent))
}

func TestResponseTextDelta(t *testing.T) {
	data := `{
    "event_id": "event_4142",
    "type": "response.text.delta",
    "response_id": "resp_001",
    "item_id": "msg_007",
    "output_index": 0,
    "content_index": 0,
    "delta": "Sure, I can h"
}`
	expected := openairt.ResponseTextDeltaEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_4142",
			Type:    openairt.ServerEventTypeResponseTextDelta,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_007",
		OutputIndex:  0,
		ContentIndex: 0,
		Delta:        "Sure, I can h",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseTextDelta, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseTextDeltaEvent))
}

func TestResponseTextDoneEvent(t *testing.T) {
	data := `{
    "event_id": "event_4344",
    "type": "response.text.done",
    "response_id": "resp_001",
    "item_id": "msg_007",
    "output_index": 0,
    "content_index": 0,
    "text": "Sure, I can help with that."
}`
	expected := openairt.ResponseTextDoneEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_4344",
			Type:    openairt.ServerEventTypeResponseTextDone,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_007",
		OutputIndex:  0,
		ContentIndex: 0,
		Text:         "Sure, I can help with that.",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseTextDone, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseTextDoneEvent))
}

func TestResponseAudioTranscriptDelta(t *testing.T) {
	data := `{
    "event_id": "event_4546",
    "type": "response.audio_transcript.delta",
    "response_id": "resp_001",
    "item_id": "msg_008",
    "output_index": 1,
    "content_index": 2,
    "delta": "Hello, how can I a"
}`
	expected := openairt.ResponseAudioTranscriptDeltaEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_4546",
			Type:    openairt.ServerEventTypeResponseAudioTranscriptDelta,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_008",
		OutputIndex:  1,
		ContentIndex: 2,
		Delta:        "Hello, how can I a",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseAudioTranscriptDelta, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseAudioTranscriptDeltaEvent))
}

func TestResponseAudioTranscriptDoneEvent(t *testing.T) {
	data := `{
    "event_id": "event_4748",
    "type": "response.audio_transcript.done",
    "response_id": "resp_001",
    "item_id": "msg_008",
    "output_index": 0,
    "content_index": 0,
    "transcript": "Hello, how can I assist you today?"
}`
	expected := openairt.ResponseAudioTranscriptDoneEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_4748",
			Type:    openairt.ServerEventTypeResponseAudioTranscriptDone,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_008",
		OutputIndex:  0,
		ContentIndex: 0,
		Transcript:   "Hello, how can I assist you today?",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseAudioTranscriptDone, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseAudioTranscriptDoneEvent))
}

func TestResponseAudioDeltaEvent(t *testing.T) {
	data := `{
    "event_id": "event_4950",
    "type": "response.audio.delta",
    "response_id": "resp_001",
    "item_id": "msg_008",
    "output_index": 0,
    "content_index": 0,
    "delta": "Base64EncodedAudioDelta"
}`
	expected := openairt.ResponseAudioDeltaEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_4950",
			Type:    openairt.ServerEventTypeResponseAudioDelta,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_008",
		OutputIndex:  0,
		ContentIndex: 0,
		Delta:        "Base64EncodedAudioDelta",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseAudioDelta, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseAudioDeltaEvent))
}

func TestResponseAudioDoneEvent(t *testing.T) {
	data := `{
    "event_id": "event_5152",
    "type": "response.audio.done",
    "response_id": "resp_001",
    "item_id": "msg_008",
    "output_index": 0,
    "content_index": 0
}`
	expected := openairt.ResponseAudioDoneEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_5152",
			Type:    openairt.ServerEventTypeResponseAudioDone,
		},
		ResponseID:   "resp_001",
		ItemID:       "msg_008",
		OutputIndex:  0,
		ContentIndex: 0,
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseAudioDone, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseAudioDoneEvent))
}

func TestResponseFunctionCallArgumentsDelta(t *testing.T) {
	data := `{
    "event_id": "event_5354",
    "type": "response.function_call_arguments.delta",
    "response_id": "resp_002",
    "item_id": "fc_001",
    "output_index": 0,
    "call_id": "call_001",
    "delta": "{\"location\": \"San\""
}`
	expected := openairt.ResponseFunctionCallArgumentsDeltaEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_5354",
			Type:    openairt.ServerEventTypeResponseFunctionCallArgumentsDelta,
		},
		ResponseID:  "resp_002",
		ItemID:      "fc_001",
		OutputIndex: 0,
		CallID:      "call_001",
		Delta:       "{\"location\": \"San\"",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseFunctionCallArgumentsDelta, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseFunctionCallArgumentsDeltaEvent))
}

func TestResponseFunctionCallArgumentsDoneEvent(t *testing.T) {
	data := `{
    "event_id": "event_5556",
    "type": "response.function_call_arguments.done",
    "response_id": "resp_002",
    "item_id": "fc_001",
    "output_index": 1,
    "call_id": "call_001",
    "arguments": "{\"location\": \"San Francisco\"}"
}`
	expected := openairt.ResponseFunctionCallArgumentsDoneEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_5556",
			Type:    openairt.ServerEventTypeResponseFunctionCallArgumentsDone,
		},
		ResponseID:  "resp_002",
		ItemID:      "fc_001",
		OutputIndex: 1,
		CallID:      "call_001",
		Arguments:   "{\"location\": \"San Francisco\"}",
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeResponseFunctionCallArgumentsDone, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.ResponseFunctionCallArgumentsDoneEvent))
}

func TestRateLimitsUpdatedEvent(t *testing.T) {
	data := `{
    "event_id": "event_5758",
    "type": "rate_limits.updated",
    "rate_limits": [
        {
            "name": "requests",
            "limit": 1000,
            "remaining": 999,
            "reset_seconds": 60.0
        },
        {
            "name": "tokens",
            "limit": 50000,
            "remaining": 49950,
            "reset_seconds": 60.0
        }
    ]
}`
	expected := openairt.RateLimitsUpdatedEvent{
		ServerEventBase: openairt.ServerEventBase{
			EventID: "event_5758",
			Type:    openairt.ServerEventTypeRateLimitsUpdated,
		},
		RateLimits: []openairt.RateLimit{
			{
				Name:         "requests",
				Limit:        1000,
				Remaining:    999,
				ResetSeconds: 60,
			},
			{
				Name:         "tokens",
				Limit:        50000,
				Remaining:    49950,
				ResetSeconds: 60,
			},
		},
	}
	actual, err := openairt.UnmarshalServerEvent([]byte(data))
	require.NoError(t, err)
	require.Equal(t, openairt.ServerEventTypeRateLimitsUpdated, actual.ServerEventType())
	require.Equal(t, expected, actual.(openairt.RateLimitsUpdatedEvent))
}
