package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/alecthomas/kong"
	todov1connect "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1/todov1connect"
	"github.com/brunoluiz/go-lab/lib/closer"
	"github.com/brunoluiz/go-lab/lib/database/postgres"
	"github.com/brunoluiz/go-lab/lib/handler/connectrpc/interceptor"
	"github.com/brunoluiz/go-lab/lib/httpx"
	"github.com/brunoluiz/go-lab/lib/otel"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/handler/connectrpc"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	defer closer.WithLogContext(ctx, logger, "failed to shutdown OpenTelemetry", otelShutdown)

	// Initialize Database
	sqlDB, err := postgres.New(cli.DBDSN, postgres.WithLiveCheck())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer closer.WithLog(ctx, logger, "failed to shutdown database/sql", sqlDB.Close)

	db := bob.NewDB(sqlDB)
	taskRepo := repository.NewTaskRepository(db, logger)
	listRepo := repository.NewListRepository(db, logger)
	validator := validator.New()
	listService := list.NewService(listRepo, logger, validator)
	todoService := todo.NewService(taskRepo, listService, logger, validator)

	// Setup Connect Handler
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return fmt.Errorf("failed to create otel interceptor: %w", err)
	}

	grpcHandler := connectrpc.NewHandler(todoService, listService)
	path, h := todov1connect.NewTodoServiceHandler(grpcHandler, connect.WithInterceptors(
		otelInterceptor,
		interceptor.ErrorLogger(logger),
	))
	mux := http.NewServeMux()
	mux.Handle(path, h)

	server := httpx.NewServer(fmt.Sprintf("%s:%d", cli.Address, cli.Port),
		otelhttp.NewHandler(mux, "server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
		httpx.WithLogger(logger),
	)
	defer closer.WithLogContext(ctx, logger, "failed to shutdown HTTP server", server.Close)

	return server.Run(ctx)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var cli CLI
	kong.Parse(&cli)

	if err := run(&cli, logger); err != nil {
		//nolint
		logger.Error("application error", "error", err)
		os.Exit(1)
	}
}
