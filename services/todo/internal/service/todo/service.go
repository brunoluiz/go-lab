package todo

import (
	"context"
	"log/slog"
	"time"

	"github.com/brunoluiz/go-lab/lib/errx"
	"github.com/brunoluiz/go-lab/services/todo/internal/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/repo"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func toDtoTask(t model.Task) Task {
	return Task{
		ID:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		ListID:      t.ListID,
	}
}

func fromDtoTask(t Task) model.Task {
	return model.Task{
		ID:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		ListID:      t.ListID,
	}
}

type ListService interface {
	GetList(ctx context.Context, req list.GetListRequest) (list.GetListResponse, error)
}

type Service struct {
	taskRepo    repo.TaskRepository
	listService ListService
	logger      *slog.Logger
	validator   *validator.Validate
}

func NewService(taskRepo repo.TaskRepository, listService ListService, logger *slog.Logger, validator *validator.Validate) *Service {
	return &Service{
		taskRepo:    taskRepo,
		listService: listService,
		logger:      logger,
		validator:   validator,
	}
}

func (s *Service) CreateTask(ctx context.Context, req CreateTaskRequest) (CreateTaskResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return CreateTaskResponse{}, errx.ErrValidation.Wrap(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateTaskResponse{}, errx.ErrUnknown.Wrap(err)
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
		return CreateTaskResponse{}, err
	}
	return CreateTaskResponse{Task: toDtoTask(created)}, nil
}

func (s *Service) GetTask(ctx context.Context, req GetTaskRequest) (GetTaskResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return GetTaskResponse{}, errx.ErrValidation.Wrap(err)
	}

	task, err := s.taskRepo.GetTask(ctx, req.TaskID)
	if err != nil {
		return GetTaskResponse{}, err
	}
	return GetTaskResponse{Task: toDtoTask(task)}, nil
}

func (s *Service) ListTasks(ctx context.Context, req ListTasksRequest) (ListTasksResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return ListTasksResponse{}, errx.ErrValidation.Wrap(err)
	}

	getListReq := list.GetListRequest{ListID: req.ListID} //nolint:staticcheck
	listResp, err := s.listService.GetList(ctx, getListReq)
	if err != nil {
		return ListTasksResponse{}, err
	}
	tasks, err := s.taskRepo.ListTasks(ctx, req.ListID)
	if err != nil {
		return ListTasksResponse{}, err
	}
	dtoTasks := make([]Task, len(tasks))
	for i, t := range tasks {
		dtoTasks[i] = toDtoTask(t)
	}
	todoList := TodoList{
		Tasks: dtoTasks,
		Name:  listResp.List.Name,
		ID:    listResp.List.ID,
	}
	return ListTasksResponse{TodoList: todoList}, nil
}

func (s *Service) UpdateTask(ctx context.Context, req UpdateTaskRequest) (UpdateTaskResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return UpdateTaskResponse{}, errx.ErrValidation.Wrap(err)
	}

	task := fromDtoTask(req.Task)
	updated, err := s.taskRepo.UpdateTask(ctx, task)
	if err != nil {
		return UpdateTaskResponse{}, err
	}
	return UpdateTaskResponse{Task: toDtoTask(updated)}, nil
}

func (s *Service) DeleteTask(ctx context.Context, req DeleteTaskRequest) (DeleteTaskResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return DeleteTaskResponse{}, errx.ErrValidation.Wrap(err)
	}

	err := s.taskRepo.DeleteTask(ctx, req.TaskID)
	if err != nil {
		return DeleteTaskResponse{}, err
	}
	return DeleteTaskResponse{}, nil
}
