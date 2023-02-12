package xlog

import (
	"os"

	"github.com/brunoluiz/go-lab/core/app"
	"golang.org/x/exp/slog"
)

type config struct {
	env app.Env
}

func WithEnv(env app.Env) func(*config) {
	return func(c *config) {
		c.env = env
	}
}

type Options func(*config)

func New(opts ...Options) *slog.Logger {
	o := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	c := &config{
		env: app.EnvLocal,
	}

	for _, opt := range opts {
		opt(c)
	}

	switch c.env {
	case app.EnvProduction:
		return slog.New(o.NewJSONHandler(os.Stdout))
	default:
		return slog.New(o.NewTextHandler(os.Stdout))
	}
}
