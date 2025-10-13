package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	v1 "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1"
	todov1connect "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1/todov1connect"
	"github.com/brunoluiz/go-lab/services/todo/internal/connectrpc"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/migration"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestServer(t *testing.T) todov1connect.TodoServiceClient {
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

	// Set up service
	sqlDB := db
	bobDB := bob.NewDB(sqlDB)
	taskRepo := repository.NewTaskRepository(bobDB, logger)
	listRepo := repository.NewListRepository(bobDB, logger)
	listService := list.NewService(listRepo, logger)
	service := todo.NewService(taskRepo, listService, logger)

	// Setup Connect Handler
	grpcHandler := connectrpc.NewHandler(service, listService)
	path, h := todov1connect.NewTodoServiceHandler(grpcHandler)
	mux := http.NewServeMux()
	mux.Handle(path, h)

	// Create test server
	server := httptest.NewServer(mux)
	t.Cleanup(func() {
		server.Close()
	})

	// Create client
	client := todov1connect.NewTodoServiceClient(
		server.Client(),
		server.URL,
	)

	return client
}

func TestIntegration(t *testing.T) {
	t.Parallel()

	client := setupTestServer(t)
	ctx := context.Background()

	// Create two lists
	list1Resp, err := client.CreateList(ctx, &v1.CreateListRequest{Name: "Work List"})
	require.NoError(t, err)
	list1 := list1Resp.List
	assert.Equal(t, "Work List", list1.Name)
	assert.NotEmpty(t, list1.Id)

	list2Resp, err := client.CreateList(ctx, &v1.CreateListRequest{Name: "Personal List"})
	require.NoError(t, err)
	list2 := list2Resp.List
	assert.Equal(t, "Personal List", list2.Name)
	assert.NotEmpty(t, list2.Id)

	// Create four tasks, two in each list
	task1Resp, err := client.CreateTask(ctx, &v1.CreateTaskRequest{Title: "Finish project", ListId: list1.Id})
	require.NoError(t, err)
	task1 := task1Resp.Task
	assert.Equal(t, "Finish project", task1.Title)
	assert.Equal(t, list1.Id, task1.ListId)
	assert.False(t, task1.IsCompleted)

	task2Resp, err := client.CreateTask(ctx, &v1.CreateTaskRequest{Title: "Review code", ListId: list1.Id})
	require.NoError(t, err)
	task2 := task2Resp.Task
	assert.Equal(t, "Review code", task2.Title)
	assert.Equal(t, list1.Id, task2.ListId)

	task3Resp, err := client.CreateTask(ctx, &v1.CreateTaskRequest{Title: "Buy groceries", ListId: list2.Id})
	require.NoError(t, err)
	task3 := task3Resp.Task
	assert.Equal(t, "Buy groceries", task3.Title)
	assert.Equal(t, list2.Id, task3.ListId)

	task4Resp, err := client.CreateTask(ctx, &v1.CreateTaskRequest{Title: "Call mom", ListId: list2.Id})
	require.NoError(t, err)
	task4 := task4Resp.Task
	assert.Equal(t, "Call mom", task4.Title)
	assert.Equal(t, list2.Id, task4.ListId)

	// Verify fetching tasks by list works
	listTasks1Resp, err := client.ListTasks(ctx, &v1.ListTasksRequest{ListId: list1.Id})
	require.NoError(t, err)
	todoList1 := listTasks1Resp.TodoList
	assert.Equal(t, "Work List", todoList1.Name)
	assert.Equal(t, list1.Id, todoList1.Id)
	assert.Len(t, todoList1.Tasks, 2)
	// Check tasks are present
	taskIDs := make(map[string]bool)
	for _, task := range todoList1.Tasks {
		taskIDs[task.Id] = true
	}
	assert.True(t, taskIDs[task1.Id])
	assert.True(t, taskIDs[task2.Id])

	listTasks2Resp, err := client.ListTasks(ctx, &v1.ListTasksRequest{ListId: list2.Id})
	require.NoError(t, err)
	todoList2 := listTasks2Resp.TodoList
	assert.Equal(t, "Personal List", todoList2.Name)
	assert.Equal(t, list2.Id, todoList2.Id)
	assert.Len(t, todoList2.Tasks, 2)
	taskIDs2 := make(map[string]bool)
	for _, task := range todoList2.Tasks {
		taskIDs2[task.Id] = true
	}
	assert.True(t, taskIDs2[task3.Id])
	assert.True(t, taskIDs2[task4.Id])

	// Update lists and tasks
	// Update list1 name
	list1.Name = "Updated Work List"
	updateList1Resp, err := client.UpdateList(ctx, &v1.UpdateListRequest{List: list1})
	require.NoError(t, err)
	assert.Equal(t, "Updated Work List", updateList1Resp.List.Name)

	// Update list2 name
	list2.Name = "Updated Personal List"
	updateList2Resp, err := client.UpdateList(ctx, &v1.UpdateListRequest{List: list2})
	require.NoError(t, err)
	assert.Equal(t, "Updated Personal List", updateList2Resp.List.Name)

	// Update task1
	task1.Title = "Finish urgent project"
	task1.IsCompleted = true
	updateTask1Resp, err := client.UpdateTask(ctx, &v1.UpdateTaskRequest{Task: task1})
	require.NoError(t, err)
	assert.Equal(t, "Finish urgent project", updateTask1Resp.Task.Title)
	assert.True(t, updateTask1Resp.Task.IsCompleted)

	// Update task3
	task3.Title = "Buy organic groceries"
	updateTask3Resp, err := client.UpdateTask(ctx, &v1.UpdateTaskRequest{Task: task3})
	require.NoError(t, err)
	assert.Equal(t, "Buy organic groceries", updateTask3Resp.Task.Title)

	// Verify updates
	// Check list1 updated
	listTasks1Resp2, err := client.ListTasks(ctx, &v1.ListTasksRequest{ListId: list1.Id})
	require.NoError(t, err)
	assert.Equal(t, "Updated Work List", listTasks1Resp2.TodoList.Name)

	// Check task1 updated
	found := false
	for _, task := range listTasks1Resp2.TodoList.Tasks {
		if task.Id == task1.Id {
			assert.Equal(t, "Finish urgent project", task.Title)
			assert.True(t, task.IsCompleted)
			found = true
		}
	}
	assert.True(t, found)

	// Check list2 updated
	listTasks2Resp2, err := client.ListTasks(ctx, &v1.ListTasksRequest{ListId: list2.Id})
	require.NoError(t, err)
	assert.Equal(t, "Updated Personal List", listTasks2Resp2.TodoList.Name)

	// Check task3 updated
	found = false
	for _, task := range listTasks2Resp2.TodoList.Tasks {
		if task.Id == task3.Id {
			assert.Equal(t, "Buy organic groceries", task.Title)
			found = true
		}
	}
	assert.True(t, found)

	// Delete tasks
	_, err = client.DeleteTask(ctx, &v1.DeleteTaskRequest{TaskId: task1.Id})
	require.NoError(t, err)
	_, err = client.DeleteTask(ctx, &v1.DeleteTaskRequest{TaskId: task2.Id})
	require.NoError(t, err)
	_, err = client.DeleteTask(ctx, &v1.DeleteTaskRequest{TaskId: task3.Id})
	require.NoError(t, err)
	_, err = client.DeleteTask(ctx, &v1.DeleteTaskRequest{TaskId: task4.Id})
	require.NoError(t, err)

	// Delete lists
	_, err = client.DeleteList(ctx, &v1.DeleteListRequest{ListId: list1.Id})
	require.NoError(t, err)
	_, err = client.DeleteList(ctx, &v1.DeleteListRequest{ListId: list2.Id})
	require.NoError(t, err)

	// Verify deletions - lists should be gone
	_, err = client.GetList(ctx, &v1.GetListRequest{ListId: list1.Id})
	assert.Error(t, err)
	_, err = client.GetList(ctx, &v1.GetListRequest{ListId: list2.Id})
	assert.Error(t, err)

	// Verify tasks are gone
	_, err = client.GetTask(ctx, &v1.GetTaskRequest{TaskId: task1.Id})
	assert.Error(t, err)
	_, err = client.GetTask(ctx, &v1.GetTaskRequest{TaskId: task2.Id})
	assert.Error(t, err)
	_, err = client.GetTask(ctx, &v1.GetTaskRequest{TaskId: task3.Id})
	assert.Error(t, err)
	_, err = client.GetTask(ctx, &v1.GetTaskRequest{TaskId: task4.Id})
	assert.Error(t, err)
}
