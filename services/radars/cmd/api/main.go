package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/xgin"
	"github.com/brunoluiz/go-lab/core/xlog"
	"github.com/brunoluiz/go-lab/services/radars/internal/handler"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

type config struct {
	app.CommonConfig
	HTTP xgin.HTTPConfig
	DB   postgres.EnvConfig
}

func main() {
	ctx := func() context.Context {
		return context.Background()
	}

	var c config
	err := envconfig.Process("radars_api", &c)
	if err != nil {
		panic(fmt.Errorf("problem reading envconfig: %s", err))
	}

	db, err := postgres.New(c.DB,
		postgres.WithMigration(repo.MigrationsFS),
		postgres.WithLiveCheck(),
	)
	if err != nil {
		panic(err)
	}

	fx.New(
		fx.Supply(c),
		fx.Supply(c.DB),
		fx.Supply(c.HTTP),
		fx.Supply(db),
		fx.Provide(
			// When context.Context is required in providers, a new context.Background will be provided
			fx.Annotate(ctx, fx.As(new(context.Context))),
			func() repo.DBTX { return db },
			// fx.Annotate(func() (*sql.DB, error) {
			//   return postgres.New(c.DB,
			//     postgres.WithMigration(repo.MigrationsFS),
			//     postgres.WithLiveCheck(),
			//   )
			// }, fx.As(new(sql.DB))),
			xlog.New,
			xgin.New,
			// postgres.New,
			fx.Annotate(
				repo.New,
				fx.As(new(repo.Querier)),
			),
			repo.New,
			repo.NewTxExec,
			handler.New,
		),
		fx.Invoke(func(h *handler.Handler, r *gin.Engine, l *slog.Logger) {
			h.Register(r)
			l.Info(fmt.Sprintf("listening at %s", c.HTTP.GetAddress()))
			r.Run(c.HTTP.GetAddress())
		}),
	).Run()
}
