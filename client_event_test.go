package openairt_test

import (
	"encoding/json"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestSessionUpdateEvent(t *testing.T) {
	temperature := float32(0.5)
	message := openairt.SessionUpdateEvent{
		EventBase: openairt.EventBase{
			EventID: "test-id",
		},
		Session: openairt.ClientSession{
			Modalities:        []openairt.Modality{openairt.ModalityText, openairt.ModalityAudio},
			Instructions:      "test-instructions",
			Voice:             openairt.VoiceAlloy,
			InputAudioFormat:  openairt.AudioFormatPcm16,
			OutputAudioFormat: openairt.AudioFormatG711Ulaw,
			InputAudioTranscription: &openairt.InputAudioTranscription{
				Model: openai.Whisper1,
			},
			TurnDetection: &openairt.ClientTurnDetection{
				Type: openairt.ClientTurnDetectionTypeServerVad,
				TurnDetectionParams: openairt.TurnDetectionParams{
					Threshold:         0.5,
					PrefixPaddingMs:   1000,
					SilenceDurationMs: 2000,
				},
			},
			Tools: []openairt.Tool{
				{
					Type: openairt.ToolTypeFunction,
					Name: "test-tool",
					Parameters: jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"location": {
								Type: jsonschema.String,
							},
						},
						Required: []string{"location"},
					},
				},
			},
			ToolChoice: openairt.ToolChoice{
				Type: openairt.ToolTypeFunction,
				Function: openairt.ToolFunction{
					Name: "test-tool",
				},
			},
			Temperature:     &temperature,
			MaxOutputTokens: 100,
		},
	}

	// data, err := json.MarshalIndent(message, "", "\t")
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{
	"event_id": "test-id",
	"session": {
			"modalities": [
					"text",
					"audio"
			],
			"instructions": "test-instructions",
			"voice": "alloy",
			"input_audio_format": "pcm16",
			"output_audio_format": "g711_ulaw",
			"input_audio_transcription": {
					"model": "whisper-1"
			},
			"turn_detection": {
					"type": "server_vad",
					"threshold": 0.5,
					"prefix_padding_ms": 1000,
					"silence_duration_ms": 2000
			},
			"tools": [
					{
							"type": "function",
							"name": "test-tool",
							"description": "",
							"parameters": {
									"type": "object",
									"properties": {
											"location": {
													"type": "string"
											}
									},
									"required": [
											"location"
									]
							}
					}
			],
			"tool_choice": {
					"type": "function",
					"function": {
							"name": "test-tool"
					}
			},
			"temperature": 0.5,
			"max_response_output_tokens": 100
	},
	"type": "session.update"
}`
	require.JSONEq(t, expected, string(data))

	message.Session.MaxOutputTokens = 0
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{
		"event_id": "test-id",
		"session": {
				"modalities": [
						"text",
						"audio"
				],
				"instructions": "test-instructions",
				"voice": "alloy",
				"input_audio_format": "pcm16",
				"output_audio_format": "g711_ulaw",
				"input_audio_transcription": {
						"model": "whisper-1"
				},
				"turn_detection": {
						"type": "server_vad",
						"threshold": 0.5,
						"prefix_padding_ms": 1000,
						"silence_duration_ms": 2000
				},
				"tools": [
						{
								"type": "function",
								"name": "test-tool",
								"description": "",
								"parameters": {
										"type": "object",
										"properties": {
												"location": {
														"type": "string"
												}
										},
										"required": [
												"location"
										]
								}
						}
				],
				"tool_choice": {
						"type": "function",
						"function": {
								"name": "test-tool"
						}
				},
				"temperature": 0.5
		},
		"type": "session.update"
	}`
	require.JSONEq(t, expected, string(data))
}

func TestSessionUpdateEventSimple(t *testing.T) {
	temperature := float32(0.5)
	message := openairt.SessionUpdateEvent{
		Session: openairt.ClientSession{
			Modalities:              []openairt.Modality{openairt.ModalityText},
			Instructions:            "test-instructions",
			Voice:                   openairt.VoiceAlloy,
			InputAudioFormat:        openairt.AudioFormatPcm16,
			OutputAudioFormat:       openairt.AudioFormatG711Ulaw,
			InputAudioTranscription: nil,
			TurnDetection:           nil,
			Tools:                   nil,
			ToolChoice:              nil,
			Temperature:             &temperature,
			MaxOutputTokens:         100,
		},
	}

	// data, err := json.MarshalIndent(message, "", "\t")
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{
	"session": {
			"modalities": [
					"text"
			],
			"instructions": "test-instructions",
			"voice": "alloy",
			"input_audio_format": "pcm16",
			"output_audio_format": "g711_ulaw",
			"turn_detection": null,
			"temperature": 0.5,
			"max_response_output_tokens": 100
	},
	"type": "session.update"
}`
	require.JSONEq(t, expected, string(data))
}

func TestInputAudioBufferAppendEvent(t *testing.T) {
	message := openairt.InputAudioBufferAppendEvent{
		Audio: "test-audio",
	}
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{"audio":"test-audio","type":"input_audio_buffer.append"}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","audio":"test-audio","type":"input_audio_buffer.append"}`
	require.JSONEq(t, expected, string(data))
}

