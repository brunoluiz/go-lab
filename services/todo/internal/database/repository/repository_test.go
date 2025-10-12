package repository_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/brunoluiz/go-lab/services/todo/internal/database/migration"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if terminateErr := pgContainer.Terminate(ctx); terminateErr != nil {
			t.Logf("failed to terminate container: %s", terminateErr)
		}
	})

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	// Run migrations
	migrator, err := migration.NewMigrator(db, logger)
	require.NoError(t, err)

	err = migrator.Up()
	require.NoError(t, err)

	return db
}

func TestTaskRepository(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bobDB := bob.NewDB(db)
	repo := repository.NewTaskRepository(bobDB, logger)
	ctx := context.Background()

	t.Run("CreateTask", func(t *testing.T) {
		t.Parallel()
		task := model.Task{
			ID:          "test-create-id",
			Title:       "Test Create Task",
			IsCompleted: false,
			CreatedAt:   time.Now(),
		}

		created, err := repo.CreateTask(ctx, task)
		require.NoError(t, err)
		assert.Equal(t, task.ID, created.ID)
		assert.Equal(t, task.Title, created.Title)
		assert.Equal(t, task.IsCompleted, created.IsCompleted)
		assert.True(t, created.CreatedAt.Equal(task.CreatedAt))
	})

	t.Run("GetTask", func(t *testing.T) {
		t.Parallel()
		// First create a task
		task := model.Task{
			ID:          "test-get-id",
			Title:       "Test Get Task",
			IsCompleted: false,
			CreatedAt:   time.Now(),
		}
		_, err := repo.CreateTask(ctx, task)
		require.NoError(t, err)

		// Now get it
		retrieved, err := repo.GetTask(ctx, task.ID)
		require.NoError(t, err)
		assert.Equal(t, task.ID, retrieved.ID)
		assert.Equal(t, task.Title, retrieved.Title)
		assert.Equal(t, task.IsCompleted, retrieved.IsCompleted)
	})

	t.Run("GetTask_NotFound", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetTask(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.Equal(t, repository.ErrTaskNotFound, err)
	})

	t.Run("ListTasks", func(t *testing.T) {
		t.Parallel()
		// Create a few tasks
		tasks := []model.Task{
			{
				ID:          "list-test-1",
				Title:       "List Test 1",
				IsCompleted: false,
				CreatedAt:   time.Now(),
			},
			{
				ID:          "list-test-2",
				Title:       "List Test 2",
				IsCompleted: true,
				CreatedAt:   time.Now(),
			},
		}

		for _, task := range tasks {
			_, err := repo.CreateTask(ctx, task)
			require.NoError(t, err)
		}

		listed, err := repo.ListTasks(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(listed), 2)

		// Check that our tasks are in the list
		found1 := false
		found2 := false
		for _, task := range listed {
			if task.ID == "list-test-1" {
				found1 = true
				assert.Equal(t, "List Test 1", task.Title)
			}
			if task.ID == "list-test-2" {
				found2 = true
				assert.Equal(t, "List Test 2", task.Title)
				assert.True(t, task.IsCompleted)
			}
		}
		assert.True(t, found1, "Task 1 should be in the list")
		assert.True(t, found2, "Task 2 should be in the list")
	})

	t.Run("UpdateTask", func(t *testing.T) {
		t.Parallel()
		// Create a task
		task := model.Task{
			ID:          "test-update-id",
			Title:       "Test Update Task",
			IsCompleted: false,
			CreatedAt:   time.Now(),
		}
		_, err := repo.CreateTask(ctx, task)
		require.NoError(t, err)

		// Update it
		updatedTask := model.Task{
			ID:          "test-update-id",
			Title:       "Updated Title",
			IsCompleted: true,
			CreatedAt:   task.CreatedAt,
		}
		updated, err := repo.UpdateTask(ctx, updatedTask)
		require.NoError(t, err)
		assert.Equal(t, "test-update-id", updated.ID)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.True(t, updated.IsCompleted)
	})

	t.Run("UpdateTask_NotFound", func(t *testing.T) {
		t.Parallel()
		task := model.Task{
			ID:          "non-existent-update",
			Title:       "Should not exist",
			IsCompleted: false,
			CreatedAt:   time.Now(),
		}
		_, err := repo.UpdateTask(ctx, task)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrTaskNotFound, err)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		t.Parallel()
		// Create a task
		task := model.Task{
			ID:          "test-delete-id",
			Title:       "Test Delete Task",
			IsCompleted: false,
			CreatedAt:   time.Now(),
		}
		_, err := repo.CreateTask(ctx, task)
		require.NoError(t, err)

		// Delete it
		err = repo.DeleteTask(ctx, task.ID)
		require.NoError(t, err)

		// Verify it's gone
		_, err = repo.GetTask(ctx, task.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrTaskNotFound, err)
	})

	t.Run("DeleteTask_NotFound", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteTask(ctx, "non-existent-delete")
		assert.Error(t, err)
		assert.Equal(t, repository.ErrTaskNotFound, err)
	})
}
