package handler

import (
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type Handler struct {
	openapi.StrictServerInterface
	Repo   repo.Querier
	WithTx repo.Tx
}

func Register(
	r *gin.Engine,
	handler *Handler,
	log *slog.Logger,
) {
	h := openapi.NewStrictHandler(handler, []openapi.StrictMiddlewareFunc{})

	schema, _ := openapi.GetSwagger()
	r.Use(
		middleware.OapiRequestValidator(schema),
	)
	openapi.RegisterHandlers(r, h)
}
