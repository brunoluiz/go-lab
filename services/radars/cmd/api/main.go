package main

import (
	"fmt"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/xgin"
	"github.com/brunoluiz/go-lab/core/xlog"
	"github.com/brunoluiz/go-lab/services/radars/internal/handler"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	app.CommonConfig
	HTTP xgin.HTTPConfig
	DB   postgres.EnvConfig
}

func main() {
	logger := xlog.New()

	var c config
	err := envconfig.Process("radars_api", &c)
	if err != nil {
		logger.Error("problem reading envconfig", err)
		return
	}

	db, err := postgres.New(c.DB.DSN,
		postgres.WithMigration(repo.MigrationsFS),
		postgres.WithLiveCheck(),
	)
	if err != nil {
		logger.Error("problem setting up postgres", err)
		return
	}

	r := xgin.New(logger)
	sqlRepo := repo.New(db)
	h := handler.New(
		sqlRepo,
		repo.NewTx(db, sqlRepo),
	)
	h.Register(r)

	logger.Info(fmt.Sprintf("listening at %s", c.HTTP.GetAddress()))
	r.Run(c.HTTP.GetAddress())
}
