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
	"github.com/brunoluiz/go-lab/lib/logger"
	"github.com/brunoluiz/go-lab/lib/o11y"
	"github.com/brunoluiz/go-lab/lib/otel"
	"golang.org/x/sync/errgroup"

	"github.com/hellofresh/health-go/v5"
)

type Env string

const (
	EnvProduction Env = "production"
	EnvLocal      Env = "local"
)

type serverConfig struct {
	LogLevel    string `enum:"debug,info,warn,error" kong:"default=info,env=LOG_LEVEL,name=log-level"`
	O11yAddress string `kong:"default=0.0.0.0,env=O11Y_ADDRESS,name=o11y-address"`
	O11yPort    int    `kong:"default=9090,env=O11Y_PORT,name=o11y-port"`
}

type ServerExec interface {
	Run(ctx context.Context, logger *slog.Logger, healthz *health.Health) error
}

func Server[T ServerExec](exec T) {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg := &serverConfig{}
	kong.Parse(exec, kong.Embed(cfg))

	logger := logger.New(logger.WithLevel(cfg.LogLevel))
	if err := runServer(ctx, cfg, logger, exec); err != nil {
		logger.ErrorContext(ctx, "application error", "error", err)
		os.Exit(1)
	}
}

func runServer[T ServerExec](ctx context.Context, cfg *serverConfig, logger *slog.Logger, exec T) error {
	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup otel: %w", err)
	}
	defer closer.WithLogContext(ctx, logger, "failed to shutdown otel", otelShutdown)

	healthz, err := health.New()
	if err != nil {
		return fmt.Errorf("failed to setup health checker: %w", err)
	}

	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(appCtx)
	eg.Go(func() error {
		return o11y.Run(ctx, logger, healthz,
			o11y.WithAddr(fmt.Sprintf("%s:%d", cfg.O11yAddress, cfg.O11yPort)),
		)
	})

	eg.Go(func() error {
		return exec.Run(ctx, logger, healthz)
	})

	if egErr := eg.Wait(); egErr != nil {
		return fmt.Errorf("application error: %w", egErr)
	}

	return nil
}
