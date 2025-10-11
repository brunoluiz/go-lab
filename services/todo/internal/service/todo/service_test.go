package todo

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoService(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	kv := database.NewKVStore()
	repo := repository.NewTaskRepository(kv, logger)
	service := NewService(repo, logger)
	ctx := context.Background()

	t.Run("CreateTask validation", func(t *testing.T) {
		req := dto.CreateTaskRequest{Title: ""}
		_, err := service.CreateTask(ctx, req)
		assert.Error(t, err)
	})

	t.Run("CreateTask success", func(t *testing.T) {
		req := dto.CreateTaskRequest{Title: "Valid Task"}
		resp, err := service.CreateTask(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Task.ID)
		assert.Equal(t, "Valid Task", resp.Task.Title)
	})

	t.Run("GetTask validation", func(t *testing.T) {
		req := dto.GetTaskRequest{TaskID: ""}
		_, err := service.GetTask(ctx, req)
		assert.Error(t, err)
	})

	t.Run("UpdateTask validation", func(t *testing.T) {
		req := dto.UpdateTaskRequest{Task: dto.Task{ID: "", Title: ""}}
		_, err := service.UpdateTask(ctx, req)
		assert.Error(t, err)
	})

	t.Run("UpdateTask title required", func(t *testing.T) {
		req := dto.UpdateTaskRequest{Task: dto.Task{ID: "123", Title: ""}}
		_, err := service.UpdateTask(ctx, req)
		assert.Error(t, err)
	})
}
