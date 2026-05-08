// Package bus defines the event bus interface.
// The in-memory implementation ships in phase 1; nats/kafka adapters are phase 2.
// Swapping implementations never changes code that publishes or subscribes.
package bus

import (
	"context"

	"github.com/YashVishwas/ixr/pkg/plugin"
	"github.com/YashVishwas/ixr/pkg/schema"
)

// Bus delivers CallEvents to all registered EventConsumers.
// Publish is non-blocking; a slow consumer must not block the request path.
type Bus interface {
	// Publish enqueues ev for delivery to all subscribers.
	Publish(ctx context.Context, ev *schema.CallEvent) error
	// Subscribe registers c to receive all future events.
	// Must be called before the first Publish.
	Subscribe(c plugin.EventConsumer)
}
