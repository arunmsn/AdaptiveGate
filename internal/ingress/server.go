// Package ingress is the transport / API boundary layer.
// It speaks HTTP, parses OpenAI-shaped requests, and hands off to the app layer.
// Nothing in here knows about providers or domain logic.
package ingress
