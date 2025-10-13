package connectrpc

import (
	"context"

	v1 "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1"
	"github.com/brunoluiz/go-lab/services/todo/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.CreateTaskResponse, error) {
	internalReq := dto.CreateTaskRequest{Title: req.Title, ListID: req.ListId}
	resp, err := h.todoService.CreateTask(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.CreateTaskResponse{Task: toProtoTask(resp.Task)}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.GetTaskResponse, error) {
	internalReq := dto.GetTaskRequest{TaskID: req.TaskId}
	resp, err := h.todoService.GetTask(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.GetTaskResponse{Task: toProtoTask(resp.Task)}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *v1.ListTasksRequest) (*v1.ListTasksResponse, error) {
	internalReq := dto.ListTasksRequest{ListID: req.ListId}
	resp, err := h.todoService.ListTasks(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	protoTasks := make([]*v1.Task, len(resp.TodoList.Tasks))
	for i, t := range resp.TodoList.Tasks {
		protoTasks[i] = toProtoTask(t)
	}
	return &v1.ListTasksResponse{TodoList: &v1.TodoList{Tasks: protoTasks, Name: resp.TodoList.Name, Id: resp.TodoList.ID}}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *v1.UpdateTaskRequest) (*v1.UpdateTaskResponse, error) {
	internalReq := dto.UpdateTaskRequest{Task: fromProtoTask(req.Task)}
	resp, err := h.todoService.UpdateTask(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateTaskResponse{Task: toProtoTask(resp.Task)}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *v1.DeleteTaskRequest) (*v1.DeleteTaskResponse, error) {
	internalReq := dto.DeleteTaskRequest{TaskID: req.TaskId}
	_, err := h.todoService.DeleteTask(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteTaskResponse{}, nil
}

func toProtoTask(t dto.Task) *v1.Task {
	return &v1.Task{
		Id:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		ListId:      t.ListID,
	}
}

func fromProtoTask(t *v1.Task) dto.Task {
	return dto.Task{
		ID:          t.Id,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt.AsTime(),
		ListID:      t.ListId,
	}
}
