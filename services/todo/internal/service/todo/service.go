package todo

import (
	"context"
	"errors"

	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/models"
)

var ErrTitleRequired = errors.New("title is required")

type Service struct {
	repo repository.TaskRepository
}

func NewService(repo repository.TaskRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(ctx context.Context, req models.CreateTaskRequest) (models.CreateTaskResponse, error) {
	if req.Title == "" {
		return models.CreateTaskResponse{}, ErrTitleRequired
	}
	return s.repo.CreateTask(ctx, req)
}

func (s *Service) GetTask(ctx context.Context, req models.GetTaskRequest) (models.GetTaskResponse, error) {
	if req.TaskID == "" {
		return models.GetTaskResponse{}, errors.New("task ID is required")
	}
	return s.repo.GetTask(ctx, req)
}

func (s *Service) ListTasks(ctx context.Context, req models.ListTasksRequest) (models.ListTasksResponse, error) {
	return s.repo.ListTasks(ctx, req)
}

func (s *Service) UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (models.UpdateTaskResponse, error) {
	if req.Task.ID == "" {
		return models.UpdateTaskResponse{}, errors.New("task ID is required")
	}
	if req.Task.Title == "" {
		return models.UpdateTaskResponse{}, ErrTitleRequired
	}
	return s.repo.UpdateTask(ctx, req)
}

func (s *Service) DeleteTask(ctx context.Context, req models.DeleteTaskRequest) (models.DeleteTaskResponse, error) {
	if req.TaskID == "" {
		return models.DeleteTaskResponse{}, errors.New("task ID is required")
	}
	return s.repo.DeleteTask(ctx, req)
}
