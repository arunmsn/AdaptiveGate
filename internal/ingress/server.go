package ingress

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Server is the HTTP listener for ixr's ingress.
// It knows nothing about providers or domain logic — only HTTP.
type Server struct {
	port int
	mux  *http.ServeMux
}

// NewServer creates a Server. Register routes on mux before calling Run.
func NewServer(port int, mux *http.ServeMux) *Server {
	return &Server{port: port, mux: mux}
}

// Run starts listening and blocks until ctx is cancelled or a fatal error occurs.
func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.mux,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("ixr listening", "port", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutCtx)
	case err := <-errCh:
		return err
	}
}
