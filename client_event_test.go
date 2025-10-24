package openairt_test

import (
	"encoding/json"
	"testing"

	openairt "github.com/WqyJh/go-openai-realtime"
	"github.com/stretchr/testify/require"
)

func TestSessionUpdateEvent(t *testing.T) {
	message := openairt.SessionUpdateEvent{
		EventBase: openairt.EventBase{
			EventID: "test-id",
		},
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				OutputModalities: []openairt.Modality{openairt.ModalityText, openairt.ModalityAudio},
				Instructions:     "test-instructions",
				Audio: openairt.RealtimeSessionAudio{
					Input: &openairt.SessionAudioInput{
						Format: openairt.AudioFormatUnion{
							PCM: &openairt.AudioFormatPCM{
								Rate: 24000,
							},
						},
						Transcription: openairt.AudioTranscription{
							Model: openairt.Whisper1,
						},
						TurnDetection: openairt.TurnDetectionUnion{
							ServerVad: &openairt.ServerVad{
								Threshold:         0.5,
								PrefixPaddingMs:   1000,
								SilenceDurationMs: 2000,
							},
						},
					},
					Output: &openairt.SessionAudioOutput{
						Format: openairt.AudioFormatUnion{
							PCM: &openairt.AudioFormatPCM{
								Rate: 24000,
							},
						},
						Voice: openairt.VoiceAlloy,
					},
				},
				Tools: []openairt.ToolUnion{
					{
						Function: &openairt.ToolFunction{
							Name:        "display_color_palette",
							Description: "Call this function when a user asks for a color palette.",
							Parameters: map[string]any{
								"type":   "object",
								"strict": true,
								"properties": map[string]any{
									"theme": map[string]any{
										"type":        "string",
										"description": "Description of the theme for the color scheme.",
									},
									"colors": map[string]any{
										"type":        "array",
										"description": "Array of five hex color codes based on the theme.",
										"items": map[string]any{
											"type":        "string",
											"description": "Hex color code",
										},
									},
								},
								"required": []string{"theme", "colors"},
							},
						},
					},
				},
				ToolChoice: openairt.ToolChoiceUnion{
					Function: &openairt.ToolChoiceFunction{
						Name: "display_color_palette",
					},
				},
				MaxOutputTokens: 100,
			},
		},
	}

	// data, err := json.MarshalIndent(message, "", "\t")
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{
  "event_id": "test-id",
  "session": {
    "type": "realtime",
    "output_modalities": ["text", "audio"],
    "instructions": "test-instructions",
    "audio": {
      "input": {
        "format": {
          "type": "audio/pcm",
          "rate": 24000
        },
        "transcription": {
          "model": "whisper-1"
        },
        "turn_detection": {
          "type": "server_vad",
          "threshold": 0.5,
          "prefix_padding_ms": 1000,
          "silence_duration_ms": 2000
        }
      },
      "output": {
        "format": {
          "type": "audio/pcm",
          "rate": 24000
        },
        "voice": "alloy"
      }
    },
    "tools": [
      {
        "type": "function",
        "name": "display_color_palette",
        "description": "Call this function when a user asks for a color palette.",
        "parameters": {
          "type": "object",
          "strict": true,
          "properties": {
            "theme": {
              "type": "string",
              "description": "Description of the theme for the color scheme."
            },
            "colors": {
              "type": "array",
              "description": "Array of five hex color codes based on the theme.",
              "items": {
                "type": "string",
                "description": "Hex color code"
              }
            }
          },
          "required": ["theme", "colors"]
        }
      }
    ],
    "tool_choice": {
      "type": "function",
	  "name": "display_color_palette"
    },
    "max_output_tokens": 100
  },
  "type": "session.update"
}`
	require.JSONEq(t, expected, string(data))

	message.Session.Realtime.MaxOutputTokens = 0
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{
  "event_id": "test-id",
  "session": {
    "type": "realtime",
    "output_modalities": ["text", "audio"],
    "instructions": "test-instructions",
    "audio": {
      "input": {
        "format": {
          "type": "audio/pcm",
          "rate": 24000
        },
        "transcription": {
          "model": "whisper-1"
        },
        "turn_detection": {
          "type": "server_vad",
          "threshold": 0.5,
          "prefix_padding_ms": 1000,
          "silence_duration_ms": 2000
        }
      },
      "output": {
        "format": {
          "type": "audio/pcm",
          "rate": 24000
        },
        "voice": "alloy"
      }
    },
    "tools": [
      {
        "type": "function",
        "name": "display_color_palette",
        "description": "Call this function when a user asks for a color palette.",
        "parameters": {
          "type": "object",
          "strict": true,
          "properties": {
            "theme": {
              "type": "string",
              "description": "Description of the theme for the color scheme."
            },
            "colors": {
              "type": "array",
              "description": "Array of five hex color codes based on the theme.",
              "items": {
                "type": "string",
                "description": "Hex color code"
              }
            }
          },
          "required": ["theme", "colors"]
        }
      }
    ],
    "tool_choice": {
      "type": "function",
      "name": "display_color_palette"
    }
  },
  "type": "session.update"
}`
	require.JSONEq(t, expected, string(data))
}

