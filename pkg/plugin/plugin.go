// Package plugin defines the interface third parties implement to extend ixr
// without forking. A plugin is a single Go file that implements EventConsumer.
package plugin

import (
	"context"

	"github.com/YashVishwas/ixr/pkg/schema"
)

// EventConsumer receives every CallEvent emitted on the bus.
// Implementations must be safe for concurrent use.
// A panicking OnEvent is caught by the plugin manager and does not take down ixr.
type EventConsumer interface {
	// Name returns a stable identifier used for logging and health checks.
	Name() string
	// OnEvent is called asynchronously after every LLM call completes.
	// Returning an error is logged but does not affect the caller.
	OnEvent(ctx context.Context, ev *schema.CallEvent) error
}
