// Package pluginmgr loads plugins at startup and dispatches events to them.
// A panicking plugin does not take down ixr; panics are caught and logged.
package pluginmgr

import (
	"github.com/ixr/ixr/pkg/bus"
	"github.com/ixr/ixr/pkg/plugin"
)

// Manager registers plugins with a Bus.
type Manager struct {
	b bus.Bus
}

// New creates a Manager that registers plugins on b.
func New(b bus.Bus) *Manager {
	return &Manager{b: b}
}

// Register subscribes c to the bus.
func (m *Manager) Register(c plugin.EventConsumer) {
	m.b.Subscribe(c)
}
