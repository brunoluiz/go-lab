package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type HTTPConfig struct {
	Address string `envconfig:"address" default:"127.0.0.1"`
	Port    string `envconfig:"port" default:"8080"`
}

func (c *HTTPConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

func NewGin() *gin.Engine {
	// dunno why the default is debug mode and why it is a global variable
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// define default middlewares
	r.Use(gin.Logger(), gin.Recovery())
	r.SetTrustedProxies(nil)

	return r
}
