package xhttp

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/brunoluiz/go-lab/services/radars"
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/config"
	"github.com/brunoluiz/go-lab/services/radars/internal/handler"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
)

func RegisterRoutes(r *gin.Engine, h *handler.Handler) error {
	schema, err := openapi.GetSwagger()
	if err != nil {
		return err
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	// NOTE: delete once kin-openapi fixes behaviour that disappears with .paths (should be able to simply use the embedded openapi.GetSwagger() instead)
	r.StaticFS("/fs", http.FS(radars.OpenAPIFS))
	r.GET("/api/openapi.yaml", func(c *gin.Context) {
		u, _ := url.Parse("http://localhost:8080/fs/openapi/main.yaml")
		doc, err := loader.LoadFromURI(u)
		if err != nil {
			c.Error(err)
			return
		}

		c.YAML(200, doc)
	})
	// NOTE: end of temporary code

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
