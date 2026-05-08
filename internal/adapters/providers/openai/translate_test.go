package openai

import (
	"testing"

	"github.com/ixr/ixr/pkg/schema"
)

func TestToWireRequest(t *testing.T) {
	tests := []struct {
		name     string
		req      *schema.RequestEnvelope
		wantMsgs int
		wantModel string
	}{
		{
			name: "single user message",
			req: &schema.RequestEnvelope{
				Model:    "gpt-4o",
				Messages: []schema.Message{{Role: "user", Content: "hello"}},
			},
			wantMsgs:  1,
			wantModel: "gpt-4o",
		},
		{
			name: "system + user messages both pass through",
			req: &schema.RequestEnvelope{
				Model: "gpt-4o",
				Messages: []schema.Message{
					{Role: "system", Content: "you are helpful"},
					{Role: "user", Content: "hi"},
				},
			},
			wantMsgs:  2,
			wantModel: "gpt-4o",
		},
		{
			name:      "empty messages",
			req:       &schema.RequestEnvelope{Model: "gpt-4o"},
			wantMsgs:  0,
			wantModel: "gpt-4o",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := toWireRequest(tc.req)
			if got.Model != tc.wantModel {
				t.Errorf("model: got %q, want %q", got.Model, tc.wantModel)
			}
			if len(got.Messages) != tc.wantMsgs {
				t.Errorf("messages: got %d, want %d", len(got.Messages), tc.wantMsgs)
			}
		})
	}
}

func TestFromWireResponse(t *testing.T) {
	wr := &wireResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Created: 1700000000,
		Model:   "gpt-4o",
		Choices: []wireChoice{
			{
				Index:        0,
				Message:      wireMessage{Role: "assistant", Content: "hello there"},
				FinishReason: "stop",
			},
		},
		Usage: wireUsage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
	}

	got := fromWireResponse(wr)

	if got.ID != "chatcmpl-123" {
		t.Errorf("ID: got %q, want %q", got.ID, "chatcmpl-123")
	}
	if len(got.Choices) != 1 {
		t.Fatalf("choices: got %d, want 1", len(got.Choices))
	}
	if got.Choices[0].Message.Content != "hello there" {
		t.Errorf("content: got %q, want %q", got.Choices[0].Message.Content, "hello there")
	}
	if got.Usage.TotalTokens != 15 {
		t.Errorf("total_tokens: got %d, want 15", got.Usage.TotalTokens)
	}
}
