package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"golang.org/x/sync/errgroup"
)

type cliConfig struct{}

type CLIExec interface {
	Run(ctx context.Context, logger *slog.Logger) error
}

func CLI[T CLIExec](exec T) {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := runCLI(ctx, logger, exec); err != nil {
		logger.ErrorContext(ctx, "application error", "error", err)
		os.Exit(1)
	}
}

func runCLI[T CLIExec](ctx context.Context, logger *slog.Logger, exec T) error {
	// Parses app flags, but also cfg (do not judge the trickery here)
	cfg := &cliConfig{}
	kong.Parse(exec, kong.Embed(cfg))

	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(appCtx)
	eg.Go(func() error {
		defer cancel()
		return exec.Run(ctx, logger)
	})

	if egErr := eg.Wait(); egErr != nil {
		return fmt.Errorf("application error: %w", egErr)
	}

	return nil
}
