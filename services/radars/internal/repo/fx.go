package repo

import (
	"database/sql"

	"go.uber.org/fx"
)

var Module = fx.Module("handler",
	fx.Supply(MigrationsFS),
	fx.Provide(
		fx.Annotate(New, fx.As(new(Querier)), fx.As(new(QuerierTx))),
		fx.Annotate(NewTxExec, fx.As(new(Tx))),
		func(db *sql.DB) DBTX { return db },
	))
