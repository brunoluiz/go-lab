package postgres

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"log/slog"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/brunoluiz/go-lab/lib/errx"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/stdlib" // registers pgx driver for database/sql
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type config struct {
	migration       fs.FS
	ping            bool
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	logger          *slog.Logger
	connTimeout     time.Duration
	maxRetries      int
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

func WithMaxOpenConns(n int) func(*config) {
	return func(c *config) {
		c.maxOpenConns = n
	}
}

func WithMaxIdleConns(n int) func(*config) {
	return func(c *config) {
		c.maxIdleConns = n
	}
}

func WithConnMaxLifetime(d time.Duration) func(*config) {
	return func(c *config) {
		c.connMaxLifetime = d
	}
}

func WithConnMaxIdleTime(d time.Duration) func(*config) {
	return func(c *config) {
		c.connMaxIdleTime = d
	}
}

func WithLogger(logger *slog.Logger) func(*config) {
	return func(c *config) {
		c.logger = logger
	}
}

func WithConnTimeout(timeout time.Duration) func(*config) {
	return func(c *config) {
		c.connTimeout = timeout
	}
}

func WithMaxRetries(retries int) func(*config) {
	return func(c *config) {
		c.maxRetries = retries
	}
}

func New(dsn string, opts ...option) (*sql.DB, error) {
	c := &config{
		maxOpenConns:    25,
		maxIdleConns:    5,
		connMaxLifetime: 5 * time.Minute,
		connMaxIdleTime: 5 * time.Minute,
		connTimeout:     30 * time.Second,
		maxRetries:      3,
	}
	for _, opt := range opts {
		opt(c)
	}

	db, err := otelsql.Open("pgx", dsn, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(c.maxOpenConns)
	db.SetMaxIdleConns(c.maxIdleConns)
	db.SetConnMaxLifetime(c.connMaxLifetime)
	db.SetConnMaxIdleTime(c.connMaxIdleTime)

	if err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	)); err != nil {
		return nil, err
	}

	if c.ping {
		if pingErr := pingWithRetry(db, c.connTimeout, c.maxRetries, c.logger); pingErr != nil {
			return nil, pingErr
		}
	}

	if c.migration != nil {
		if migrateErr := up(db, c.migration, c.logger); migrateErr != nil {
			return nil, migrateErr
		}
	}

	return db, nil
}

func up(db *sql.DB, fs fs.FS, _ *slog.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errx.ErrInternal.Wrap(err)
	}

	src, err := iofs.New(fs, "migrations")
	if err != nil {
		return errx.ErrInternal.New("migrations directory not found in embedded filesystem")
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return errx.ErrInternal.Wrap(err)
	}

	if upErr := m.Up(); !errors.Is(upErr, migrate.ErrNoChange) {
		return errx.ErrInternal.Wrapf(err, "migration failed")
	}

	return nil
}

func pingWithRetry(db *sql.DB, timeout time.Duration, maxRetries int, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.Default()
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			lastErr = err
			if attempt < maxRetries {
				logger.WarnContext(context.Background(), "database ping failed, retrying",
					"attempt", attempt+1,
					"max_retries", maxRetries,
					"error", err)
				time.Sleep(time.Duration(attempt+1) * time.Second) // Exponential backoff
				continue
			}
		} else {
			return nil
		}
	}
	return lastErr
}

// HealthCheck performs a comprehensive health check on the database connection
func HealthCheck(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	// Basic ping check
	if err := db.PingContext(ctx); err != nil {
		return errx.ErrInternal.Wrapf(err, "database health check failed: ping")
	}

	// Check connection pool stats
	stats := db.Stats()
	if logger != nil && logger.Enabled(ctx, slog.LevelDebug) {
		logger.DebugContext(ctx, "database connection pool stats",
			"open_connections", stats.OpenConnections,
			"in_use", stats.InUse,
			"idle", stats.Idle,
			"wait_count", stats.WaitCount,
			"wait_duration", stats.WaitDuration,
			"max_idle_closed", stats.MaxIdleClosed,
			"max_lifetime_closed", stats.MaxLifetimeClosed,
		)
	}

	// Simple query to verify database is responsive
	if err := db.QueryRowContext(ctx, "SELECT 1").Scan(new(int)); err != nil {
		return errx.ErrInternal.Wrapf(err, "database health check failed: query")
	}

	return nil
}
