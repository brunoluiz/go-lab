package main

import (
	"log"
	"net/http"

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

	log.Println("Starting server on :4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
