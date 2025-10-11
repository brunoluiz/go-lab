package todo

import (
	"context"
	"testing"

	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoService(t *testing.T) {
	kv := database.NewKVStore()
	repo := repository.NewTaskRepository(kv)
	service := NewService(repo)
	ctx := context.Background()

	t.Run("CreateTask validation", func(t *testing.T) {
		req := model.CreateTaskRequest{Title: ""}
		_, err := service.CreateTask(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, ErrTitleRequired, err)
	})

	t.Run("CreateTask success", func(t *testing.T) {
		req := model.CreateTaskRequest{Title: "Valid Task"}
		resp, err := service.CreateTask(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Task.ID)
		assert.Equal(t, "Valid Task", resp.Task.Title)
	})

	t.Run("GetTask validation", func(t *testing.T) {
		req := model.GetTaskRequest{TaskID: ""}
		_, err := service.GetTask(ctx, req)
		assert.Error(t, err)
	})

	t.Run("UpdateTask validation", func(t *testing.T) {
		req := model.UpdateTaskRequest{Task: model.Task{ID: "", Title: ""}}
		_, err := service.UpdateTask(ctx, req)
		assert.Error(t, err)
	})

	t.Run("UpdateTask title required", func(t *testing.T) {
		req := model.UpdateTaskRequest{Task: model.Task{ID: "123", Title: ""}}
		_, err := service.UpdateTask(ctx, req)
		assert.Equal(t, ErrTitleRequired, err)
	})
}
