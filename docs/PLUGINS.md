# Writing ixr plugins

A plugin is a Go type that implements `pkg/plugin.EventConsumer`. It receives every `CallEvent` emitted after an LLM call completes.

## the interface

```go
type EventConsumer interface {
    Name() string
    OnEvent(ctx context.Context, ev *schema.CallEvent) error
}
```

## example — under 30 lines

```go
package costcap

import (
    "context"
    "fmt"

    "github.com/ixr/ixr/pkg/plugin"
    "github.com/ixr/ixr/pkg/schema"
)

type CostCapPlugin struct {
    LimitUSD float64
}

func (c *CostCapPlugin) Name() string { return "cost-cap" }

func (c *CostCapPlugin) OnEvent(ctx context.Context, ev *schema.CallEvent) error {
    if ev.Cost.TotalUSD > c.LimitUSD {
        return fmt.Errorf("cost cap exceeded: %.4f > %.4f", ev.Cost.TotalUSD, c.LimitUSD)
    }
    return nil
}

// Verify interface compliance at compile time.
var _ plugin.EventConsumer = (*CostCapPlugin)(nil)
```

## guarantees

- `OnEvent` is called **asynchronously** — it never blocks the caller's response
- A **panicking plugin does not take down ixr** — the plugin manager catches and logs panics
- Returning an error from `OnEvent` is logged but does not affect the caller

## registering a plugin

Pass it to `ixr.Start()` via options (phase 1 day 5):

```go
ixr.Start(ixr.WithPlugins(&costcap.CostCapPlugin{LimitUSD: 1.00}))
```

## reference plugin

`plugins/audit-log/` ships with ixr. It writes every `CallEvent` as a JSON line to stdout.
It is the minimal proof that the bus works and a useful starting point for your own plugin.
