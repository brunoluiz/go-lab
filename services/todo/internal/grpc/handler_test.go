package grpc

import (
	"context"
	"log/slog"
	"os"
	"testing"

	v1 "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1"
	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	kv := database.NewKVStore()
	repo := repository.NewTaskRepository(kv, logger)
	service := todo.NewService(repo, logger)
	handler := NewHandler(service)
	ctx := context.Background()

	tests := []struct {
		name    string
		prepare func() string
		run     func(t *testing.T, title string)
	}{
		{
			name: "CreateTask",
			prepare: func() string {
				return "Handler Test"
			},
			run: func(t *testing.T, title string) {
				req := &v1.CreateTaskRequest{Title: title}
				resp, err := handler.CreateTask(ctx, req)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Task.Id)
				assert.Equal(t, title, resp.Task.Title)
				assert.False(t, resp.Task.IsCompleted)
				assert.NotNil(t, resp.Task.CreatedAt)
			},
		},
		{
			name: "GetTask",
			prepare: func() string {
				return "Get Handler Test"
			},
			run: func(t *testing.T, title string) {
				createReq := &v1.CreateTaskRequest{Title: title}
				createResp, err := handler.CreateTask(ctx, createReq)
				require.NoError(t, err)

				getReq := &v1.GetTaskRequest{TaskId: createResp.Task.Id}
				getResp, err := handler.GetTask(ctx, getReq)
				require.NoError(t, err)
				assert.Equal(t, createResp.Task, getResp.Task)
			},
		},
		{
			name: "ListTasks",
			prepare: func() string {
				return ""
			},
			run: func(t *testing.T, title string) {
				req := &v1.ListTasksRequest{}
				resp, err := handler.ListTasks(ctx, req)
				require.NoError(t, err)
				assert.Equal(t, "default", resp.TodoList.Name)
				assert.GreaterOrEqual(t, len(resp.TodoList.Tasks), 1)
			},
		},
		{
			name: "UpdateTask",
			prepare: func() string {
				return "Update Handler Test"
			},
			run: func(t *testing.T, title string) {
				createReq := &v1.CreateTaskRequest{Title: title}
				createResp, err := handler.CreateTask(ctx, createReq)
				require.NoError(t, err)

				task := createResp.Task
				task.Title = "Updated"
				task.IsCompleted = true
				updateReq := &v1.UpdateTaskRequest{Task: task}
				updateResp, err := handler.UpdateTask(ctx, updateReq)
				require.NoError(t, err)
				assert.Equal(t, "Updated", updateResp.Task.Title)
				assert.True(t, updateResp.Task.IsCompleted)
			},
		},
		{
			name: "DeleteTask",
			prepare: func() string {
				return "Delete Handler Test"
			},
			run: func(t *testing.T, title string) {
				createReq := &v1.CreateTaskRequest{Title: title}
				createResp, err := handler.CreateTask(ctx, createReq)
				require.NoError(t, err)

				deleteReq := &v1.DeleteTaskRequest{TaskId: createResp.Task.Id}
				_, err = handler.DeleteTask(ctx, deleteReq)
				require.NoError(t, err)

				getReq := &v1.GetTaskRequest{TaskId: createResp.Task.Id}
				_, err = handler.GetTask(ctx, getReq)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title := tt.prepare()
			tt.run(t, title)
		})
	}
}
