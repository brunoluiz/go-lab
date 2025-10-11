package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/model"
	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, req model.CreateTaskRequest) (model.CreateTaskResponse, error)
	GetTask(ctx context.Context, req model.GetTaskRequest) (model.GetTaskResponse, error)
	ListTasks(ctx context.Context, req model.ListTasksRequest) (model.ListTasksResponse, error)
	UpdateTask(ctx context.Context, req model.UpdateTaskRequest) (model.UpdateTaskResponse, error)
	DeleteTask(ctx context.Context, req model.DeleteTaskRequest) (model.DeleteTaskResponse, error)
}

type taskRepository struct {
	kv *database.KVStore
}

func NewTaskRepository(kv *database.KVStore) TaskRepository {
	return &taskRepository{kv: kv}
}

func (r *taskRepository) CreateTask(ctx context.Context, req model.CreateTaskRequest) (model.CreateTaskResponse, error) {
	_ = ctx
	id, err := uuid.NewV7()
	if err != nil {
		return model.CreateTaskResponse{}, fmt.Errorf("failed to generate UUID: %w", err)
	}
	idStr := id.String()
	task := model.Task{
		ID:          idStr,
		Title:       req.Title,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}
	data, err := json.Marshal(task)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal task", "error", err)
		return model.CreateTaskResponse{}, fmt.Errorf("failed to marshal task: %w", err)
	}
	key := fmt.Sprintf("task:%s", idStr)
	r.kv.Set(key, data)
	return model.CreateTaskResponse{Task: task}, nil
}

func (r *taskRepository) GetTask(ctx context.Context, req model.GetTaskRequest) (model.GetTaskResponse, error) {
	_ = ctx
	key := fmt.Sprintf("task:%s", req.TaskID)
	data, err := r.kv.Get(key)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get task", "error", err, "task_id", req.TaskID)
		return model.GetTaskResponse{}, fmt.Errorf("failed to get task: %w", err)
	}
	var task model.Task
	if err := json.Unmarshal(data, &task); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal task", "error", err, "task_id", req.TaskID)
		return model.GetTaskResponse{}, fmt.Errorf("failed to unmarshal task: %w", err)
	}
	return model.GetTaskResponse{Task: task}, nil
}

func (r *taskRepository) ListTasks(ctx context.Context, req model.ListTasksRequest) (model.ListTasksResponse, error) {
	_ = ctx
	_ = req
	tasksMap, err := r.kv.List("task:")
	if err != nil {
		return model.ListTasksResponse{}, fmt.Errorf("failed to list tasks: %w", err)
	}
	var tasks []model.Task
	for _, data := range tasksMap {
		var task model.Task
		if err := json.Unmarshal(data, &task); err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal task in list", "error", err)
			return model.ListTasksResponse{}, fmt.Errorf("failed to unmarshal task: %w", err)
		}
		tasks = append(tasks, task)
	}
	todoList := model.TodoList{
		Tasks: tasks,
		Name:  "default",
	}
	return model.ListTasksResponse{TodoList: todoList}, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, req model.UpdateTaskRequest) (model.UpdateTaskResponse, error) {
	_ = ctx
	data, err := json.Marshal(req.Task)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal task for update", "error", err, "task_id", req.Task.ID)
		return model.UpdateTaskResponse{}, fmt.Errorf("failed to marshal task: %w", err)
	}
	key := fmt.Sprintf("task:%s", req.Task.ID)
	r.kv.Set(key, data)
	return model.UpdateTaskResponse{Task: req.Task}, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, req model.DeleteTaskRequest) (model.DeleteTaskResponse, error) {
	_ = ctx
	key := fmt.Sprintf("task:%s", req.TaskID)
	if err := r.kv.Delete(key); err != nil {
		slog.ErrorContext(ctx, "failed to delete task", "error", err, "task_id", req.TaskID)
		return model.DeleteTaskResponse{}, fmt.Errorf("failed to delete task: %w", err)
	}
	return model.DeleteTaskResponse{}, nil
}
