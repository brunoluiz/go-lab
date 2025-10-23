package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	todov1connect "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1/todov1connect"
	"github.com/brunoluiz/go-lab/lib/app"
	"github.com/brunoluiz/go-lab/lib/closer"
	"github.com/brunoluiz/go-lab/lib/database/postgres"
	"github.com/brunoluiz/go-lab/lib/handler/connectrpc/interceptor"
	"github.com/brunoluiz/go-lab/lib/httpx"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/handler/connectrpc"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
	"github.com/go-playground/validator/v10"
	"github.com/hellofresh/health-go/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type CLI struct {
	Address string `kong:"default=0.0.0.0,env=ADDRESS"`
	Port    int    `kong:"default=4000,env=PORT"`
	DBDSN   string `kong:"default=postgres://todo_user:todo_pass@localhost:5432/todo?sslmode=disable,env=DB_DSN"`
}

func (cli *CLI) Run(ctx context.Context, logger *slog.Logger, healthz *health.Health) error {
	sqlDB, err := postgres.New(cli.DBDSN, logger, postgres.WithHealthChecker(healthz))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer closer.WithLog(ctx, logger, "failed to shutdown database/sql", sqlDB.Conn.Close)

	db := bob.NewDB(sqlDB.Conn)
	validator := validator.New()
	taskRepo := repository.NewTaskRepository(db, logger)
	listRepo := repository.NewListRepository(db, logger)
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

	server := httpx.New(fmt.Sprintf("%s:%d", cli.Address, cli.Port),
		otelhttp.NewHandler(mux, "server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
		httpx.WithLogger(logger),
	)
	defer closer.WithLogContext(ctx, logger, "failed to shutdown HTTP server", server.Close)

	return server.Run(ctx)
}

func main() {
	app.Run(&CLI{})
}
