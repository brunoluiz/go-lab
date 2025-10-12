package todo_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/dto"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/todo/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTodoService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "CreateTask",
			run: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
				mockRepo := mock.NewMockTaskRepository(ctrl)
				service := todo.NewService(mockRepo, logger)

				subTests := []struct {
					name    string
					prepare func() (dto.CreateTaskRequest, bool)
					assert  func(t *testing.T, req dto.CreateTaskRequest, resp dto.CreateTaskResponse, err error, wantErr bool)
				}{
					{
						name: "validation empty title",
						prepare: func() (dto.CreateTaskRequest, bool) {
							return dto.CreateTaskRequest{Title: ""}, true
						},
						assert: func(t *testing.T, _ dto.CreateTaskRequest, _ dto.CreateTaskResponse, err error, _ bool) {
							assert.Error(t, err)
						},
					},
					{
						name: "success",
						prepare: func() (dto.CreateTaskRequest, bool) {
							mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(model.Task{
								ID:          "test-id",
								Title:       "Valid Task",
								IsCompleted: false,
								CreatedAt:   time.Now(),
							}, nil)
							return dto.CreateTaskRequest{Title: "Valid Task"}, false
						},
						assert: func(t *testing.T, _ dto.CreateTaskRequest, resp dto.CreateTaskResponse, err error, _ bool) {
							require.NoError(t, err)
							assert.NotEmpty(t, resp.Task.ID)
							assert.Equal(t, resp.Task.Title, "Valid Task")
						},
					},
				}

				for _, tt := range subTests {
					t.Run(tt.name, func(t *testing.T) {
						req, wantErr := tt.prepare()
						resp, err := service.CreateTask(ctx, req)
						tt.assert(t, req, resp, err, wantErr)
					})
				}
			},
		},
		{
			name: "GetTask",
			run: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
				mockRepo := mock.NewMockTaskRepository(ctrl)
				service := todo.NewService(mockRepo, logger)

				subTests := []struct {
					name    string
					prepare func() (dto.GetTaskRequest, bool)
					assert  func(t *testing.T, req dto.GetTaskRequest, resp dto.GetTaskResponse, err error, wantErr bool)
				}{
					{
						name: "validation empty id",
						prepare: func() (dto.GetTaskRequest, bool) {
							return dto.GetTaskRequest{TaskID: ""}, true
						},
						assert: func(t *testing.T, _ dto.GetTaskRequest, _ dto.GetTaskResponse, err error, _ bool) {
							assert.Error(t, err)
						},
					},
					{
						name: "success",
						prepare: func() (dto.GetTaskRequest, bool) {
							mockRepo.EXPECT().GetTask(ctx, "test-id").Return(model.Task{
								ID:          "test-id",
								Title:       "Test Task",
								IsCompleted: false,
								CreatedAt:   time.Now(),
							}, nil)
							return dto.GetTaskRequest{TaskID: "test-id"}, false
						},
						assert: func(t *testing.T, _ dto.GetTaskRequest, resp dto.GetTaskResponse, err error, _ bool) {
							require.NoError(t, err)
							assert.Equal(t, "test-id", resp.Task.ID)
							assert.Equal(t, "Test Task", resp.Task.Title)
						},
					},
				}

				for _, tt := range subTests {
					t.Run(tt.name, func(t *testing.T) {
						req, wantErr := tt.prepare()
						resp, err := service.GetTask(ctx, req)
						tt.assert(t, req, resp, err, wantErr)
					})
				}
			},
		},
		{
			name: "UpdateTask",
			run: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
				mockRepo := mock.NewMockTaskRepository(ctrl)
				service := todo.NewService(mockRepo, logger)

				subTests := []struct {
					name    string
					prepare func() (dto.UpdateTaskRequest, bool)
					assert  func(t *testing.T, req dto.UpdateTaskRequest, resp dto.UpdateTaskResponse, err error, wantErr bool)
				}{
					{
						name: "validation empty id and title",
						prepare: func() (dto.UpdateTaskRequest, bool) {
							return dto.UpdateTaskRequest{Task: dto.Task{ID: "", Title: ""}}, true
						},
						assert: func(t *testing.T, _ dto.UpdateTaskRequest, _ dto.UpdateTaskResponse, err error, _ bool) {
							assert.Error(t, err)
						},
					},
					{
						name: "validation title required",
						prepare: func() (dto.UpdateTaskRequest, bool) {
							return dto.UpdateTaskRequest{Task: dto.Task{ID: "123", Title: ""}}, true
						},
						assert: func(t *testing.T, _ dto.UpdateTaskRequest, _ dto.UpdateTaskResponse, err error, _ bool) {
							assert.Error(t, err)
						},
					},
				}

				for _, tt := range subTests {
					t.Run(tt.name, func(t *testing.T) {
						req, wantErr := tt.prepare()
						resp, err := service.UpdateTask(ctx, req)
						tt.assert(t, req, resp, err, wantErr)
					})
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.run(t)
		})
	}
}