func TestInputAudioBufferCommitEvent(t *testing.T) {
	message := openairt.InputAudioBufferCommitEvent{}
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{"type":"input_audio_buffer.commit"}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","type":"input_audio_buffer.commit"}`
	require.JSONEq(t, expected, string(data))
}

func TestInputAudioBufferClearEvent(t *testing.T) {
	message := openairt.InputAudioBufferClearEvent{}
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{"type":"input_audio_buffer.clear"}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","type":"input_audio_buffer.clear"}`
	require.JSONEq(t, expected, string(data))
}

func TestConversationItemCreateEvent(t *testing.T) {
	message := openairt.ConversationItemCreateEvent{
		PreviousItemID: "test-previous-item-id",
		Item: openairt.MessageItem{
			ID:     "test-id",
			Type:   openairt.MessageItemTypeMessage,
			Status: openairt.ItemStatusCompleted,
			Role:   openairt.MessageRoleUser,
			Content: []openairt.MessageContentPart{
				{Type: openairt.MessageContentTypeText, Text: "test-content"},
				{Type: openairt.MessageContentTypeAudio, Audio: "test-audio"},
				{Type: openairt.MessageContentTypeTranscript, Transcript: "test-transcript"},
			},
		},
	}
	data, err := json.MarshalIndent(message, "", "\t")
	require.NoError(t, err)
	expected := `{
	"previous_item_id": "test-previous-item-id",
	"item": {
			"id": "test-id",
			"type": "message",
			"status": "completed",
			"role": "user",
			"content": [
					{
							"type": "text",
							"text": "test-content"
					},
					{
							"type": "audio",
							"audio": "test-audio"
					},
					{
							"type": "transcript",
							"transcript": "test-transcript"
					}
			]
	},
	"type": "conversation.item.create"
}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{
	"event_id": "test-id",
	"previous_item_id": "test-previous-item-id",
	"item": {
			"id": "test-id",
			"type": "message",
			"status": "completed",
			"role": "user",
			"content": [
					{
							"type": "text",
							"text": "test-content"
					},
					{
							"type": "audio",
							"audio": "test-audio"
					},
					{
							"type": "transcript",
							"transcript": "test-transcript"
					}
			]
	},
	"type": "conversation.item.create"
}`
	require.JSONEq(t, expected, string(data))
}

func TestConversationItemTruncateEvent(t *testing.T) {
	message := openairt.ConversationItemTruncateEvent{
		ItemID:       "test-item-id",
		ContentIndex: 1,
		AudioEndMs:   2000,
	}
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{"item_id":"test-item-id","content_index":1,"audio_end_ms":2000,"type":"conversation.item.truncate"}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","item_id":"test-item-id","content_index":1,"audio_end_ms":2000,"type":"conversation.item.truncate"}`
	require.JSONEq(t, expected, string(data))
}
func TestConversationItemDeleteEvent(t *testing.T) {
	message := openairt.ConversationItemDeleteEvent{
		ItemID: "test-item-id",
	}
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{"item_id":"test-item-id","type":"conversation.item.delete"}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","item_id":"test-item-id","type":"conversation.item.delete"}`
	require.JSONEq(t, expected, string(data))
}

func TestResponseCreateEvent(t *testing.T) {
	message := openairt.ResponseCreateEvent{
		Response: openairt.ResponseCreateParams{
			Modalities:        []openairt.Modality{openairt.ModalityText, openairt.ModalityAudio},
			Instructions:      "test-instructions",
			Voice:             openairt.VoiceAlloy,
			OutputAudioFormat: openairt.AudioFormatG711Ulaw,
			Tools:             nil,
			ToolChoice:        openairt.ToolChoiceAuto,
			Temperature:       nil,
			MaxOutputTokens:   100,
		},
	}
	data, err := json.MarshalIndent(message, "", "\t")
	require.NoError(t, err)
	expected := `{
	"response": {
		"modalities": [
				"text",
				"audio"
		],
		"instructions": "test-instructions",
		"voice": "alloy",
		"output_audio_format": "g711_ulaw",
		"tool_choice": "auto",
		"max_output_tokens": 100
	},
	"type": "response.create"
}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	message.Response.MaxOutputTokens = openairt.Inf
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","response":{"modalities":["text","audio"],"instructions":"test-instructions","voice":"alloy","output_audio_format":"g711_ulaw","tool_choice":"auto","max_output_tokens":"inf"},"type":"response.create"}`
	require.JSONEq(t, expected, string(data))
}

func TestResponseCancelEvent(t *testing.T) {
	message := openairt.ResponseCancelEvent{}
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{"type":"response.cancel"}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","type":"response.cancel"}`
	require.JSONEq(t, expected, string(data))
}
