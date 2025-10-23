package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/brunoluiz/go-lab/lib/closer"
	"github.com/brunoluiz/go-lab/lib/o11y"
	"github.com/brunoluiz/go-lab/lib/otel"
	"golang.org/x/sync/errgroup"

	"github.com/hellofresh/health-go/v5"
)

type Env string

const (
	EnvProduction Env = "production"
	EnvLocal      Env = "local"
	EnvTest       Env = "test"
)

type Exec interface {
	Run(ctx context.Context, logger *slog.Logger, healthz *health.Health) error
}

func Run[T Exec](exec T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	kong.Parse(exec)

	// Initialize OpenTelemetry
	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to setup otel", "error", err)
		os.Exit(1)
	}
	defer closer.WithLogContext(ctx, logger, "failed to shutdown otel", otelShutdown)

	healthz, err := health.New()
	if err != nil {
		logger.ErrorContext(ctx, "failed to setup health checker", "error", err)
		os.Exit(1)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return o11y.Run(ctx, logger, healthz)
	})

	eg.Go(func() error {
		return exec.Run(ctx, logger, healthz)
	})

	if err := eg.Wait(); err != nil {
		logger.ErrorContext(ctx, "application error", "error", err)
		os.Exit(1)
	}
}
