package app

import (
	"context"
	"fmt"
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

type global struct {
	O11yAddress string `kong:"default=0.0.0.0,env=O11Y_ADDRESS,name=o11y-address"`
	O11yPort    int    `kong:"default=9090,env=O11Y_PORT,name=o11y-port"`
}

type Exec interface {
	Run(ctx context.Context, logger *slog.Logger, healthz *health.Health) error
}

func Run[T Exec](exec T) {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(ctx, logger, exec); err != nil {
		logger.ErrorContext(ctx, "application error", "error", err)
		os.Exit(1)
	}
}

func run[T Exec](ctx context.Context, logger *slog.Logger, exec T) error {
	// Parses app flags, but also cfg (do not judge the trickery here)
	cfg := &global{}
	kong.Parse(exec, kong.Embed(cfg))

	// Initialize OpenTelemetry
	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup otel: %w", err)
	}
	defer closer.WithLogContext(ctx, logger, "failed to shutdown otel", otelShutdown)

	healthz, err := health.New()
	if err != nil {
		return fmt.Errorf("failed to setup health checker: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return o11y.Run(ctx, logger, healthz, o11y.WithAddr(fmt.Sprintf("%s:%d", cfg.O11yAddress, cfg.O11yPort)))
	})

	eg.Go(func() error {
		return exec.Run(ctx, logger, healthz)
	})

	if egErr := eg.Wait(); egErr != nil {
		return fmt.Errorf("application error: %w", egErr)
	}

	return nil
}
