package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/services/todo/gen/openapi"
	"github.com/brunoluiz/go-lab/services/todo/internal/repo"
	"github.com/segmentio/ksuid"
)

func (h *Handler) AddList(ctx context.Context, req openapi.AddListRequestObject) (openapi.AddListResponseObject, error) {
	out, err := h.repo.SaveList(ctx, repo.SaveListParams{
		UniqID: ksuid.New().String(),
		Title:  req.Body.Title,
	})
	if err != nil {
		return nil, err
	}

	return openapi.AddList201JSONResponse{
		Title:     out.Title,
		UniqId:    out.UniqID,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}, nil
}

func (h *Handler) UpdateList(ctx context.Context, req openapi.UpdateListRequestObject) (openapi.UpdateListResponseObject, error) {
	out, err := h.repo.SaveList(ctx, repo.SaveListParams{
		UniqID: req.ListId,
		Title:  req.Body.Title,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &app.ErrNotFound{}
		}
		return nil, err
	}

	return openapi.UpdateList200JSONResponse{
		Title:     out.Title,
		UniqId:    out.UniqID,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}, nil
}

func (h *Handler) DeleteList(ctx context.Context, req openapi.DeleteListRequestObject) (openapi.DeleteListResponseObject, error) {
	err := h.repo.DeleteList(ctx, req.ListId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &app.ErrNotFound{}
		}
		return nil, err
	}

	return openapi.DeleteList200JSONResponse{
		Success: true,
	}, nil
}

func (h *Handler) GetListById(ctx context.Context, req openapi.GetListByIdRequestObject) (openapi.GetListByIdResponseObject, error) {
	out, err := h.repo.GetListByID(ctx, req.ListId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &app.ErrNotFound{}
		}
		return nil, err
	}

	return openapi.GetListById200JSONResponse{
		UniqId:    out.UniqID,
		Title:     out.Title,
		CreatedAt: out.CreatedAt,
	}, nil
}
