package todo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/dto"
	"github.com/google/uuid"
)

var (
	ErrTitleRequired    = errors.New("title is required")
	ErrTaskNotFound     = errors.New("task not found")
	ErrValidationFailed = errors.New("validation failed")
	ErrInternal         = errors.New("internal error")
)

func toDtoTask(t model.Task) dto.Task {
	return dto.Task{
		ID:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
	}
}

func fromDtoTask(t dto.Task) model.Task {
	return model.Task{
		ID:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
	}
}

type Service struct {
	repo   repository.TaskRepository
	logger *slog.Logger
}

func NewService(repo repository.TaskRepository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (dto.CreateTaskResponse, error) {
	if req.Title == "" {
		return dto.CreateTaskResponse{}, fmt.Errorf("%w: %w", ErrValidationFailed, ErrTitleRequired)
	}
	id, err := uuid.NewV7()
	if err != nil {
		return dto.CreateTaskResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	task := model.Task{
		ID:          id.String(),
		Title:       req.Title,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}
	created, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return dto.CreateTaskResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return dto.CreateTaskResponse{Task: toDtoTask(created)}, nil
}

func (s *Service) GetTask(ctx context.Context, req dto.GetTaskRequest) (dto.GetTaskResponse, error) {
	if req.TaskID == "" {
		err := errors.New("task ID is required")
		return dto.GetTaskResponse{}, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	task, err := s.repo.GetTask(ctx, req.TaskID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return dto.GetTaskResponse{}, ErrTaskNotFound
		}
		return dto.GetTaskResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return dto.GetTaskResponse{Task: toDtoTask(task)}, nil
}

func (s *Service) ListTasks(ctx context.Context, _ dto.ListTasksRequest) (dto.ListTasksResponse, error) {
	tasks, err := s.repo.ListTasks(ctx)
	if err != nil {
		return dto.ListTasksResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	dtoTasks := make([]dto.Task, len(tasks))
	for i, t := range tasks {
		dtoTasks[i] = toDtoTask(t)
	}
	todoList := dto.TodoList{
		Tasks: dtoTasks,
		Name:  "default",
	}
	return dto.ListTasksResponse{TodoList: todoList}, nil
}

func (s *Service) UpdateTask(ctx context.Context, req dto.UpdateTaskRequest) (dto.UpdateTaskResponse, error) {
	if req.Task.ID == "" {
		err := errors.New("task ID is required")
		return dto.UpdateTaskResponse{}, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	if req.Task.Title == "" {
		return dto.UpdateTaskResponse{}, fmt.Errorf("%w: %w", ErrValidationFailed, ErrTitleRequired)
	}
	task := fromDtoTask(req.Task)
	updated, err := s.repo.UpdateTask(ctx, task)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return dto.UpdateTaskResponse{}, ErrTaskNotFound
		}
		return dto.UpdateTaskResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return dto.UpdateTaskResponse{Task: toDtoTask(updated)}, nil
}

func (s *Service) DeleteTask(ctx context.Context, req dto.DeleteTaskRequest) (dto.DeleteTaskResponse, error) {
	if req.TaskID == "" {
		err := errors.New("task ID is required")
		return dto.DeleteTaskResponse{}, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	err := s.repo.DeleteTask(ctx, req.TaskID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return dto.DeleteTaskResponse{}, ErrTaskNotFound
		}
		return dto.DeleteTaskResponse{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return dto.DeleteTaskResponse{}, nil
}

func (s *Service) Hello() string {
	return "oi"
}
