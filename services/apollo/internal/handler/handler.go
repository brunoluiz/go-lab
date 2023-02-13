package handler

import (
	"context"

	"github.com/brunoluiz/go-lab/services/apollo/gen/openapi"
	"github.com/brunoluiz/go-lab/services/apollo/internal/repo"
	"github.com/davecgh/go-spew/spew"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type Handler struct {
	openapi.StrictServerInterface

	repo repo.Querier
}

func (h *Handler) AddList(ctx context.Context, req openapi.AddListRequestObject) (openapi.AddListResponseObject, error) {
	out, err := h.repo.SaveList(ctx, repo.SaveListParams{
		UniqID: ksuid.New().String(),
		Title:  req.Body.Title,
	})
	if err != nil {
		return openapi.AddList400JSONResponse{
			Message: err.Error(),
		}, nil
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
		return openapi.UpdateList400JSONResponse{
			Message: err.Error(),
		}, nil
	}
	spew.Dump(req)

	return openapi.UpdateList200JSONResponse{
		Title:     out.Title,
		UniqId:    out.UniqID,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}, nil
}

func (h *Handler) DeleteList(ctx context.Context, req openapi.DeleteListRequestObject) (openapi.DeleteListResponseObject, error) {
	if err := h.repo.DeleteList(ctx, req.ListId); err != nil {
		return openapi.DeleteList404JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return openapi.DeleteList200JSONResponse{
		Success: true,
	}, nil
}

func (h *Handler) GetListById(ctx context.Context, req openapi.GetListByIdRequestObject) (openapi.GetListByIdResponseObject, error) {
	out, err := h.repo.GetListByID(ctx, req.ListId)
	if err != nil {
		return openapi.GetListById404JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return openapi.GetListById200JSONResponse{
		UniqId:    out.UniqID,
		Title:     out.Title,
		CreatedAt: out.CreatedAt,
	}, nil
}

func Register(r *gin.Engine, l repo.Querier) {
	h := openapi.NewStrictHandler(&Handler{repo: l}, nil)

	schema, _ := openapi.GetSwagger()
	r.Use(middleware.OapiRequestValidator(schema))
	r = openapi.RegisterHandlers(r, h)
}
