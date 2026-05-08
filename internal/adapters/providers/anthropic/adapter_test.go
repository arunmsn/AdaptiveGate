package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YashVishwas/ixr/pkg/schema"
)

func TestAdapter_Chat_HappyPath(t *testing.T) {
	fixture := wireResponse{
		ID:    "msg-abc",
		Model: "claude-3-5-sonnet-20241022",
		Content: []wireContent{
			{Type: "text", Text: "hello from claude"},
		},
		StopReason: "end_turn",
		Usage:      wireUsage{InputTokens: 6, OutputTokens: 4},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") == "" {
			t.Error("missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") != anthropicVersion {
			t.Errorf("anthropic-version: got %q, want %q", r.Header.Get("anthropic-version"), anthropicVersion)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(fixture)
	}))
	defer srv.Close()

	a := New("test-key", srv.URL)
	req := &schema.RequestEnvelope{
		Model:    "claude-3-5-sonnet-20241022",
		Messages: []schema.Message{{Role: "user", Content: "hi"}},
	}

	resp, err := a.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Choices) != 1 {
		t.Fatalf("choices: got %d, want 1", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "hello from claude" {
		t.Errorf("content: got %q", resp.Choices[0].Message.Content)
	}
	if resp.Choices[0].FinishReason != "stop" {
		t.Errorf("finish_reason: got %q, want stop", resp.Choices[0].FinishReason)
	}
}

func TestAdapter_Chat_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"type":"error","error":{"type":"authentication_error"}}`, http.StatusUnauthorized)
	}))
	defer srv.Close()

	a := New("bad-key", srv.URL)
	_, err := a.Chat(context.Background(), &schema.RequestEnvelope{
		Model:    "claude-3-5-sonnet-20241022",
		Messages: []schema.Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error for 401, got nil")
	}
}
