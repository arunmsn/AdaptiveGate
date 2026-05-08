package ingress

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ixr/ixr/pkg/provider"
	"github.com/ixr/ixr/pkg/schema"
)

// Router picks a provider for a given model name.
type Router func(model string) (provider.Provider, error)

// ChatHandler handles POST /v1/chat/completions.
// It is OpenAI-compatible: existing SDKs point at ixr with no code changes.
type ChatHandler struct {
	router Router
}

// NewChatHandler creates a handler that delegates to router for provider selection.
func NewChatHandler(router Router) *ChatHandler {
	return &ChatHandler{router: router}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is supported")
		return
	}

	var req schema.RequestEnvelope
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request_body", "could not parse request JSON")
		return
	}

	if req.Model == "" {
		writeError(w, http.StatusBadRequest, "missing_model", "model field is required")
		return
	}

	// Streaming is phase 2 — reject early with a clear message.
	if req.Stream {
		writeError(w, http.StatusNotImplemented, "streaming_not_supported", "streaming is not yet supported; set stream=false")
		return
	}

	p, err := h.router(req.Model)
	if err != nil {
		writeError(w, http.StatusBadRequest, "no_provider", err.Error())
		return
	}

	resp, err := p.Chat(r.Context(), &req)
	if err != nil {
		slog.Error("provider error", "provider", p.Name(), "model", req.Model, "err", err)
		writeError(w, http.StatusBadGateway, "provider_error", "upstream provider returned an error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to write response", "err", err)
	}
}

// apiError matches the OpenAI error envelope so existing SDKs parse it correctly.
type apiError struct {
	Error apiErrorBody `json:"error"`
}

type apiErrorBody struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
}

func writeError(w http.ResponseWriter, status int, errType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError{
		Error: apiErrorBody{Message: message, Type: errType},
	})
}
