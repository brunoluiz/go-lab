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
)

type Env string

const (
	EnvProduction Env = "production"
	EnvLocal      Env = "local"
	EnvTest       Env = "test"
)

type Exec interface {
	Run(ctx context.Context, logger *slog.Logger) error
}

func Run[T Exec](exec T) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize OpenTelemetry
	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to setup otel", "error", err)
		os.Exit(1)
	}
	defer closer.WithLogContext(ctx, logger, "failed to shutdown otel", otelShutdown)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		obsServer := o11y.New()
		defer closer.WithLogContext(ctx, logger, "failed to shutdown o11y server", obsServer.Close)
		return obsServer.Run(ctx)
	})

	eg.Go(func() error {
		kong.Parse(exec)
		return exec.Run(ctx, logger)
	})

	if err := eg.Wait(); err != nil {
		logger.ErrorContext(ctx, "application error", "error", err)
		os.Exit(1)
	}
}
