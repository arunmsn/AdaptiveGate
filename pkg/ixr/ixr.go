// Package ixr is the one-line entry point for embedding ixr in a Go service.
//
//	import ixr "github.com/ixr/ixr/pkg/ixr"
//
//	func main() {
//	    go ixr.Start()
//	}
package ixr

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	auditlog "github.com/ixr/ixr/plugins/audit-log"

	"github.com/ixr/ixr/internal/adapters/bus"
	"github.com/ixr/ixr/internal/adapters/pluginmgr"
	"github.com/ixr/ixr/internal/adapters/providers/anthropic"
	"github.com/ixr/ixr/internal/adapters/providers/openai"
	"github.com/ixr/ixr/internal/ingress"
	"github.com/ixr/ixr/pkg/provider"
)

// Option configures the ixr instance.
type Option func(*config)

type config struct {
	port int
}

// WithPort overrides the listen port (default: 7000).
func WithPort(port int) Option {
	return func(c *config) { c.port = port }
}

// Start starts the ixr proxy and blocks until the process receives SIGINT/SIGTERM
// or a fatal error occurs. It is the one-line entry point for embedding ixr.
func Start(opts ...Option) error {
	cfg := &config{port: 7000}
	for _, o := range opts {
		o(cfg)
	}

	// Build provider registry from environment.
	// Config-file loading comes in phase 1 day 5.
	registry := map[string]provider.Provider{}

	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		registry["openai"] = openai.New(key, "")
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		registry["anthropic"] = anthropic.New(key, "")
	}

	if len(registry) == 0 {
		return fmt.Errorf("ixr: no providers configured — set OPENAI_API_KEY and/or ANTHROPIC_API_KEY")
	}

	// Model-prefix router: gpt-* → openai, claude-* → anthropic (day 3).
	router := ingress.Router(func(model string) (provider.Provider, error) {
		switch {
		case strings.HasPrefix(model, "gpt-") || strings.HasPrefix(model, "o1") || strings.HasPrefix(model, "o3"):
			p, ok := registry["openai"]
			if !ok {
				return nil, fmt.Errorf("openai provider not configured; set OPENAI_API_KEY")
			}
			return p, nil
		case strings.HasPrefix(model, "claude-"):
			p, ok := registry["anthropic"]
			if !ok {
				return nil, fmt.Errorf("anthropic provider not configured; set ANTHROPIC_API_KEY")
			}
			return p, nil
		default:
			return nil, fmt.Errorf("no provider found for model %q", model)
		}
	})

	// Event bus + plugins.
	memBus := bus.NewMemory(0)
	mgr := pluginmgr.New(memBus)
	mgr.Register(&auditlog.Plugin{})

	mux := http.NewServeMux()
	mux.Handle("POST /v1/chat/completions", ingress.NewChatHandler(router, memBus))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go memBus.Start(ctx)

	return ingress.NewServer(cfg.port, mux).Run(ctx)
}
