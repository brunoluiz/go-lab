package xgin

import (
	"fmt"

	"github.com/brunoluiz/go-lab/core/xgin/xmiddleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
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
	r.SetTrustedProxies(nil)

	return r
}
