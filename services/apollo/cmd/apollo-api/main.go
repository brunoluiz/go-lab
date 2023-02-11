package main

import (
	"fmt"
	"log"

	"github.com/brunoluiz/go-lab/services/apollo/internal/handler"
	"github.com/brunoluiz/go-lab/services/apollo/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	HTTP struct {
		Address string `envconfig:"address" default:"127.0.0.1"`
		Port    string `envconfig:"port" default:"8080"`
	}
}

func main() {
	r := gin.Default()
	var c config
	err := envconfig.Process("apollo_api", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	handler.Register(r,
		handler.List(repo.List()),
		handler.Task(),
	)

	r.Run(fmt.Sprintf("%s:%s", c.HTTP.Address, c.HTTP.Port))
}
