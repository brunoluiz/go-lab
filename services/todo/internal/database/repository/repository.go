package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInternal     = errors.New("internal error")
)

type TaskRepository interface {
	CreateTask(ctx context.Context, req model.Task) (model.Task, error)
	GetTask(ctx context.Context, id string) (model.Task, error)
	ListTasks(ctx context.Context) ([]model.Task, error)
	UpdateTask(ctx context.Context, task model.Task) (model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskRepository struct {
	kv     *database.KVStore
	logger *slog.Logger
}

func NewTaskRepository(kv *database.KVStore, logger *slog.Logger) TaskRepository {
	return &taskRepository{kv: kv, logger: logger}
}

func (r *taskRepository) CreateTask(ctx context.Context, task model.Task) (model.Task, error) {
	_ = ctx
	data, err := json.Marshal(task)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to marshal task", "error", err)
		return model.Task{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	key := fmt.Sprintf("task:%s", task.ID)
	r.kv.Set(key, data)
	return task, nil
}

func (r *taskRepository) GetTask(ctx context.Context, id string) (model.Task, error) {
	_ = ctx
	key := fmt.Sprintf("task:%s", id)
	data, err := r.kv.Get(key)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			r.logger.ErrorContext(ctx, "task not found", "task_id", id)
			return model.Task{}, ErrTaskNotFound
		}
		r.logger.ErrorContext(ctx, "failed to get task", "error", err, "task_id", id)
		return model.Task{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	var task model.Task
	if err := json.Unmarshal(data, &task); err != nil {
		r.logger.ErrorContext(ctx, "failed to unmarshal task", "error", err, "task_id", id)
		return model.Task{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return task, nil
}

func (r *taskRepository) ListTasks(ctx context.Context) ([]model.Task, error) {
	_ = ctx
	tasksMap, err := r.kv.List("task:")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	var tasks []model.Task
	for _, data := range tasksMap {
		var task model.Task
		if err := json.Unmarshal(data, &task); err != nil {
			r.logger.ErrorContext(ctx, "failed to unmarshal task in list", "error", err)
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task model.Task) (model.Task, error) {
	_ = ctx
	data, err := json.Marshal(task)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to marshal task for update", "error", err, "task_id", task.ID)
		return model.Task{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	key := fmt.Sprintf("task:%s", task.ID)
	r.kv.Set(key, data)
	return task, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id string) error {
	_ = ctx
	key := fmt.Sprintf("task:%s", id)
	if err := r.kv.Delete(key); err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			r.logger.ErrorContext(ctx, "task not found for deletion", "task_id", id)
			return ErrTaskNotFound
		}
		r.logger.ErrorContext(ctx, "failed to delete task", "error", err, "task_id", id)
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return nil
}
