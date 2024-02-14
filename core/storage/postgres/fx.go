package postgres

import (
	"database/sql"
	"embed"

	"go.uber.org/fx"
)

var Module = fx.Module("postgres", fx.Provide(
	func(params struct {
		fx.In

		Cfg       EnvConfig
		Migration embed.FS
	},
	) (*sql.DB, error) {
		return New(params.Cfg,
			WithMigration(params.Migration),
			WithLiveCheck(),
		)
	},
))
