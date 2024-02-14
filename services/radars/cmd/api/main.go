package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/xgin"
	"github.com/brunoluiz/go-lab/core/xlog"
	"github.com/brunoluiz/go-lab/services/radars/internal/config"
	"github.com/brunoluiz/go-lab/services/radars/internal/handler"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		xlog.Module,
		xgin.Module,
		repo.Module,
		handler.Module,
		postgres.Module,
		fx.Provide(
			fx.Annotate(func() context.Context {
				return context.Background()
			}, fx.As(new(context.Context))),
		),
		fx.Invoke(func(c *config.Config, h *handler.Handler, r *gin.Engine, l *slog.Logger) {
			h.Register(r)
			l.Info(fmt.Sprintf("listening at %s", c.HTTP.GetAddress()))
			r.Run(c.HTTP.GetAddress())
		}),
	).Run()
}
