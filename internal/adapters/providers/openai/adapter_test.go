package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ixr/ixr/pkg/schema"
)

func TestAdapter_Chat_HappyPath(t *testing.T) {
	fixture := wireResponse{
		ID:      "chatcmpl-abc",
		Object:  "chat.completion",
		Created: 1700000000,
		Model:   "gpt-4o",
		Choices: []wireChoice{
			{Index: 0, Message: wireMessage{Role: "assistant", Content: "pong"}, FinishReason: "stop"},
		},
		Usage: wireUsage{PromptTokens: 5, CompletionTokens: 2, TotalTokens: 7},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(fixture)
	}))
	defer srv.Close()

	a := New("test-key", srv.URL)
	req := &schema.RequestEnvelope{
		Model:    "gpt-4o",
		Messages: []schema.Message{{Role: "user", Content: "ping"}},
	}

	resp, err := a.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "chatcmpl-abc" {
		t.Errorf("ID: got %q, want chatcmpl-abc", resp.ID)
	}
	if len(resp.Choices) != 1 || resp.Choices[0].Message.Content != "pong" {
		t.Errorf("unexpected choices: %+v", resp.Choices)
	}
	if resp.Usage.TotalTokens != 7 {
		t.Errorf("total_tokens: got %d, want 7", resp.Usage.TotalTokens)
	}
}

func TestAdapter_Chat_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"invalid key"}}`, http.StatusUnauthorized)
	}))
	defer srv.Close()

	a := New("bad-key", srv.URL)
	_, err := a.Chat(context.Background(), &schema.RequestEnvelope{Model: "gpt-4o"})
	if err == nil {
		t.Fatal("expected error for 401, got nil")
	}
}

func TestAdapter_Chat_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	a := New("key", srv.URL)
	_, err := a.Chat(context.Background(), &schema.RequestEnvelope{Model: "gpt-4o"})
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}
