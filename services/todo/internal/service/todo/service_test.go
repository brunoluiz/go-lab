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
	t.Parallel()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	kv := database.NewKVStore()
	repo := repository.NewTaskRepository(kv, logger)
	service := NewService(repo, logger)
	ctx := context.Background()

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "CreateTask",
			run: func(t *testing.T) {
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
						assert: func(t *testing.T, req dto.CreateTaskRequest, resp dto.CreateTaskResponse, err error, wantErr bool) {
							if wantErr {
								assert.Error(t, err)
							} else {
								require.NoError(t, err)
								assert.NotEmpty(t, resp.Task.ID)
								assert.Equal(t, resp.Task.Title, req.Title)
							}
						},
					},
					{
						name: "success",
						prepare: func() (dto.CreateTaskRequest, bool) {
							return dto.CreateTaskRequest{Title: "Valid Task"}, false
						},
						assert: func(t *testing.T, req dto.CreateTaskRequest, resp dto.CreateTaskResponse, err error, wantErr bool) {
							if wantErr {
								assert.Error(t, err)
							} else {
								require.NoError(t, err)
								assert.NotEmpty(t, resp.Task.ID)
								assert.Equal(t, resp.Task.Title, req.Title)
							}
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
						assert: func(t *testing.T, req dto.GetTaskRequest, resp dto.GetTaskResponse, err error, wantErr bool) {
							assert.Error(t, err)
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
						assert: func(t *testing.T, req dto.UpdateTaskRequest, resp dto.UpdateTaskResponse, err error, wantErr bool) {
							assert.Error(t, err)
						},
					},
					{
						name: "validation title required",
						prepare: func() (dto.UpdateTaskRequest, bool) {
							return dto.UpdateTaskRequest{Task: dto.Task{ID: "123", Title: ""}}, true
						},
						assert: func(t *testing.T, req dto.UpdateTaskRequest, resp dto.UpdateTaskResponse, err error, wantErr bool) {
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
