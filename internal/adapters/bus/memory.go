package bus

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ixr/ixr/pkg/plugin"
	"github.com/ixr/ixr/pkg/schema"
)

const defaultBufferSize = 1024

// Memory is a non-blocking in-process event bus backed by a buffered channel.
// Start must be called once before the first Publish.
type Memory struct {
	ch          chan *schema.CallEvent
	subscribers []plugin.EventConsumer
	mu          sync.RWMutex
}

// NewMemory creates a Memory bus with the given buffer size.
// Pass 0 to use the default (1024).
func NewMemory(bufSize int) *Memory {
	if bufSize <= 0 {
		bufSize = defaultBufferSize
	}
	return &Memory{ch: make(chan *schema.CallEvent, bufSize)}
}

// Subscribe registers c to receive events. Must be called before Start.
func (m *Memory) Subscribe(c plugin.EventConsumer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers = append(m.subscribers, c)
}

// Publish enqueues ev for async delivery. Drops the event if the buffer is full.
func (m *Memory) Publish(_ context.Context, ev *schema.CallEvent) error {
	select {
	case m.ch <- ev:
	default:
		slog.Warn("bus: buffer full, dropping event", "id", ev.ID)
	}
	return nil
}

// Start launches the dispatch loop and blocks until ctx is cancelled.
// Call it in a goroutine: go bus.Start(ctx).
func (m *Memory) Start(ctx context.Context) {
	for {
		select {
		case ev := <-m.ch:
			m.dispatch(ctx, ev)
		case <-ctx.Done():
			// Drain remaining events before exit.
			for {
				select {
				case ev := <-m.ch:
					m.dispatch(ctx, ev)
				default:
					return
				}
			}
		}
	}
}

func (m *Memory) dispatch(ctx context.Context, ev *schema.CallEvent) {
	m.mu.RLock()
	subs := m.subscribers
	m.mu.RUnlock()

	for _, s := range subs {
		func(c plugin.EventConsumer) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("plugin panic", "plugin", c.Name(), "panic", r)
				}
			}()
			if err := c.OnEvent(ctx, ev); err != nil {
				slog.Error("plugin error", "plugin", c.Name(), "err", err)
			}
		}(s)
	}
}
