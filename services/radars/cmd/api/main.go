package main

import (
	"github.com/brunoluiz/go-lab/services/radars/cmd/api/provider"
	"github.com/brunoluiz/go-lab/services/radars/internal/config"
	"github.com/brunoluiz/go-lab/services/radars/internal/xhttp"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		provider.InjectApp(),
		fx.Invoke(xhttp.Serve),
	).Run()
}
