package todo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/brunoluiz/go-lab/lib/errx"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/dto"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ErrTitleRequired = errors.New("title is required")

func toDtoTask(t model.Task) dto.Task {
	return dto.Task{
		ID:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		ListID:      t.ListID,
	}
}

func fromDtoTask(t dto.Task) model.Task {
	return model.Task{
		ID:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		ListID:      t.ListID,
	}
}

type Service struct {
	taskRepo    repository.TaskRepository
	listService *list.Service
	logger      *slog.Logger
	validator   *validator.Validate
}

func NewService(taskRepo repository.TaskRepository, listService *list.Service, logger *slog.Logger) *Service {
	return &Service{
		taskRepo:    taskRepo,
		listService: listService,
		logger:      logger,
		validator:   validator.New(),
	}
}

func (s *Service) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (dto.CreateTaskResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.CreateTaskResponse{}, errx.ErrValidation.Wrapf(err, "validation error")
	}
	id, err := uuid.NewV7()
	if err != nil {
		return dto.CreateTaskResponse{}, errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	task := model.Task{
		ID:          id.String(),
		Title:       req.Title,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		ListID:      req.ListID,
	}
	created, err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return dto.CreateTaskResponse{}, err
	}
	return dto.CreateTaskResponse{Task: toDtoTask(created)}, nil
}

func (s *Service) GetTask(ctx context.Context, req dto.GetTaskRequest) (dto.GetTaskResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.GetTaskResponse{}, errx.ErrValidation.Wrapf(err, "validation error")
	}
	task, err := s.taskRepo.GetTask(ctx, req.TaskID)
	if err != nil {
		return dto.GetTaskResponse{}, err
	}
	return dto.GetTaskResponse{Task: toDtoTask(task)}, nil
}

func (s *Service) ListTasks(ctx context.Context, req dto.ListTasksRequest) (dto.ListTasksResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.ListTasksResponse{}, errx.ErrValidation.Wrapf(err, "validation error")
	}
	getListReq := dto.GetListRequest{ListID: req.ListID} //nolint:staticcheck
	listResp, err := s.listService.GetList(ctx, getListReq)
	if err != nil {
		return dto.ListTasksResponse{}, err
	}
	tasks, err := s.taskRepo.ListTasks(ctx, req.ListID)
	if err != nil {
		return dto.ListTasksResponse{}, err
	}
	dtoTasks := make([]dto.Task, len(tasks))
	for i, t := range tasks {
		dtoTasks[i] = toDtoTask(t)
	}
	todoList := dto.TodoList{
		Tasks: dtoTasks,
		Name:  listResp.List.Name,
		ID:    listResp.List.ID,
	}
	return dto.ListTasksResponse{TodoList: todoList}, nil
}

func (s *Service) UpdateTask(ctx context.Context, req dto.UpdateTaskRequest) (dto.UpdateTaskResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.UpdateTaskResponse{}, errx.ErrValidation.Wrapf(err, "validation error")
	}
	task := fromDtoTask(req.Task)
	updated, err := s.taskRepo.UpdateTask(ctx, task)
	if err != nil {
		return dto.UpdateTaskResponse{}, err
	}
	return dto.UpdateTaskResponse{Task: toDtoTask(updated)}, nil
}

func (s *Service) DeleteTask(ctx context.Context, req dto.DeleteTaskRequest) (dto.DeleteTaskResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.DeleteTaskResponse{}, errx.ErrValidation.Wrapf(err, "validation error")
	}
	err := s.taskRepo.DeleteTask(ctx, req.TaskID)
	if err != nil {
		return dto.DeleteTaskResponse{}, err
	}
	return dto.DeleteTaskResponse{}, nil
}
