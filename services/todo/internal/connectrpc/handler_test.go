package connectrpc_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	v1 "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1"
	"github.com/brunoluiz/go-lab/services/todo/internal/connectrpc"
	"github.com/brunoluiz/go-lab/services/todo/internal/connectrpc/mock"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name    string
		prepare func(mockRepo *mock.MockTaskRepository) string
		run     func(t *testing.T, title string, handler *connectrpc.Handler)
	}{
		{
			name: "CreateTask",
			prepare: func(mockRepo *mock.MockTaskRepository) string {
				mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(model.Task{
					ID:          "test-id",
					Title:       "Handler Test",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil)
				return "Handler Test"
			},
			run: func(t *testing.T, title string, handler *connectrpc.Handler) {
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
			prepare: func(mockRepo *mock.MockTaskRepository) string {
				createdTime := time.Now()
				mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(model.Task{
					ID:          "get-test-id",
					Title:       "Get Handler Test",
					IsCompleted: false,
					CreatedAt:   createdTime,
				}, nil)
				mockRepo.EXPECT().GetTask(ctx, "get-test-id").Return(model.Task{
					ID:          "get-test-id",
					Title:       "Get Handler Test",
					IsCompleted: false,
					CreatedAt:   createdTime,
				}, nil)
				return "Get Handler Test"
			},
			run: func(t *testing.T, title string, handler *connectrpc.Handler) {
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
			prepare: func(mockRepo *mock.MockTaskRepository) string {
				mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(model.Task{
					ID:          "list-test-id",
					Title:       "List Handler Test",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil)
				mockRepo.EXPECT().ListTasks(ctx).Return([]model.Task{
					{
						ID:          "list-test-id",
						Title:       "List Handler Test",
						IsCompleted: false,
						CreatedAt:   time.Now(),
					},
				}, nil)
				return "List Handler Test"
			},
			run: func(t *testing.T, title string, handler *connectrpc.Handler) {
				createReq := &v1.CreateTaskRequest{Title: title}
				_, err := handler.CreateTask(ctx, createReq)
				require.NoError(t, err)

				req := &v1.ListTasksRequest{}
				resp, err := handler.ListTasks(ctx, req)
				require.NoError(t, err)
				assert.Equal(t, "default", resp.TodoList.Name)
				assert.GreaterOrEqual(t, len(resp.TodoList.Tasks), 1)
			},
		},
		{
			name: "UpdateTask",
			prepare: func(mockRepo *mock.MockTaskRepository) string {
				mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(model.Task{
					ID:          "update-test-id",
					Title:       "Update Handler Test",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil)
				mockRepo.EXPECT().UpdateTask(ctx, gomock.Any()).Return(model.Task{
					ID:          "update-test-id",
					Title:       "Updated",
					IsCompleted: true,
					CreatedAt:   time.Now(),
				}, nil)
				return "Update Handler Test"
			},
			run: func(t *testing.T, title string, handler *connectrpc.Handler) {
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
			prepare: func(mockRepo *mock.MockTaskRepository) string {
				mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(model.Task{
					ID:          "delete-test-id",
					Title:       "Delete Handler Test",
					IsCompleted: false,
					CreatedAt:   time.Now(),
				}, nil)
				mockRepo.EXPECT().DeleteTask(ctx, "delete-test-id").Return(nil)
				mockRepo.EXPECT().GetTask(ctx, "delete-test-id").Return(model.Task{}, repository.ErrTaskNotFound)
				return "Delete Handler Test"
			},
			run: func(t *testing.T, title string, handler *connectrpc.Handler) {
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
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			mockRepo := mock.NewMockTaskRepository(ctrl)
			service := todo.NewService(mockRepo, logger)
			handler := connectrpc.NewHandler(service)

			title := tt.prepare(mockRepo)
			tt.run(t, title, handler)
		})
	}
}
