package dto

import "time"

type CreateListRequest struct {
	Name string `validate:"required"`
}

type CreateListResponse struct {
	List List
}

type GetListRequest struct {
	ListID string `validate:"required"`
}

type GetListResponse struct {
	List List
}

type ListListsRequest struct {
	// Empty for now
}

type ListListsResponse struct {
	Lists []List
}

type UpdateListRequest struct {
	List List `validate:"required"`
}

type UpdateListResponse struct {
	List List
}

type DeleteListRequest struct {
	ListID string `validate:"required"`
}

type DeleteListResponse struct {
	// Empty
}

type CreateTaskRequest struct {
	Title  string `validate:"required"`
	ListID string `validate:"required"`
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
	ListID string `validate:"required"`
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

type List struct {
	ID        string    `validate:"required"`
	Name      string    `validate:"required"`
	CreatedAt time.Time `validate:"-"`
}

type Task struct {
	ID          string    `validate:"required"`
	Title       string    `validate:"required"`
	IsCompleted bool      `validate:"-"`
	CreatedAt   time.Time `validate:"-"`
	ListID      string    `validate:"required"`
}

type TodoList struct {
	Tasks []Task
	Name  string
	ID    string
}
