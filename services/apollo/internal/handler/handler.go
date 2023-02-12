package handler

import (
	"context"

	"github.com/brunoluiz/go-lab/services/apollo/gen/openapi"
	"github.com/brunoluiz/go-lab/services/apollo/gen/sqlc/lists"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type Handler struct {
	openapi.StrictServerInterface

	repo lists.Querier
}

func (h *Handler) AddList(ctx context.Context, req openapi.AddListRequestObject) (openapi.AddListResponseObject, error) {
	out, err := h.repo.Create(ctx, lists.CreateParams{
		UID:   ksuid.New().String(),
		Title: req.Body.Title,
	})
	if err != nil {
		return openapi.AddList400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return openapi.AddList200JSONResponse{
		Title: out.Title,
		Uid:   &out.UID,
	}, nil
}

func (*Handler) UpdateList(ctx context.Context, request openapi.UpdateListRequestObject) (openapi.UpdateListResponseObject, error) {
	panic("not implemented")
}

func (*Handler) DeleteList(ctx context.Context, request openapi.DeleteListRequestObject) (openapi.DeleteListResponseObject, error) {
	panic("not implemented")
}

func (h *Handler) GetListById(ctx context.Context, req openapi.GetListByIdRequestObject) (openapi.GetListByIdResponseObject, error) {
	out, err := h.repo.ByUID(ctx, req.ListId)
	if err != nil {
		return openapi.GetListById404JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return openapi.GetListById200JSONResponse{
		Uid:   &out.UID,
		Title: out.Title,
	}, nil
}

func Register(r *gin.Engine, l lists.Querier) {
	h := openapi.NewStrictHandler(&Handler{repo: l}, nil)

	schema, _ := openapi.GetSwagger()
	r.Use(middleware.OapiRequestValidator(schema))
	r = openapi.RegisterHandlers(r, h)
}
