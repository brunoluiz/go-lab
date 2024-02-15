package main

import (
	"context"

	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/xgin"
	"github.com/brunoluiz/go-lab/core/xlog"
	"github.com/brunoluiz/go-lab/services/radars/internal/config"
	"github.com/brunoluiz/go-lab/services/radars/internal/handler"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
	"github.com/brunoluiz/go-lab/services/radars/internal/xhttp"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		xlog.Module,
		xgin.Module,
		postgres.Module,
		repo.Module,
		handler.Module,
		xhttp.Module,
		fx.Provide(
			fx.Annotate(func() context.Context {
				return context.Background()
			}, fx.As(new(context.Context))),
		),
		fx.Invoke(xhttp.Serve),
	).Run()
}
