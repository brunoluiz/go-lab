package handler

import (
	"github.com/brunoluiz/go-lab/services/todo/gen/openapi"
	"github.com/brunoluiz/go-lab/services/todo/internal/repo"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type Handler struct {
	openapi.StrictServerInterface

	repo repo.Querier
}

func Register(r *gin.Engine, repo repo.Querier, log *slog.Logger) {
	h := openapi.NewStrictHandler(&Handler{repo: repo}, []openapi.StrictMiddlewareFunc{})

	schema, _ := openapi.GetSwagger()
	r.Use(
		middleware.OapiRequestValidator(schema),
	)
	r = openapi.RegisterHandlers(r, h)
}
