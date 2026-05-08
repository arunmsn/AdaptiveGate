package pluginmgr

import (
	"github.com/YashVishwas/ixr/pkg/plugin"
)

// Registry holds a list of plugins to be registered at startup.
type Registry struct {
	plugins []plugin.EventConsumer
}

// Add appends c to the registry.
func (r *Registry) Add(c plugin.EventConsumer) {
	r.plugins = append(r.plugins, c)
}

// RegisterAll subscribes every plugin in the registry with m.
func (r *Registry) RegisterAll(m *Manager) {
	for _, p := range r.plugins {
		m.Register(p)
	}
}
