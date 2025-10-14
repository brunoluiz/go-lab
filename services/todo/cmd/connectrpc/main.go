package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/alecthomas/kong"
	"github.com/brunoluiz/go-lab/core/storage/postgres"
	todov1connect "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1/todov1connect"
	"github.com/brunoluiz/go-lab/services/todo/internal/connectrpc"
	"github.com/brunoluiz/go-lab/services/todo/internal/connectrpc/interceptor"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/otel"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
)

type CLI struct {
	Address string `kong:"default=0.0.0.0,env=ADDRESS"`
	Port    int    `kong:"default=4000,env=PORT"`
	DBDSN   string `kong:"default=postgres://postgres:password@localhost:5432/todo?sslmode=disable,env=DB_DSN"`
}

func run(cli *CLI, logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Initialize OpenTelemetry
	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup OpenTelemetry: %w", err)
	}
	defer func() {
		if shutdownErr := otelShutdown(ctx); shutdownErr != nil {
			logger.ErrorContext(ctx, "failed to shutdown OpenTelemetry", "error", shutdownErr)
		}
	}()

	sqlDB, err := postgres.New(postgres.EnvConfig{
		DSN: cli.DBDSN,
	}, postgres.WithLiveCheck())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer sqlDB.Close()

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return fmt.Errorf("failed to create otel interceptor: %w", err)
	}

	db := bob.NewDB(sqlDB)
	taskRepo := repository.NewTaskRepository(db, logger)
	listRepo := repository.NewListRepository(db, logger)
	listService := list.NewService(listRepo, logger)
	service := todo.NewService(taskRepo, listService, logger)

	// Setup Connect Handler
	grpcHandler := connectrpc.NewHandler(service, listService)
	path, h := todov1connect.NewTodoServiceHandler(grpcHandler, connect.WithInterceptors(
		interceptor.Error(logger),
		otelInterceptor,
	))
	mux := http.NewServeMux()
	mux.Handle(path, h)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cli.Address, cli.Port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		Protocols:         p,
	}

	go func() {
		logger.InfoContext(ctx, "Starting server", slog.String("address", cli.Address), slog.Int("port", cli.Port))
		if serr := server.ListenAndServe(); serr != nil && serr != http.ErrServerClosed {
			logger.ErrorContext(ctx, "Failure to serve", "error", serr)
		}
	}()

	<-ctx.Done()
	logger.InfoContext(ctx, "Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if serr := server.Shutdown(shutdownCtx); serr != nil {
		return fmt.Errorf("failure to shutdown: %w", serr)
	}
	logger.InfoContext(ctx, "Shutdown complete")
	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var cli CLI
	kong.Parse(&cli)

	if err := run(&cli, logger); err != nil {
		//nolint
		logger.Error("Application error", "error", err)
		os.Exit(1)
	}
}
