package repository

import (
	"context"
	"testing"

	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository(t *testing.T) {
	kv := database.NewKVStore()
	repo := NewTaskRepository(kv)
	ctx := context.Background()

	t.Run("CreateTask", func(t *testing.T) {
		req := model.CreateTaskRequest{Title: "Test Task"}
		resp, err := repo.CreateTask(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Task.ID)
		assert.Equal(t, "Test Task", resp.Task.Title)
		assert.False(t, resp.Task.IsCompleted)
		assert.NotZero(t, resp.Task.CreatedAt)
	})

	t.Run("GetTask", func(t *testing.T) {
		req := model.CreateTaskRequest{Title: "Get Test"}
		createResp, err := repo.CreateTask(ctx, req)
		require.NoError(t, err)

		getReq := model.GetTaskRequest{TaskID: createResp.Task.ID}
		getResp, err := repo.GetTask(ctx, getReq)
		require.NoError(t, err)
		assert.Equal(t, createResp.Task.ID, getResp.Task.ID)
		assert.Equal(t, createResp.Task.Title, getResp.Task.Title)
		assert.Equal(t, createResp.Task.IsCompleted, getResp.Task.IsCompleted)
		assert.True(t, createResp.Task.CreatedAt.Equal(getResp.Task.CreatedAt))
	})

	t.Run("ListTasks", func(t *testing.T) {
		req := model.ListTasksRequest{}
		resp, err := repo.ListTasks(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "default", resp.TodoList.Name)
		assert.GreaterOrEqual(t, len(resp.TodoList.Tasks), 1)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		req := model.CreateTaskRequest{Title: "Update Test"}
		createResp, err := repo.CreateTask(ctx, req)
		require.NoError(t, err)

		task := createResp.Task
		task.Title = "Updated Title"
		task.IsCompleted = true
		updateReq := model.UpdateTaskRequest{Task: task}
		updateResp, err := repo.UpdateTask(ctx, updateReq)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updateResp.Task.Title)
		assert.True(t, updateResp.Task.IsCompleted)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		req := model.CreateTaskRequest{Title: "Delete Test"}
		createResp, err := repo.CreateTask(ctx, req)
		require.NoError(t, err)

		deleteReq := model.DeleteTaskRequest{TaskID: createResp.Task.ID}
		_, err = repo.DeleteTask(ctx, deleteReq)
		require.NoError(t, err)

		getReq := model.GetTaskRequest{TaskID: createResp.Task.ID}
		_, err = repo.GetTask(ctx, getReq)
		assert.Error(t, err)
	})
}
