// Package auditlog is the reference EventConsumer plugin shipped with ixr.
// It writes every CallEvent as a JSON line to stdout — proof that the bus works
// and a useful primitive for any operator who needs a simple audit trail.
package auditlog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/YashVishwas/ixr/pkg/schema"
)

// Plugin implements plugin.EventConsumer.
type Plugin struct{}

// Name returns the stable plugin identifier.
func (p *Plugin) Name() string { return "audit-log" }

// OnEvent writes ev as a single JSON line to stdout.
func (p *Plugin) OnEvent(_ context.Context, ev *schema.CallEvent) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("audit-log: marshal: %w", err)
	}
	fmt.Println(string(b))
	return nil
}
