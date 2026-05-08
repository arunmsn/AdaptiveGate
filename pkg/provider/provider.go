// Package provider defines the interface for LLM backend adapters.
// Each provider (openai, anthropic, bedrock, ...) is a self-contained adapter
// that implements this interface. Adding a new provider never touches ixr core.
package provider

import (
	"context"

	"github.com/ixr/ixr/pkg/schema"
)

// Provider translates ixr's canonical schema to and from a specific LLM provider's
// wire format, and executes the call.
type Provider interface {
	// Name returns the provider identifier used in config and CallEvent.Provider.
	Name() string
	// Chat executes a chat completion request and returns the normalized response.
	Chat(ctx context.Context, req *schema.RequestEnvelope) (*schema.ResponseEnvelope, error)
}
