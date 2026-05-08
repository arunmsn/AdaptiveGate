package anthropic

import (
	"testing"

	"github.com/ixr/ixr/pkg/schema"
)

func TestToWireRequest_SystemLifted(t *testing.T) {
	req := &schema.RequestEnvelope{
		Model: "claude-3-5-sonnet-20241022",
		Messages: []schema.Message{
			{Role: "system", Content: "be concise"},
			{Role: "user", Content: "hello"},
			{Role: "assistant", Content: "hi"},
			{Role: "user", Content: "how are you"},
		},
	}

	got := toWireRequest(req)

	if got.System != "be concise" {
		t.Errorf("system: got %q, want %q", got.System, "be concise")
	}
	if len(got.Messages) != 3 {
		t.Errorf("messages: got %d, want 3 (system must be removed)", len(got.Messages))
	}
	if got.Messages[0].Role != "user" {
		t.Errorf("first message role: got %q, want user", got.Messages[0].Role)
	}
	if got.MaxTokens != defaultMaxTokens {
		t.Errorf("max_tokens: got %d, want %d", got.MaxTokens, defaultMaxTokens)
	}
}

func TestToWireRequest_NoSystem(t *testing.T) {
	req := &schema.RequestEnvelope{
		Model:    "claude-3-5-sonnet-20241022",
		Messages: []schema.Message{{Role: "user", Content: "hello"}},
	}

	got := toWireRequest(req)

	if got.System != "" {
		t.Errorf("system: got %q, want empty", got.System)
	}
	if len(got.Messages) != 1 {
		t.Errorf("messages: got %d, want 1", len(got.Messages))
	}
}

func TestFromWireResponse(t *testing.T) {
	wr := &wireResponse{
		ID:    "msg_123",
		Model: "claude-3-5-sonnet-20241022",
		Content: []wireContent{
			{Type: "text", Text: "hello there"},
		},
		StopReason: "end_turn",
		Usage:      wireUsage{InputTokens: 8, OutputTokens: 4},
	}

	got, err := fromWireResponse(wr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Choices) != 1 {
		t.Fatalf("choices: got %d, want 1", len(got.Choices))
	}
	if got.Choices[0].Message.Content != "hello there" {
		t.Errorf("content: got %q, want %q", got.Choices[0].Message.Content, "hello there")
	}
	if got.Choices[0].FinishReason != "stop" {
		t.Errorf("finish_reason: got %q, want stop", got.Choices[0].FinishReason)
	}
	if got.Usage.PromptTokens != 8 {
		t.Errorf("prompt_tokens: got %d, want 8", got.Usage.PromptTokens)
	}
}

func TestFromWireResponse_NoTextBlock(t *testing.T) {
	wr := &wireResponse{
		Content: []wireContent{{Type: "tool_use", Text: ""}},
	}
	_, err := fromWireResponse(wr)
	if err == nil {
		t.Fatal("expected error for missing text block, got nil")
	}
}

func TestNormalizeStopReason(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"end_turn", "stop"},
		{"max_tokens", "length"},
		{"stop_sequence", "stop"},
		{"tool_use", "tool_use"},
		{"unknown_reason", "unknown_reason"},
	}
	for _, tc := range tests {
		got := normalizeStopReason(tc.in)
		if got != tc.want {
			t.Errorf("normalizeStopReason(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