func TestSessionUpdateEventSimple(t *testing.T) {
	message := openairt.SessionUpdateEvent{
		Session: openairt.SessionUnion{
			Realtime: &openairt.RealtimeSession{
				OutputModalities: []openairt.Modality{openairt.ModalityText},
				Instructions:     "test-instructions",
				Audio: openairt.RealtimeSessionAudio{
					Output: &openairt.SessionAudioOutput{
						Voice: openairt.VoiceAlloy,
					},
				},
				MaxOutputTokens: 100,
			},
		},
	}

	// data, err := json.MarshalIndent(message, "", "\t")
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{
	"session": {
			"type": "realtime",
			"output_modalities": [
					"text"
			],
			"instructions": "test-instructions",
			"audio": {
					"output": {
							"voice": "alloy"
					}
			},
			"max_output_tokens": 100
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
		Item: openairt.MessageItemUnion{
			User: &openairt.MessageItemUser{
				ID: "test-id",
				Content: []openairt.MessageContentInput{
					{Type: openairt.MessageContentTypeText, Text: "test-content"},
					{Type: openairt.MessageContentTypeAudio, Audio: "test-audio"},
					{Type: openairt.MessageContentTypeTranscript, Transcript: "test-transcript"},
				},
				Status: openairt.ItemStatusCompleted,
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
			Instructions:     "test-instructions",
			Tools:            []openairt.ToolUnion{},
			OutputModalities: []openairt.Modality{openairt.ModalityText},
			Metadata: map[string]string{
				"response_purpose": "summarization",
			},
			Conversation: "none",
			Input: []openairt.MessageItemUnion{
				{
					User: &openairt.MessageItemUser{
						Content: []openairt.MessageContentInput{
							{
								Type: openairt.MessageContentTypeInputText,
								Text: "Summarize the above message in one sentence.",
							},
						},
					},
				},
			},
		},
	}
	data, err := json.MarshalIndent(message, "", "\t")
	require.NoError(t, err)
	expected := `{
  "type": "response.create",
  "response": {
    "instructions": "test-instructions",
    "tools": [],
    "conversation": "none",
    "output_modalities": ["text"],
    "metadata": {
      "response_purpose": "summarization"
    },
    "input": [
      {
        "type": "message",
        "role": "user",
        "content": [
          {
            "type": "input_text",
            "text": "Summarize the above message in one sentence."
          }
        ]
      }
    ]
  }
}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	message.Response.MaxOutputTokens = openairt.Inf
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{
  "type": "response.create",
  "event_id": "test-id",
  "response": {
    "instructions": "test-instructions",
    "tools": [],
    "conversation": "none",
    "output_modalities": ["text"],
    "metadata": {
      "response_purpose": "summarization"
    },
    "input": [
      {
        "type": "message",
        "role": "user",
        "content": [
          {
            "type": "input_text",
            "text": "Summarize the above message in one sentence."
          }
        ]
      }
    ],
	"max_output_tokens": "inf"
  }
}`
	require.JSONEq(t, expected, string(data))
}

func TestResponseCancelEvent(t *testing.T) {
	message := openairt.ResponseCancelEvent{
		ResponseID: "test-response-id",
	}
	data, err := json.Marshal(message)
	require.NoError(t, err)
	expected := `{"response_id":"test-response-id","type":"response.cancel"}`
	require.JSONEq(t, expected, string(data))

	message.EventBase.EventID = "test-id"
	data, err = json.Marshal(message)
	require.NoError(t, err)
	expected = `{"event_id":"test-id","response_id":"test-response-id","type":"response.cancel"}`
	require.JSONEq(t, expected, string(data))
}
