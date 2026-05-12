// Package ixr is the one-line entry point for embedding ixr in a Go service.
//
//	import ixr "github.com/YashVishwas/ixr/pkg/ixr"
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

	auditlog "github.com/YashVishwas/ixr/plugins/audit-log"

	"github.com/YashVishwas/ixr/internal/adapters/bus"
	cfgloader "github.com/YashVishwas/ixr/internal/adapters/config"
	"github.com/YashVishwas/ixr/internal/adapters/pluginmgr"
	"github.com/YashVishwas/ixr/internal/adapters/providers/anthropic"
	"github.com/YashVishwas/ixr/internal/adapters/providers/deepseek"
	"github.com/YashVishwas/ixr/internal/adapters/providers/googleai"
	"github.com/YashVishwas/ixr/internal/adapters/providers/llama"
	"github.com/YashVishwas/ixr/internal/adapters/providers/openai"
	"github.com/YashVishwas/ixr/internal/ingress"
	"github.com/YashVishwas/ixr/pkg/provider"
)

// Option configures the ixr instance.
type Option func(*config)

type config struct {
	port       int
	configFile string
}

// WithPort overrides the listen port (default: 7000).
func WithPort(port int) Option {
	return func(c *config) { c.port = port }
}

// WithConfigFile loads configuration from the given ixr.yaml path.
// Provider credentials in the file may use ${ENV_VAR} syntax.
func WithConfigFile(path string) Option {
	return func(c *config) { c.configFile = path }
}

// Start starts the ixr proxy and blocks until the process receives SIGINT/SIGTERM
// or a fatal error occurs. It is the one-line entry point for embedding ixr.
func Start(opts ...Option) error {
	cfg := &config{port: 7000}
	for _, o := range opts {
		o(cfg)
	}

	registry, port, err := buildRegistry(cfg)
	if err != nil {
		return err
	}

	router := ingress.Router(func(model string) (provider.Provider, error) {
		m := strings.ToLower(model)
		switch {
		case strings.HasPrefix(model, "gpt-") || strings.HasPrefix(model, "o1") || strings.HasPrefix(model, "o3"):
			p, ok := registry["openai"]
			if !ok {
				return nil, fmt.Errorf("openai provider not configured")
			}
			return p, nil
		case strings.HasPrefix(model, "claude-"):
			p, ok := registry["anthropic"]
			if !ok {
				return nil, fmt.Errorf("anthropic provider not configured")
			}
			return p, nil
		case strings.HasPrefix(m, "gemma") || strings.Contains(m, "gemma-"):
			p, ok := registry["gemma"]
			if !ok {
				return nil, fmt.Errorf("gemma provider not configured")
			}
			return p, nil
		case strings.HasPrefix(m, "gemini"):
			p, ok := registry["gemini"]
			if !ok {
				return nil, fmt.Errorf("gemini provider not configured")
			}
			return p, nil
		case strings.Contains(m, "llama"):
			p, ok := registry["llama"]
			if !ok {
				return nil, fmt.Errorf("llama provider not configured (use GROQ_API_KEY or ixr.yaml providers.llama)")
			}
			return p, nil
		case strings.HasPrefix(m, "deepseek"):
			p, ok := registry["deepseek"]
			if !ok {
				return nil, fmt.Errorf("deepseek provider not configured")
			}
			return p, nil
		default:
			return nil, fmt.Errorf("no provider found for model %q", model)
		}
	})

	memBus := bus.NewMemory(0)
	mgr := pluginmgr.New(memBus)
	mgr.Register(&auditlog.Plugin{})

	mux := http.NewServeMux()
	mux.Handle("POST /v1/chat/completions", ingress.NewChatHandler(router, memBus))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go memBus.Start(ctx)

	return ingress.NewServer(port, mux).Run(ctx)
}

// buildRegistry constructs the provider map and effective port from config file or env vars.
func buildRegistry(cfg *config) (map[string]provider.Provider, int, error) {
	// Try config file first: explicit path → auto-discover → fall back to env.
	var fileCfg *cfgloader.Config
	var err error

	if cfg.configFile != "" {
		fileCfg, err = cfgloader.Load(cfg.configFile)
		if err != nil {
			return nil, 0, err
		}
	} else {
		fileCfg, err = cfgloader.Discover()
		if err != nil {
			return nil, 0, err
		}
	}

	registry := map[string]provider.Provider{}
	port := cfg.port

	if fileCfg != nil {
		if fileCfg.Server.Port != 0 && cfg.port == 7000 {
			port = fileCfg.Server.Port
		}
		for name, pc := range fileCfg.Providers {
			switch name {
			case "openai":
				if pc.APIKey != "" {
					registry["openai"] = openai.New(pc.APIKey, pc.BaseURL)
				}
			case "anthropic":
				if pc.APIKey != "" {
					registry["anthropic"] = anthropic.New(pc.APIKey, pc.BaseURL)
				}
			case "gemini":
				if pc.APIKey != "" {
					registry["gemini"] = googleai.NewGemini(pc.APIKey, pc.BaseURL)
				}
			case "gemma":
				if pc.APIKey != "" {
					registry["gemma"] = googleai.NewGemma(pc.APIKey, pc.BaseURL)
				}
			case "llama":
				if pc.APIKey != "" {
					registry["llama"] = llama.New(pc.APIKey, pc.BaseURL)
				}
			case "deepseek":
				if pc.APIKey != "" {
					registry["deepseek"] = deepseek.New(pc.APIKey, pc.BaseURL)
				}
			}
		}
	}

	// Env vars supplement or override config file.
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		registry["openai"] = openai.New(key, "")
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		registry["anthropic"] = anthropic.New(key, "")
	}
	googleKey := os.Getenv("GOOGLE_API_KEY")
	if googleKey == "" {
		googleKey = os.Getenv("GEMINI_API_KEY")
	}
	if googleKey != "" {
		registry["gemini"] = googleai.NewGemini(googleKey, "")
		registry["gemma"] = googleai.NewGemma(googleKey, "")
	}
	if key := os.Getenv("GROQ_API_KEY"); key != "" {
		registry["llama"] = llama.New(key, "")
	}
	if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
		registry["deepseek"] = deepseek.New(key, "")
	}

	if len(registry) == 0 {
		return nil, 0, fmt.Errorf("ixr: no providers configured — set API keys (e.g. OPENAI_API_KEY, ANTHROPIC_API_KEY, GOOGLE_API_KEY, GROQ_API_KEY, DEEPSEEK_API_KEY) or provide ixr.yaml")
	}

	return registry, port, nil
}
