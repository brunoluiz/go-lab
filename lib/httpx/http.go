package httpx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	*http.Server

	logger          *slog.Logger
	shutdownTimeout time.Duration
}

type ServerOption func(*Server)

func WithShutdownTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.shutdownTimeout = d
	}
}

func WithLogger(logger *slog.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

func NewServer(addr string, handler http.Handler, opts ...ServerOption) *Server {
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	s := &Server{
		Server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			Protocols:         p,
		},
		shutdownTimeout: 5 * time.Second,
		logger:          slog.Default(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		s.logger.InfoContext(ctx, "starting server", slog.String("address", s.Addr))
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		s.logger.InfoContext(ctx, "shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failure to shutdown: %w", err)
		}
		s.logger.InfoContext(ctx, "shutdown complete")
		return nil
	}
}

func (s *Server) Close(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()
	return s.Shutdown(shutdownCtx)
}
