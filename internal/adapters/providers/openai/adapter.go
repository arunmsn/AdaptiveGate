package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/YashVishwas/ixr/pkg/schema"
)

const defaultBaseURL = "https://api.openai.com/v1"

// Adapter implements pkg/provider.Provider for OpenAI.
type Adapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// New creates an Adapter using the given API key.
// baseURL is overridable for testing; pass "" to use the default.
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

func (a *Adapter) Name() string { return "openai" }

// Chat sends req to OpenAI and returns the normalized response.
func (a *Adapter) Chat(ctx context.Context, req *schema.RequestEnvelope) (*schema.ResponseEnvelope, error) {
	body, err := json.Marshal(toWireRequest(req))
	if err != nil {
		return nil, fmt.Errorf("openai: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)

	httpResp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai: do request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: read response body: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai: status %d: %s", httpResp.StatusCode, respBody)
	}

	var wireResp wireResponse
	if err := json.Unmarshal(respBody, &wireResp); err != nil {
		return nil, fmt.Errorf("openai: decode response: %w", err)
	}

	return fromWireResponse(&wireResp), nil
}
