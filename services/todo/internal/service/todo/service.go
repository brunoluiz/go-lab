package todo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/model"
)

var ErrTitleRequired = errors.New("title is required")

type Service struct {
	repo repository.TaskRepository
}

func NewService(repo repository.TaskRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(ctx context.Context, req model.CreateTaskRequest) (model.CreateTaskResponse, error) {
	if req.Title == "" {
		slog.ErrorContext(ctx, "create task validation failed", "error", ErrTitleRequired)
		return model.CreateTaskResponse{}, ErrTitleRequired
	}
	return s.repo.CreateTask(ctx, req)
}

func (s *Service) GetTask(ctx context.Context, req model.GetTaskRequest) (model.GetTaskResponse, error) {
	if req.TaskID == "" {
		err := errors.New("task ID is required")
		slog.ErrorContext(ctx, "get task validation failed", "error", err)
		return model.GetTaskResponse{}, err
	}
	return s.repo.GetTask(ctx, req)
}

func (s *Service) ListTasks(ctx context.Context, req model.ListTasksRequest) (model.ListTasksResponse, error) {
	return s.repo.ListTasks(ctx, req)
}

func (s *Service) UpdateTask(ctx context.Context, req model.UpdateTaskRequest) (model.UpdateTaskResponse, error) {
	if req.Task.ID == "" {
		err := errors.New("task ID is required")
		slog.ErrorContext(ctx, "update task validation failed", "error", err)
		return model.UpdateTaskResponse{}, err
	}
	if req.Task.Title == "" {
		slog.ErrorContext(ctx, "update task validation failed", "error", ErrTitleRequired)
		return model.UpdateTaskResponse{}, ErrTitleRequired
	}
	return s.repo.UpdateTask(ctx, req)
}

func (s *Service) DeleteTask(ctx context.Context, req model.DeleteTaskRequest) (model.DeleteTaskResponse, error) {
	if req.TaskID == "" {
		err := errors.New("task ID is required")
		slog.ErrorContext(ctx, "delete task validation failed", "error", err)
		return model.DeleteTaskResponse{}, err
	}
	return s.repo.DeleteTask(ctx, req)
}
