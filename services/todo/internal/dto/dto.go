package dto

import "time"

type CreateTaskRequest struct {
	Title string
}

type CreateTaskResponse struct {
	Task Task
}

type GetTaskRequest struct {
	TaskID string
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
	Task Task
}

type UpdateTaskResponse struct {
	Task Task
}

type DeleteTaskRequest struct {
	TaskID string
}

type DeleteTaskResponse struct {
	// Empty
}

type Task struct {
	ID          string
	Title       string
	IsCompleted bool
	CreatedAt   time.Time
}

type TodoList struct {
	Tasks []Task
	Name  string
}
