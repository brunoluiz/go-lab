package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/models"
	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, req models.CreateTaskRequest) (models.CreateTaskResponse, error)
	GetTask(ctx context.Context, req models.GetTaskRequest) (models.GetTaskResponse, error)
	ListTasks(ctx context.Context, req models.ListTasksRequest) (models.ListTasksResponse, error)
	UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (models.UpdateTaskResponse, error)
	DeleteTask(ctx context.Context, req models.DeleteTaskRequest) (models.DeleteTaskResponse, error)
}

type taskRepository struct {
	kv *database.KVStore
}

func NewTaskRepository(kv *database.KVStore) TaskRepository {
	return &taskRepository{kv: kv}
}

func (r *taskRepository) CreateTask(ctx context.Context, req models.CreateTaskRequest) (models.CreateTaskResponse, error) {
	_ = ctx
	id, err := uuid.NewV7()
	if err != nil {
		return models.CreateTaskResponse{}, fmt.Errorf("failed to generate UUID: %w", err)
	}
	idStr := id.String()
	task := models.Task{
		ID:          idStr,
		Title:       req.Title,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}
	data, err := json.Marshal(task)
	if err != nil {
		return models.CreateTaskResponse{}, fmt.Errorf("failed to marshal task: %w", err)
	}
	key := fmt.Sprintf("task:%s", idStr)
	r.kv.Set(key, data)
	return models.CreateTaskResponse{Task: task}, nil
}

func (r *taskRepository) GetTask(ctx context.Context, req models.GetTaskRequest) (models.GetTaskResponse, error) {
	_ = ctx
	key := fmt.Sprintf("task:%s", req.TaskID)
	data, err := r.kv.Get(key)
	if err != nil {
		return models.GetTaskResponse{}, fmt.Errorf("failed to get task: %w", err)
	}
	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return models.GetTaskResponse{}, fmt.Errorf("failed to unmarshal task: %w", err)
	}
	return models.GetTaskResponse{Task: task}, nil
}

func (r *taskRepository) ListTasks(ctx context.Context, req models.ListTasksRequest) (models.ListTasksResponse, error) {
	_ = ctx
	_ = req
	tasksMap, err := r.kv.List("task:")
	if err != nil {
		return models.ListTasksResponse{}, fmt.Errorf("failed to list tasks: %w", err)
	}
	var tasks []models.Task
	for _, data := range tasksMap {
		var task models.Task
		if err := json.Unmarshal(data, &task); err != nil {
			return models.ListTasksResponse{}, fmt.Errorf("failed to unmarshal task: %w", err)
		}
		tasks = append(tasks, task)
	}
	todoList := models.TodoList{
		Tasks: tasks,
		Name:  "default",
	}
	return models.ListTasksResponse{TodoList: todoList}, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (models.UpdateTaskResponse, error) {
	_ = ctx
	data, err := json.Marshal(req.Task)
	if err != nil {
		return models.UpdateTaskResponse{}, fmt.Errorf("failed to marshal task: %w", err)
	}
	key := fmt.Sprintf("task:%s", req.Task.ID)
	r.kv.Set(key, data)
	return models.UpdateTaskResponse{Task: req.Task}, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, req models.DeleteTaskRequest) (models.DeleteTaskResponse, error) {
	_ = ctx
	key := fmt.Sprintf("task:%s", req.TaskID)
	if err := r.kv.Delete(key); err != nil {
		return models.DeleteTaskResponse{}, fmt.Errorf("failed to delete task: %w", err)
	}
	return models.DeleteTaskResponse{}, nil
}
