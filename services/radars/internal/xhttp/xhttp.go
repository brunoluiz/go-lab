package xhttp

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brunoluiz/go-lab/services/radars"
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

	// NOTE: delete once kin-openapi fixes behaviour that disappears with .paths (should be able to simply use the embedded openapi.GetSwagger() instead)
	r.StaticFS("/__", http.FS(radars.OpenAPIFS))
	// NOTE: end of temporary code

	openapi.RegisterHandlersWithOptions(r,
		openapi.NewStrictHandler(h, []openapi.StrictMiddlewareFunc{}),
		openapi.GinServerOptions{
			Middlewares: []openapi.MiddlewareFunc{
				openapi.MiddlewareFunc(ginmiddleware.OapiRequestValidatorWithOptions(schema, &ginmiddleware.Options{
					ErrorHandler: func(c *gin.Context, message string, statusCode int) {
						c.JSON(statusCode, map[string]string{
							"status":  "fail",
							"message": message,
						})
					},
				})),
			},
		})
	return nil
}

func Serve(c *config.Config, r *gin.Engine, l *slog.Logger) {
	l.Info(fmt.Sprintf("listening at %s", c.HTTP.GetAddress()))
	r.Run(c.HTTP.GetAddress())
}
