package repository

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/brunoluiz/go-lab/services/todo/internal/database"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	kv := database.NewKVStore()
	repo := NewTaskRepository(kv, logger)
	ctx := context.Background()

	tests := []struct {
		name    string
		prepare func() (model.Task, error)
		run     func(t *testing.T, task model.Task)
	}{
		{
			name: "CreateTask",
			prepare: func() (model.Task, error) {
				id1, err := uuid.NewV7()
				if err != nil {
					return model.Task{}, err
				}
				return model.Task{
					ID:          id1.String(),
					Title:       "Test Task",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil
			},
			run: func(t *testing.T, task model.Task) {
				resp, err := repo.CreateTask(ctx, task)
				require.NoError(t, err)
				assert.Equal(t, task.ID, resp.ID)
				assert.Equal(t, "Test Task", resp.Title)
				assert.False(t, resp.IsCompleted)
			},
		},
		{
			name: "GetTask",
			prepare: func() (model.Task, error) {
				id1, err := uuid.NewV7()
				if err != nil {
					return model.Task{}, err
				}
				return model.Task{
					ID:          id1.String(),
					Title:       "Get Test",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil
			},
			run: func(t *testing.T, task model.Task) {
				_, err := repo.CreateTask(ctx, task)
				require.NoError(t, err)

				getResp, err := repo.GetTask(ctx, task.ID)
				require.NoError(t, err)
				assert.Equal(t, task.ID, getResp.ID)
				assert.Equal(t, task.Title, getResp.Title)
				assert.Equal(t, task.IsCompleted, getResp.IsCompleted)
			},
		},
		{
			name: "ListTasks",
			prepare: func() (model.Task, error) {
				return model.Task{}, nil
			},
			run: func(t *testing.T, task model.Task) {
				resp, err := repo.ListTasks(ctx)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(resp), 1)
			},
		},
		{
			name: "UpdateTask",
			prepare: func() (model.Task, error) {
				id2, err := uuid.NewV7()
				if err != nil {
					return model.Task{}, err
				}
				return model.Task{
					ID:          id2.String(),
					Title:       "Update Test",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil
			},
			run: func(t *testing.T, task model.Task) {
				_, err := repo.CreateTask(ctx, task)
				require.NoError(t, err)

				task.Title = "Updated Title"
				task.IsCompleted = true
				updateResp, err := repo.UpdateTask(ctx, task)
				require.NoError(t, err)
				assert.Equal(t, "Updated Title", updateResp.Title)
				assert.True(t, updateResp.IsCompleted)
			},
		},
		{
			name: "DeleteTask",
			prepare: func() (model.Task, error) {
				id, err := uuid.NewV7()
				if err != nil {
					return model.Task{}, err
				}
				return model.Task{
					ID:          id.String(),
					Title:       "Delete Test",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil
			},
			run: func(t *testing.T, task model.Task) {
				_, err := repo.CreateTask(ctx, task)
				require.NoError(t, err)

				err = repo.DeleteTask(ctx, task.ID)
				require.NoError(t, err)

				_, err = repo.GetTask(ctx, task.ID)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			task, err := tt.prepare()
			require.NoError(t, err)
			tt.run(t, task)
		})
	}
}
