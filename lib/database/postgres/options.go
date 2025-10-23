package postgres

import (
	"io/fs"
	"time"
)

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
