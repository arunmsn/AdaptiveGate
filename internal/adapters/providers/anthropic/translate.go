package anthropic

import (
	"fmt"
	"time"

	"github.com/ixr/ixr/pkg/schema"
)

// Anthropic Messages API wire types — internal to this adapter.

type wireRequest struct {
	Model     string        `json:"model"`
	Messages  []wireMessage `json:"messages"`
	System    string        `json:"system,omitempty"`
	MaxTokens int           `json:"max_tokens"`
}

type wireMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type wireResponse struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Role       string         `json:"role"`
	Content    []wireContent  `json:"content"`
	Model      string         `json:"model"`
	StopReason string         `json:"stop_reason"`
	Usage      wireUsage      `json:"usage"`
}

type wireContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type wireUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// toWireRequest converts ixr's canonical envelope to the Anthropic Messages API format.
// System messages are lifted out of the messages array into the top-level system field,
// as required by the Anthropic API.
func toWireRequest(req *schema.RequestEnvelope) wireRequest {
	const defaultMaxTokens = 4096

	var system string
	var msgs []wireMessage

	for _, m := range req.Messages {
		if m.Role == "system" {
			system = m.Content
			continue
		}
		msgs = append(msgs, wireMessage{Role: m.Role, Content: m.Content})
	}

	return wireRequest{
		Model:     req.Model,
		Messages:  msgs,
		System:    system,
		MaxTokens: defaultMaxTokens,
	}
}

// fromWireResponse converts an Anthropic response to ixr's canonical envelope.
func fromWireResponse(wr *wireResponse) (*schema.ResponseEnvelope, error) {
	text, err := extractText(wr.Content)
	if err != nil {
		return nil, err
	}

	return &schema.ResponseEnvelope{
		ID:      wr.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   wr.Model,
		Choices: []schema.Choice{
			{
				Index:        0,
				Message:      schema.Message{Role: "assistant", Content: text},
				FinishReason: normalizeStopReason(wr.StopReason),
			},
		},
		Usage: schema.Usage{
			PromptTokens:     wr.Usage.InputTokens,
			CompletionTokens: wr.Usage.OutputTokens,
			TotalTokens:      wr.Usage.InputTokens + wr.Usage.OutputTokens,
		},
	}, nil
}

func extractText(content []wireContent) (string, error) {
	for _, c := range content {
		if c.Type == "text" {
			return c.Text, nil
		}
	}
	return "", fmt.Errorf("anthropic: no text content block in response")
}

// normalizeStopReason maps Anthropic stop reasons to OpenAI finish reasons
// so callers that branch on finish_reason don't need provider-specific logic.
func normalizeStopReason(reason string) string {
	switch reason {
	case "end_turn":
		return "stop"
	case "max_tokens":
		return "length"
	case "stop_sequence":
		return "stop"
	default:
		return reason
	}
}
