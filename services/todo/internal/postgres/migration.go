package postgres

import (
	"database/sql"
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	pgdriver "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migration/*.sql
var migrations embed.FS

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(db *sql.DB, _ *slog.Logger) (*Migrator, error) {
	source, err := iofs.New(&migrations, ".")
	if err != nil {
		return nil, err
	}

	driver, err := pgdriver.WithInstance(db, &pgdriver.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return nil, err
	}

	return &Migrator{migrate: m}, nil
}

func (m *Migrator) Up() error {
	return m.migrate.Up()
}
