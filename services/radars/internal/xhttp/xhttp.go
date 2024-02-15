package xhttp

import (
	"fmt"
	"log/slog"

	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/config"
	"github.com/brunoluiz/go-lab/services/radars/internal/handler"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
)

func RegisterRoutes(r *gin.Engine, h *handler.Handler) error {
	schema, err := openapi.GetSwagger()
	if err != nil {
		return err
	}

	r.Use(
		ginmiddleware.OapiRequestValidator(schema),
	)
	openapi.RegisterHandlers(r, openapi.NewStrictHandler(h, []openapi.StrictMiddlewareFunc{}))
	return nil
}

func Serve(c *config.Config, r *gin.Engine, l *slog.Logger) {
	l.Info(fmt.Sprintf("listening at %s", c.HTTP.GetAddress()))
	r.Run(c.HTTP.GetAddress())
}
