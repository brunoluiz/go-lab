package app

type CommonConfig struct {
	Env string `envconfig:"env" default:"production"`
}

type Env string

const (
	EnvProduction Env = "production"
	EnvLocal      Env = "local"
)
