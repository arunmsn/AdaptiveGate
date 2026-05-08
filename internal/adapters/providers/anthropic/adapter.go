// Package anthropic implements pkg/provider.Provider for the Anthropic Messages API.
// It translates ixr's canonical schema to/from Anthropic's wire format and normalises
// the response back to OpenAI shape so callers need no provider-specific logic.
package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/YashVishwas/ixr/pkg/schema"
)

const (
	defaultBaseURL        = "https://api.anthropic.com/v1"
	anthropicVersion      = "2023-06-01"
)

// Adapter implements pkg/provider.Provider for Anthropic.
type Adapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// New creates an Adapter using the given API key.
// Pass baseURL="" to use the default; override in tests.
func New(apiKey, baseURL string) *Adapter {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Adapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (a *Adapter) Name() string { return "anthropic" }

// Chat sends req to the Anthropic Messages API and returns a normalised response.
func (a *Adapter) Chat(ctx context.Context, req *schema.RequestEnvelope) (*schema.ResponseEnvelope, error) {
	body, err := json.Marshal(toWireRequest(req))
	if err != nil {
		return nil, fmt.Errorf("anthropic: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("anthropic: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", a.apiKey)
	httpReq.Header.Set("anthropic-version", anthropicVersion)

	httpResp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic: do request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("anthropic: read response body: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic: status %d: %s", httpResp.StatusCode, respBody)
	}

	var wireResp wireResponse
	if err := json.Unmarshal(respBody, &wireResp); err != nil {
		return nil, fmt.Errorf("anthropic: decode response: %w", err)
	}

	resp, err := fromWireResponse(&wireResp)
	if err != nil {
		return nil, fmt.Errorf("anthropic: normalise response: %w", err)
	}
	return resp, nil
}
