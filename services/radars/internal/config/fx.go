package config

import (
	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"go.uber.org/fx"
)

var Module = fx.Module("config", fx.Provide(
	New,
	func(c *Config) postgres.EnvConfig {
		return c.DB
	},
))
