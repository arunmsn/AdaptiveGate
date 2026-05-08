// Package ixr is the one-line entry point for embedding ixr in a Go service.
//
//	import ixr "github.com/ixr/ixr/pkg/ixr"
//
//	func main() {
//	    go ixr.Start()
//	}
package ixr

// Option configures the ixr instance.
type Option func(*config)

type config struct {
	configPath string
	port       int
}

// WithConfig sets the path to the ixr.yaml config file.
func WithConfig(path string) Option {
	return func(c *config) { c.configPath = path }
}

// WithPort overrides the listen port (default: 7000).
func WithPort(port int) Option {
	return func(c *config) { c.port = port }
}

// Start starts the ixr proxy with the given options and blocks until the
// context is cancelled or a fatal error occurs. It is the one-line entry
// point for embedding ixr in a Go service.
func Start(opts ...Option) error {
	cfg := &config{port: 7000}
	for _, o := range opts {
		o(cfg)
	}
	// TODO: wire up ingress, providers, bus, and plugins — phase 1 day 2–5
	panic("ixr.Start: not yet implemented")
}
