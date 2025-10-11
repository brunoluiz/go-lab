package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	todov1connect "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1/todov1connect"
	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/grpc"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
)

func main() {
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
		Addr:              ":4000",
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
