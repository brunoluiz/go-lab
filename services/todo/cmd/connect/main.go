package main

import (
	"context"
	"log"
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
	kv := database.NewKVStore()
	repo := repository.NewTaskRepository(kv)
	service := todo.NewService(repo)
	handler := grpc.NewHandler(service)

	mux := http.NewServeMux()
	path, h := todov1connect.NewTodoServiceHandler(handler)
	mux.Handle(path, h)

	ctx, cancel := context.WithCancel(context.Background())

	server := &http.Server{
		Addr:              ":4000",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Println("Starting server on :4000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("Server stopped")
}
