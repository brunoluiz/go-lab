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
)

type ListRepository interface {
	CreateList(ctx context.Context, req model.List) (model.List, error)
	GetList(ctx context.Context, id string) (model.List, error)
	ListLists(ctx context.Context) ([]model.List, error)
	UpdateList(ctx context.Context, list model.List) (model.List, error)
	DeleteList(ctx context.Context, id string) error
}

type listRepository struct {
	db     bob.Executor
	logger *slog.Logger
}

func NewListRepository(db bob.Executor, logger *slog.Logger) ListRepository {
	return &listRepository{db: db, logger: logger}
}

func (r *listRepository) CreateList(ctx context.Context, list model.List) (model.List, error) {
	setter := models.ListSetter{
		ID:        omit.From(list.ID),
		Name:      omit.From(list.Name),
		CreatedAt: omit.From(list.CreatedAt),
	}
	created, err := models.Lists.Insert(&setter).One(ctx, r.db)
	if err != nil {
		return model.List{}, errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	return model.List{
		ID:        created.ID,
		Name:      created.Name,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (r *listRepository) GetList(ctx context.Context, id string) (model.List, error) {
	list, err := models.FindList(ctx, r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.List{}, errx.ErrNotFound.Wrapf(err, "resource not found")
		}
		return model.List{}, errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	return model.List{
		ID:        list.ID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
	}, nil
}

func (r *listRepository) ListLists(ctx context.Context) ([]model.List, error) {
	lists, err := models.Lists.Query().All(ctx, r.db)
	if err != nil {
		return nil, errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	var result []model.List
	for _, list := range lists {
		result = append(result, model.List{
			ID:        list.ID,
			Name:      list.Name,
			CreatedAt: list.CreatedAt,
		})
	}
	return result, nil
}

func (r *listRepository) UpdateList(ctx context.Context, list model.List) (model.List, error) {
	bobList, err := models.FindList(ctx, r.db, list.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.List{}, errx.ErrNotFound.Wrapf(err, "resource not found")
		}
		return model.List{}, errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	setter := models.ListSetter{
		Name:      omit.From(list.Name),
		CreatedAt: omit.From(list.CreatedAt),
	}
	err = bobList.Update(ctx, r.db, &setter)
	if err != nil {
		return model.List{}, errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	return model.List{
		ID:        bobList.ID,
		Name:      bobList.Name,
		CreatedAt: bobList.CreatedAt,
	}, nil
}

func (r *listRepository) DeleteList(ctx context.Context, id string) error {
	bobList, err := models.FindList(ctx, r.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errx.ErrNotFound.Wrapf(err, "resource not found")
		}
		return errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	err = bobList.Delete(ctx, r.db)
	if err != nil {
		return errx.ErrUnknown.Wrapf(err, "unknown error")
	}
	return nil
}
