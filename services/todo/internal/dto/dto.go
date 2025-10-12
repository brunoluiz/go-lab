package dto

import "time"

type CreateTaskRequest struct {
	Title string `validate:"required"`
}

type CreateTaskResponse struct {
	Task Task
}

type GetTaskRequest struct {
	TaskID string `validate:"required"`
}

type GetTaskResponse struct {
	Task Task
}

type ListTasksRequest struct {
	// Empty for now
}

type ListTasksResponse struct {
	TodoList TodoList
}

type UpdateTaskRequest struct {
	Task Task `validate:"required"`
}

type UpdateTaskResponse struct {
	Task Task
}

type DeleteTaskRequest struct {
	TaskID string `validate:"required"`
}

type DeleteTaskResponse struct {
	// Empty
}

type Task struct {
	ID          string    `validate:"required"`
	Title       string    `validate:"required"`
	IsCompleted bool      `validate:"-"`
	CreatedAt   time.Time `validate:"-"`
}

type TodoList struct {
	Tasks []Task
	Name  string
}
