package handler

import (
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	openapi.StrictServerInterface
	Repo repo.Querier
	Tx   repo.Tx
}

func New(
	q repo.Querier,
	withTx repo.Tx,
) *Handler {
	return &Handler{
		Repo: q,
		Tx:   withTx,
	}
}

func (h *Handler) Register(r *gin.Engine) {
	schema, _ := openapi.GetSwagger()
	r.Use(
		middleware.OapiRequestValidator(schema),
	)
	openapi.RegisterHandlers(r, openapi.NewStrictHandler(h, []openapi.StrictMiddlewareFunc{}))
}
