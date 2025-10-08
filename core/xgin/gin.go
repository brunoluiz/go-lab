package xgin

import (
	"fmt"
	"log/slog"

	"github.com/brunoluiz/go-lab/core/xgin/xmiddleware"
	"github.com/gin-gonic/gin"
)

type HTTPConfig struct {
	Address string `envconfig:"address" default:"127.0.0.1"`
	Port    string `envconfig:"port" default:"8080"`
}

func (c *HTTPConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

func New(log *slog.Logger) *gin.Engine {
	// dunno why the default is debug mode and why it is a global variable
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// define default middlewares
	r.Use(
		gin.Recovery(),
		xmiddleware.ErrorHandler(log),
	)
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Warn("failed to set trusted proxies", "error", err)
	}

	return r
}
