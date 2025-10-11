package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	todov1connect "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1/todov1connect"
	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/grpc"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
)

type CLI struct {
	Address string `kong:"default=0.0.0.0,env=ADDRESS"`
	Port    int    `kong:"default=4000,env=PORT"`
}

func main() {
	var cli CLI
	kong.Parse(&cli)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	kv := database.NewKVStore()
	repo := repository.NewTaskRepository(kv, logger)
	service := todo.NewService(repo, logger)
	handler := grpc.NewHandler(service)

	mux := http.NewServeMux()
	path, h := todov1connect.NewTodoServiceHandler(handler)
	mux.Handle(path, h)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cli.Address, cli.Port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Starting server on :4000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
