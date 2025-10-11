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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	kv := database.NewKVStore()
	repo := repository.NewTaskRepository(kv, logger)
	service := todo.NewService(repo, logger)
	handler := NewHandler(service)
	ctx := context.Background()

	t.Run("CreateTask", func(t *testing.T) {
		req := &v1.CreateTaskRequest{Title: "Handler Test"}
		resp, err := handler.CreateTask(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Task.Id)
		assert.Equal(t, "Handler Test", resp.Task.Title)
		assert.False(t, resp.Task.IsCompleted)
		assert.NotNil(t, resp.Task.CreatedAt)
	})

	t.Run("GetTask", func(t *testing.T) {
		createReq := &v1.CreateTaskRequest{Title: "Get Handler Test"}
		createResp, err := handler.CreateTask(ctx, createReq)
		require.NoError(t, err)

		getReq := &v1.GetTaskRequest{TaskId: createResp.Task.Id}
		getResp, err := handler.GetTask(ctx, getReq)
		require.NoError(t, err)
		assert.Equal(t, createResp.Task, getResp.Task)
	})

	t.Run("ListTasks", func(t *testing.T) {
		req := &v1.ListTasksRequest{}
		resp, err := handler.ListTasks(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "default", resp.TodoList.Name)
		assert.GreaterOrEqual(t, len(resp.TodoList.Tasks), 1)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		createReq := &v1.CreateTaskRequest{Title: "Update Handler Test"}
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
	})

	t.Run("DeleteTask", func(t *testing.T) {
		createReq := &v1.CreateTaskRequest{Title: "Delete Handler Test"}
		createResp, err := handler.CreateTask(ctx, createReq)
		require.NoError(t, err)

		deleteReq := &v1.DeleteTaskRequest{TaskId: createResp.Task.Id}
		_, err = handler.DeleteTask(ctx, deleteReq)
		require.NoError(t, err)

		getReq := &v1.GetTaskRequest{TaskId: createResp.Task.Id}
		_, err = handler.GetTask(ctx, getReq)
		assert.Error(t, err)
	})
}
