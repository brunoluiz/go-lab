package connectrpc

import (
	todov1connect "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1/todov1connect"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
)

type Handler struct {
	todov1connect.UnimplementedTodoServiceHandler

	todoService *todo.Service
	listService *list.Service
}

func NewHandler(todoService *todo.Service, listService *list.Service) *Handler {
	return &Handler{
		todoService: todoService,
		listService: listService,
	}
}
