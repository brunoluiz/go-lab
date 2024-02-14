package config

import (
	"fmt"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/xgin"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	app.CommonConfig
	HTTP xgin.HTTPConfig
	DB   postgres.EnvConfig
}

func New() (*Config, error) {
	var c Config
	if err := envconfig.Process("radars_api", &c); err != nil {
		return nil, fmt.Errorf("problem reading envconfig: %s", err)
	}
	return &c, nil
}
