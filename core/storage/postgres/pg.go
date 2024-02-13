package postgres

import (
	"database/sql"
	"errors"
	"io/fs"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/stdlib"
)

type EnvConfig struct {
	DSN string `envconfig:"db_dsn" required:"true"`
}

type config struct {
	migration fs.FS
	ping      bool
}

type option func(*config)

func WithMigration(embed fs.FS) func(*config) {
	return func(c *config) {
		c.migration = embed
	}
}

func WithLiveCheck() func(*config) {
	return func(c *config) {
		c.ping = true
	}
}

func New(cfg EnvConfig, opts ...option) (*sql.DB, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if c.ping {
		if err := db.Ping(); err != nil {
			return nil, err
		}
	}

	if c.migration != nil {
		if err := up(db, c.migration); err != nil {
			return nil, err
		}
	}

	return db, err
}

func up(db *sql.DB, fs fs.FS) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	src, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		spew.Dump(err)
		return err
	}

	return nil
}
