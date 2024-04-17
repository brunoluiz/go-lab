package provider

import (
	"context"

	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/xgin"
	"github.com/brunoluiz/go-lab/core/xlog"
	"github.com/brunoluiz/go-lab/services/radars/internal/handler"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
	"github.com/brunoluiz/go-lab/services/radars/internal/xhttp"
	"go.uber.org/fx"
)

func InjectApp() fx.Option {
	return fx.Module("app",
		fx.Provide(
			fx.Annotate(func() context.Context {
				return context.Background()
			}, fx.As(new(context.Context))),
		),
		xlog.Module,
		xgin.Module,
		postgres.Module,
		repo.Module,
		handler.Module,
		xhttp.Module,
	)
}
