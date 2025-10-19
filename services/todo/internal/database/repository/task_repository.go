package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/aarondl/opt/omit"
	"github.com/brunoluiz/go-lab/lib/errx"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/bob/models"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, req model.Task) (model.Task, error)
	GetTask(ctx context.Context, id string) (model.Task, error)
	ListTasks(ctx context.Context, listID string) ([]model.Task, error)
	UpdateTask(ctx context.Context, task model.Task) (model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskRepository struct {
	db     bob.Executor
	logger *slog.Logger
}

func NewTaskRepository(db bob.Executor, logger *slog.Logger) TaskRepository {
	return &taskRepository{db: db, logger: logger}
}

func (r *taskRepository) CreateTask(ctx context.Context, task model.Task) (model.Task, error) {
	setter := models.TaskSetter{
		ID:          omit.From(task.ID),
		Title:       omit.From(task.Title),
		IsCompleted: omit.From(task.IsCompleted),
		CreatedAt:   omit.From(task.CreatedAt),
		ListID:      omit.From(task.ListID),
	}
	created, err := models.Tasks.Insert(&setter).One(ctx, r.db)
	if err != nil {
		return model.Task{}, errx.ErrUnknown.Wrap(err)
	}
	return model.Task{
		ID:          created.ID,
		Title:       created.Title,
		IsCompleted: created.IsCompleted,
		CreatedAt:   created.CreatedAt,
		ListID:      created.ListID,
	}, nil
}

func (r *taskRepository) GetTask(ctx context.Context, id string) (model.Task, error) {
	task, err := models.FindTask(ctx, r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, errx.ErrNotFound.Wrap(err)
		}
		return model.Task{}, errx.ErrUnknown.Wrap(err)
	}
	return model.Task{
		ID:          task.ID,
		Title:       task.Title,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		ListID:      task.ListID,
	}, nil
}

func (r *taskRepository) ListTasks(ctx context.Context, listID string) ([]model.Task, error) {
	tasks, err := models.Tasks.Query(
		sm.Where(models.Tasks.Columns.ListID.EQ(psql.Arg(listID))),
	).All(ctx, r.db)
	if err != nil {
		return nil, errx.ErrUnknown.Wrap(err)
	}
	var result []model.Task
	for _, task := range tasks {
		result = append(result, model.Task{
			ID:          task.ID,
			Title:       task.Title,
			IsCompleted: task.IsCompleted,
			CreatedAt:   task.CreatedAt,
			ListID:      task.ListID,
		})
	}
	return result, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task model.Task) (model.Task, error) {
	bobTask, err := models.FindTask(ctx, r.db, task.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, errx.ErrNotFound.Wrap(err)
		}
		return model.Task{}, errx.ErrUnknown.Wrap(err)
	}
	setter := models.TaskSetter{
		Title:       omit.From(task.Title),
		IsCompleted: omit.From(task.IsCompleted),
		CreatedAt:   omit.From(task.CreatedAt),
		ListID:      omit.From(task.ListID),
	}
	err = bobTask.Update(ctx, r.db, &setter)
	if err != nil {
		return model.Task{}, errx.ErrUnknown.Wrap(err)
	}
	return model.Task{
		ID:          bobTask.ID,
		Title:       bobTask.Title,
		IsCompleted: bobTask.IsCompleted,
		CreatedAt:   bobTask.CreatedAt,
		ListID:      bobTask.ListID,
	}, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id string) error {
	bobTask, err := models.FindTask(ctx, r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errx.ErrNotFound.Wrap(err)
		}
		return errx.ErrUnknown.Wrap(err)
	}
	err = bobTask.Delete(ctx, r.db)
	if err != nil {
		return errx.ErrUnknown.Wrap(err)
	}
	return nil
}
