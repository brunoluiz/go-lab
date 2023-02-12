package main

import (
	"fmt"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/xlog"
	"github.com/brunoluiz/go-lab/services/apollo/gen/sqlc/lists"
	"github.com/brunoluiz/go-lab/services/apollo/internal/db"
	"github.com/brunoluiz/go-lab/services/apollo/internal/handler"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	app.CommonConfig
	HTTP app.HTTPConfig
	DB   postgres.EnvConfig
}

func main() {
	logger := xlog.New()

	var c config
	err := envconfig.Process("apollo_api", &c)
	if err != nil {
		logger.Error("problem reading envconfig", err)
		return
	}

	db, err := postgres.New(c.DB.DSN,
		postgres.WithMigration(db.MigrationsFS),
		postgres.WithLiveCheck(),
	)
	if err != nil {
		logger.Error("problem setting up postgres", err)
		return
	}

	r := app.NewGin()
	handler.Register(r, lists.New(db))

	logger.Info(fmt.Sprintf("listening at %s", c.HTTP.GetAddress()))
	r.Run(c.HTTP.GetAddress())
}
